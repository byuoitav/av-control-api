package state

import (
	"context"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type setPower struct{}

func (s *setPower) GenerateActions(ctx context.Context, room api.Room, env string, stateReq api.StateRequest) generatedActions {
	var resp generatedActions

	responses := make(chan actionResponse)

	var devices []api.Device
	for k, v := range stateReq.Devices {
		if v.PoweredOn != nil {
			for i := range room.Devices {
				if room.Devices[i].ID == k {
					devices = append(devices, room.Devices[i])
					break
				}
			}
		}
	}

	for _, dev := range devices {
		var cmd string
		if *stateReq.Devices[dev.ID].PoweredOn == true {
			cmd = "PowerOn"
		} else {
			cmd = "Standby"
		}
		url, order, err := getCommand(dev, cmd, env)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setPower",
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
				Field: "setPower",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setPower",
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

type powered struct {
	PoweredOn bool `json:"poweredOn"`
}

func (s *setPower) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		received++
		var state powered

		// since we get back {"power": "standby"} we're doing this for now
		if string(resp.Body) == "{\"power\":\"on\"}" {
			state.PoweredOn = true
		} else if string(resp.Body) == "{\"power\":\"standby\"}" {
			state.PoweredOn = false
		} else {
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setPower",
				Error: fmt.Sprintf("unexpected response from driver:\n%s", resp.Body),
			}

			resp.Updates <- DeviceStateUpdate{}
			continue
		}

		// I guess ideally we'd do this but not for now...
		// if err := json.Unmarshal(aResp.Body, &state); err != nil {
		//  resp.Errors <- api.DeviceStateError{
		// 	ID: resp.Action.ID,
		// 	Field: "setPower",
		// 	Error: fmt.Sprintf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body),
		// }

		// resp.Updates <- OutputStateUpdate{}
		// continue
		// }

		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				PoweredOn: &state.PoweredOn,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
