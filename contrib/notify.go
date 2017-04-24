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

package contrib

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/longXboy/lunnel/msg"
	"gopkg.in/redis.v5"
)

var redisCli *redis.Client
var notifyKey string

func InitNotify(notifyUrl string, nk string) error {
	if notifyUrl == "" {
		log.Fatalln("notifyUrl is empty")
	}
	if nk == "" {
		log.Fatalln("notifyKey is empty")
	}
	notifyKey = nk
	redisCli = redis.NewClient(&redis.Options{
		Addr:     notifyUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisCli.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorln("redis ping failed!")
		return err
	}
	log.WithFields(log.Fields{"pong": pong}).Infoln("init Notify,ping redis success!")
	cmd := redisCli.Del(notifyKey)
	if cmd.Err() != nil {
		log.WithFields(log.Fields{"err": cmd.Err()}).Errorln("init redis,delete key failed!")
		return cmd.Err()
	}
	return nil
}

type PublishTunnelRequest struct {
	PublicUrl string `json:"public_url"`
	LocalAddr string `json:"local_addr"`
}

func AddTunnel(domain string, tunnel msg.Tunnel, clientId string) error {
	var ret error
	var retry int
	for {
		retry++
		if retry > 2 {
			break
		}
		cmd := redisCli.SAdd(notifyKey, tunnel.PublicAddr())
		if cmd.Err() != nil {
			ret = cmd.Err()
			continue
		}
		url := fmt.Sprintf("%s/v2/daomonits/%s/tunnels", daoKeeperUrl, clientId)
		var req PublishTunnelRequest
		req.LocalAddr = tunnel.LocalAddr()
		req.PublicUrl = tunnel.PublicAddr()
		content, err := json.Marshal(req)
		if err != nil {
			ret = err
			continue
		}
		reader := bytes.NewReader(content)
		resp, err := httpClient.Post(url, "application/json", reader)
		if err != nil {
			ret = err
			continue
		}
		if resp.Body != nil {
			resp.Body.Close()
		}
		if resp.StatusCode != 200 {
			ret = err
			continue
		}
		return nil
	}
	if ret != nil {
		log.WithFields(log.Fields{"key": notifyKey, "client_id": clientId, "member": tunnel.PublicAddr()}).Errorln("redis add pub url success!")
	}
	return ret
}

func RemoveTunnel(domain string, tunnel msg.Tunnel, clientId string) error {
	cmd := redisCli.SRem(notifyKey, tunnel.PublicAddr())
	if cmd.Err() != nil {
		return cmd.Err()
	}
	log.WithFields(log.Fields{"key": notifyKey, "member": tunnel.PublicAddr()}).Debugln("redis remove pub url success!")
	return nil
}
