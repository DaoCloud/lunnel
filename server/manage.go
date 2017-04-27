// Copyright 2017 longXboy, longxboyhi@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/longXboy/lunnel/msg"
	"github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

type tunnelsStateReq struct {
	PublicUrl string
}

type tunnelsStateResp struct {
	Tunnels []string
}

type clientState struct {
	Id             string
	LastRead       uint64
	EncryptMode    string
	EnableCompress bool
	Version        string
	Tunnels        map[string]tunnelState
	TotalPipes     int64
}

type tunnelState struct {
	Tunnel   msg.Tunnel
	IsClosed bool
}

type clientsStateResp struct {
	Clients []clientState
}

func serveManage() {
	r := gin.New()
	if serverConf.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	r.GET("/v1/tunnels", tunnelsQuery)
	r.POST("/v1/tunnel", tunnelQuery)

	r.GET("/v1/clients", clientsQuery)
	r.GET("/v1/clients/clientId", clientQuery)

	r.Run(fmt.Sprintf("%s:%d", serverConf.ListenIP, serverConf.ManagePort))

	http.ListenAndServe(fmt.Sprintf("%s:%d", serverConf.ListenIP, serverConf.ManagePort), nil)
}

func tunnelQuery(c *gin.Context) {
	var query tunnelsStateReq
	err := c.BindJSON(&query)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("unmarshal req body failed!"))
		return
	}

	var tunnelStats tunnelsStateResp = tunnelsStateResp{Tunnels: []string{}}
	if query.PublicUrl != "" {
		TunnelMapLock.RLock()
		tunnel, isok := TunnelMap[query.PublicUrl]
		TunnelMapLock.RUnlock()
		if isok {
			tunnelStats.Tunnels = append(tunnelStats.Tunnels, tunnel.tunnelConfig.PublicAddr())
		}
	}

	c.JSON(http.StatusOK, tunnelStats)
}

func tunnelsQuery(c *gin.Context) {
	var tunnelStats tunnelsStateResp = tunnelsStateResp{Tunnels: []string{}}

	TunnelMapLock.RLock()
	for _, v := range TunnelMap {
		tunnelStats.Tunnels = append(tunnelStats.Tunnels, v.tunnelConfig.PublicAddr())
	}
	TunnelMapLock.RUnlock()

	c.JSON(http.StatusOK, tunnelStats)
}

func clientQuery(c *gin.Context) {
	clientId := c.Param("clientId")
	uuid, err := uuid.FromString(clientId)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("invalid uuid"))
		return
	}
	var clientStates clientsStateResp = clientsStateResp{Clients: []clientState{}}
	ControlMapLock.RLock()
	ctlClient := ControlMap[uuid]
	ControlMapLock.RUnlock()
	var client clientState
	client.LastRead = atomic.LoadUint64(&client.LastRead)
	client.TotalPipes = atomic.LoadInt64(&ctlClient.totalPipes)
	client.Tunnels = make(map[string]tunnelState)
	ctlClient.tunnelLock.Lock()
	for _, v := range ctlClient.tunnels {
		client.Tunnels[v.name] = tunnelState{Tunnel: v.tunnelConfig, IsClosed: v.isClosed}
	}
	ctlClient.tunnelLock.Unlock()
	client.EnableCompress = ctlClient.enableCompress
	client.EncryptMode = ctlClient.encryptMode
	client.Id = ctlClient.ClientID.String()
	client.Version = ctlClient.version
	clientStates.Clients = append(clientStates.Clients, client)
	c.JSON(http.StatusOK, clientStates)
}

func clientsQuery(c *gin.Context) {
	var clientStates clientsStateResp = clientsStateResp{Clients: []clientState{}}
	clients := make([]*Control, 0)
	ControlMapLock.RLock()
	for _, v := range ControlMap {
		clients = append(clients, v)
	}
	ControlMapLock.RUnlock()
	for _, c := range clients {
		var client clientState
		client.LastRead = atomic.LoadUint64(&c.lastRead)
		client.TotalPipes = atomic.LoadInt64(&c.totalPipes)
		client.Tunnels = make(map[string]tunnelState)
		c.tunnelLock.Lock()
		for _, v := range c.tunnels {
			client.Tunnels[v.name] = tunnelState{Tunnel: v.tunnelConfig, IsClosed: v.isClosed}
		}
		c.tunnelLock.Unlock()
		client.EnableCompress = c.enableCompress
		client.EncryptMode = c.encryptMode
		client.Id = c.ClientID.String()
		client.Version = c.version
		clientStates.Clients = append(clientStates.Clients, client)
	}
	c.JSON(http.StatusOK, clientStates)
}
