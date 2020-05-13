package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type getBlanked struct{}

func (g *getBlanked) GenerateActions(ctx context.Context, room []api.Device, env string) generateActionsResponse {
	var resp generateActionsResponse

	// just doing basic get blanked for now
	for _, dev := range room {
		url, order, err := getCommand(dev, "GetBlanked", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: err.Error(),
			})

			continue
		}

		// replace values
		params := map[string]string{
			"address": dev.Address,
		}
		url, err = fillURL(url, params)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		// build http request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
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

type blanked struct {
	Blanked *bool `json:"blanked"`
}

func (g *getBlanked) handleResponse(respChan chan actionResponse) {
	aResp := <-respChan
	close(respChan)

	handleErr := func(err error) {
		aResp.Errors <- api.DeviceStateError{
			ID:    aResp.Action.ID,
			Field: "blanked",
			Error: err.Error(),
		}

		aResp.Updates <- DeviceStateUpdate{}
	}

	if aResp.Error != nil {
		handleErr(fmt.Errorf("unable to make http request: %w", aResp.Error))
		return
	}

	if aResp.StatusCode/100 != 2 {
		handleErr(fmt.Errorf("%v response from driver: %s", aResp.Error, aResp.Body))
		return
	}

	var state blanked
	if err := json.Unmarshal(aResp.Body, &state); err != nil {
		handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body))
		return
	}

	aResp.Updates <- DeviceStateUpdate{
		ID: aResp.Action.ID,
		DeviceState: api.DeviceState{
			Blanked: state.Blanked,
		},
	}
}
