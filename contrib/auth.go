package contrib

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/longXboy/lunnel/msg"
)

var httpClient *http.Client

var daoKeeperUrl string

func InitAuth(authUrl string) error {
	if authUrl == "" {
		return fmt.Errorf("auth url not be empty")
	}
	daoKeeperUrl = authUrl

	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   8 * time.Second,
			KeepAlive: 90 * time.Second,
		}).DialContext,
		MaxIdleConns:          12,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   8 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: trans,
	}

	return nil
}

func Auth(chello *msg.ControlClientHello) (bool, error) {
	clientId := ""
	if chello.ClientID != nil {
		clientId = chello.ClientID.String()
	}
	resp, err := httpClient.PostForm(fmt.Sprintf("%s/v1/ngrokd/auth", daoKeeperUrl), url.Values{
		"user": {chello.AuthToken},
		"id":   {clientId},
	})
	if err != nil {
		return false, fmt.Errorf("Request daokeeper error %s,%v", fmt.Sprintf("%s/v1/ngrokd/auth", daoKeeperUrl), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{"authtoken": chello.AuthToken}).Infoln("client auth token success!")
		return true, nil
	} else {
		log.WithFields(log.Fields{"authtoken": chello.AuthToken, "statuscode": resp.StatusCode}).Errorln("client auth token failed!")
		return false, fmt.Errorf("Response daokeeper code %d,%v", resp.StatusCode)
	}
	return true, nil
}
