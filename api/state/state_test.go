package state

import (
	"net/http"
	"net/url"

	"github.com/byuoitav/av-control-api/api"
)

type stateTest struct {
	name          string
	room          string
	deviceService api.DeviceService
	env           string
	resp          generateActionsResponse
}

func urlParse(rawurl string) *url.URL {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err.Error())
	}

	return url
}

func newRequest(method string, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}

	return req
}

func boolP(b bool) *bool {
	return &b
}

func intP(i int) *int {
	return &i
}

func deviceID(s string) *api.DeviceID {
	v := api.DeviceID(s)
	return &v
}
