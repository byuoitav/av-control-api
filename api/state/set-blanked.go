package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type setBlanked struct{}

func (s *setBlanked) GenerateActions(ctx context.Context, room []api.Device, env string, stateReq api.StateRequest) generatedActions {
	var resp generatedActions

	responses := make(chan actionResponse)

	var devices []api.Device
	for k, v := range stateReq.OutputGroups {
		if v.Blanked != nil {
			for i := range room {
				if room[i].ID == k {
					devices = append(devices, room[i])
					break
				}
			}
		}
	}

	if len(devices) == 0 {
		return resp
	}

	for _, dev := range devices {
		var cmd string
		if *stateReq.OutputGroups[dev.ID].Blanked == true {
			cmd = "BlankDisplay"
		} else {
			cmd = "UnblankDisplay"
		}

		url, order, err := getCommand(dev, cmd, env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setBlanked",
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
				Field: "setBlanked",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setBlanked",
				Error: fmt.Sprintf("unable to build http request: %s", err),
			})

			continue
		}

		act := action{
			ID:       dev.ID,
			Req:      req,
			Order:    order,
			Response: responses,
		}

		resp.Actions = append(resp.Actions, act)
		resp.ExpectedUpdates++
	}

	if resp.ExpectedUpdates == 0 {
		return generatedActions{}
	}

	if len(resp.Actions) > 0 {
		go s.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates)
	}

	return resp
}

type blank struct {
	Blanked bool `json:"blanked"`
}

func (s *setBlanked) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		received++
		var state blank
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setBlanked",
				Error: fmt.Sprintf("unable to parse response from driver: %v. response:\n%s", err, resp.Body),
			}

			resp.Updates <- OutputStateUpdate{}
			continue
		}

		resp.Updates <- OutputStateUpdate{
			ID: resp.Action.ID,
			OutputState: api.OutputState{
				Blanked: &state.Blanked,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
