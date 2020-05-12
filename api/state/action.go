package state

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api"
)

type action struct {
	ID    api.DeviceID
	Req   *http.Request
	Order *int

	Response chan actionResponse
}

type actionResponse struct {
	Action *action
	Error  error

	StatusCode int
	Header     http.Header
	Body       []byte

	Errors  chan api.DeviceStateError
	Updates chan DeviceStateUpdate
}

func executeActions(actions []action, updates chan DeviceStateUpdate, errors chan api.DeviceStateError) {
	for i := range actions {
		go func(action action) {
			aResp := actionResponse{
				Action:  &action,
				Errors:  errors,
				Updates: updates,
			}

			fmt.Printf("sending request to %s\n", action.Req.URL.String())

			resp, err := http.DefaultClient.Do(action.Req)
			if err != nil {
				aResp.Error = err
				action.Response <- aResp
				return
			}
			defer resp.Body.Close()

			aResp.StatusCode = resp.StatusCode
			aResp.Header = resp.Header

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				aResp.Error = err
				action.Response <- aResp
				return
			}

			aResp.Body = body
			action.Response <- aResp
		}(actions[i])

		time.Sleep(50 * time.Millisecond)
	}
}
