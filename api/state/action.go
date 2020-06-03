package state

import (
	"context"
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

func (a *action) Equal(b *action) bool {
	switch {
	case a == b:
		return true
	case a == nil && b != nil:
		return false
	case a != nil && b == nil:
		return false
	case a.ID != b.ID:
		return false
	case a.Order == nil && b.Order != nil:
		return false
	case a.Order != nil && b.Order == nil:
		return false
	case a.Order != nil && *a.Order != *b.Order:
		return false
	case a.Response != b.Response:
		return false
	case a.Req == nil && b.Req != nil:
		return false
	case a.Req != nil && b.Req == nil:
		return false
	case a.Req != nil && a.Req.Method != b.Req.Method:
		return false
	case a.Req != nil && a.Req.URL.String() != b.Req.URL.String():
		return false
	}

	return true
}

type actionResponse struct {
	Action *action
	Error  error

	StatusCode int
	Header     http.Header
	Body       []byte

	Errors  chan api.DeviceStateError
	Updates chan OutputStateUpdate
}

func executeActions(ctx context.Context, actions []action, updates chan OutputStateUpdate, errors chan api.DeviceStateError) {
	for i := range actions {
		go func(action action) {
			aResp := actionResponse{
				Action:  &action,
				Errors:  errors,
				Updates: updates,
			}

			fmt.Printf("sending request to %s\n", action.Req.URL.String())

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			action.Req = action.Req.WithContext(ctx)

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

func uniqueActions(actions []action) []action {
	var unique []action

	for _, action := range actions {
		add := true
		for _, u := range unique {
			if action.Equal(&u) {
				add = false
				break
			}
		}

		if add {
			unique = append(unique, action)
		}
	}

	return unique
}
