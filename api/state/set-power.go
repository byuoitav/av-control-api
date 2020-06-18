package state

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type setPower struct {
	Logger      api.Logger
	Environment string
}

func (s *setPower) GenerateActions(ctx context.Context, room api.Room, stateReq api.StateRequest) generatedActions {
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
		url, order, err := getCommand(dev, "SetPower", s.Environment)
		if err != nil {
			s.Logger.Warn("unable to get command", zap.String("command", "SetPower"), zap.Any("device", dev.ID), zap.Error(err))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setPower",
				Error: err.Error(),
			})

			continue
		}

		params := map[string]string{
			"address": dev.Address,
			"power":   strconv.FormatBool(*stateReq.Devices[dev.ID].PoweredOn),
		}
		url, err = fillURL(url, params)
		if err != nil {
			s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setPower",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
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

		s.Logger.Info("Successfully built action", zap.Any("device", dev.ID))

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
		handleErr := func(err error) {
			s.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setPower",
				Error: err.Error(),
			}

			resp.Updates <- DeviceStateUpdate{}
		}
		received++

		if resp.Error != nil {
			handleErr(fmt.Errorf("unable t o make http requeset: %w", resp.Error))
			continue
		}

		if resp.StatusCode/100 != 2 {
			handleErr(fmt.Errorf("%v response from driver: %s", resp.StatusCode, resp.Body))
			continue
		}
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
		// if err := json.Unmarshal(resp.Body, &state); err != nil {
		// 	handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
		// 	continue
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
