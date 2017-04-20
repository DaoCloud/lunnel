package contrib

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
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

func Auth(authToken string) (bool, error) {
	resp, err := httpClient.PostForm(fmt.Sprintf("%s/v1/ngrokd/auth", daoKeeperUrl), url.Values{"user": {authToken}})
	if err != nil {
		return false, fmt.Errorf("Request daokeeper error %s,%v", fmt.Sprintf("%s/v1/ngrokd/auth", daoKeeperUrl), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{"authtoken": authToken}).Infoln("client auth token success!")
		return true, nil
	} else {
		log.WithFields(log.Fields{"authtoken": authToken, "statuscode": resp.StatusCode}).Errorln("client auth token failed!")
		return false, fmt.Errorf("Response daokeeper code %d,%v", resp.StatusCode)
	}
	return true, nil
}
