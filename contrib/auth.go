package contrib

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

var daoUrl *url.URL

func InitAuth(authUrl string) error {
	if authUrl == "" {
		return fmt.Errorf("auth url not be empty")
	}
	var err error
	daoUrl, err = url.Parse(authUrl)
	if err != nil {
		return fmt.Errorf("DAOKEEPER_URL Parse error %s,%v", authUrl, err)
	}
	daoUrl.Path = "/v1/ngrokd/auth"
	return nil
}

func Auth(authToken string) (bool, error) {
	resp, err := http.PostForm(daoUrl.String(), url.Values{"user": {authToken}})
	if err != nil {
		return false, fmt.Errorf("Request daokeeper error %s,%v", daoUrl.String(), err)
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
