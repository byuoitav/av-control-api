package state

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type getPower struct{}

func (g *getPower) GenerateActions(ctx context.Context, room []api.Device, env string) generateActionsResponse {
	var resp generateActionsResponse

	for _, dev := range room {
		url, order, err := getCommand(dev, "GetPower", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "poweredOn",
				Error: err.Error(),
			})

			continue
		}

		params := map[string]string{
			"address": dev.Address,
		}
		url, err = fillURL(url, params)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "poweredOn",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		// build http request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "poweredOn",
				Error: fmt.Sprintf("unable to build http request: %s", err),
			})

			continue
		}
		act := action{
			ID:       dev.ID,
			Req:      req,
			Order:    order,
			Response: make(chan actionResponse),
		}
		go g.handleResponse(act.Response)

		resp.Actions = append(resp.Actions, act)
		resp.ExpectedUpdates++
	}

	return resp
}

type poweredOn struct {
	PoweredOn bool `json:"poweredOn"`
}

func (g *getPower) handleResponse(respChan chan actionResponse) {
	aResp := <-respChan
	close(respChan)

	handleErr := func(err error) {
		aResp.Errors <- api.DeviceStateError{
			ID:    aResp.Action.ID,
			Field: "poweredOn",
			Error: err.Error(),
		}

		aResp.Updates <- OutputStateUpdate{}
	}

	if aResp.Error != nil {
		handleErr(fmt.Errorf("unable to make http request: %w", aResp.Error))
		return
	}

	var state poweredOn

	// since we get back {"power": "standby"} we're doing this for now
	if string(aResp.Body) == "{\"power\":\"on\"}" {
		state.PoweredOn = true
	} else if string(aResp.Body) == "{\"power\":\"standby\"}" {
		state.PoweredOn = false
	} else {
		handleErr(fmt.Errorf("unexpected response from driver:\n%s", aResp.Body))
	}

	// I guess ideally we'd do this but not for now...
	// if err := json.Unmarshal(aResp.Body, &state); err != nil {
	// 	handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body))
	// 	return
	// }

	aResp.Updates <- OutputStateUpdate{
		ID: aResp.Action.ID,
		OutputState: api.OutputState{
			PoweredOn: &state.PoweredOn,
		},
	}
}
