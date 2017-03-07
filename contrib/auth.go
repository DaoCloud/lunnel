package contrib

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var daoUrl *url.URL

func InitAuth(authUrl string) error {
	if authUrl == "" {
		return fmt.Errorf("auth url not be empty")
	}
	var err error
	daoUrl, err = url.Parse(authUrl)
	if err != nil {
		return false, fmt.Errorf("DAOKEEPER_URL Parse error %s,%v", daokeeper, err)
	}
	daoUrl.Path = "/v1/ngrokd/auth"
	return nil
}

func Auth(authToken string) (bool, error) {
	resp, err := http.PostForm(u.String(), url.Values{"user": {authToken}})
	if err != nil {
		return false, fmt.Errorf("Request daokeeper error %s,%v", u.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, fmt.Errorf("Response daokeeper code %d,%v", resp.StatusCode)
	}
	return true, nil
}
