package state

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type action struct {
	ID    api.DeviceID
	Req   *http.Request
	Order *int
	Data  interface{}

	Response chan actionResponse
}

// Equal reports if two actions are equal. Equal is defined as the following:
// * a.ID == b.ID
// * a.Order == b.Order
// * a.Req.Method == b.Req.Method
// * a.Req.URL.String() == b.Req.URL.String()
//
// It does not compare the Response channels of the two structs.
func (a action) Equal(b action) bool {
	switch {
	case a.ID != b.ID:
		return false
	case a.Order == nil && b.Order != nil:
		return false
	case a.Order != nil && b.Order == nil:
		return false
	case a.Order != nil && *a.Order != *b.Order:
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
	Updates chan DeviceStateUpdate
}

func (gs *GetSetter) executeActions(ctx context.Context, actions []action, updates chan DeviceStateUpdate, errors chan api.DeviceStateError) {
	for i := range actions {
		go func(action action) {
			aResp := actionResponse{
				Action:  &action,
				Errors:  errors,
				Updates: updates,
			}

			log := gs.Logger.With(zap.Any("device", action.ID), zap.String("method", action.Req.Method), zap.String("url", action.Req.URL.String()))
			log.Info("Sending request")

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			action.Req = action.Req.WithContext(ctx)

			resp, err := http.DefaultClient.Do(action.Req)
			if err != nil {
				log.Warn("unable to make request", zap.Error(err))
				aResp.Error = err
				action.Response <- aResp
				return
			}
			defer resp.Body.Close()

			aResp.StatusCode = resp.StatusCode
			aResp.Header = resp.Header

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Warn("unable to read response", zap.Error(err))
				aResp.Error = err
				action.Response <- aResp
				return
			}

			log.Debug("response", zap.Int("statusCode", resp.StatusCode), zap.ByteString("body", body))

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
			if action.Equal(u) {
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
