package state

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type setMuted struct {
	Logger      api.Logger
	Environment string
}

func (s *setMuted) GenerateActions(ctx context.Context, room api.Room, stateReq api.StateRequest) generatedActions {
	var resp generatedActions

	var devices []api.Device
	for k, v := range stateReq.Devices {
		if v.Mutes != nil {
			for i := range room.Devices {
				if room.Devices[i].ID == k {
					devices = append(devices, room.Devices[i])
					break
				}
			}
		}
	}

	if len(devices) == 0 {
		return resp
	}
	responses := make(chan actionResponse)

	for _, dev := range devices {
		url, order, err := getCommand(dev, "SetMute", s.Environment)
		if err != nil {
			s.Logger.Warn("unable to get command", zap.String("command", "SetMute"), zap.Any("device", dev.ID), zap.Error(err))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "setMute",
				Error: err.Error(),
			})
			continue
		}
		for k, v := range stateReq.Devices[dev.ID].Mutes {
			params := map[string]string{
				"address": dev.Address,
				"block":   k,
				"muted":   strconv.FormatBool(v),
			}

			url, err = fillURL(url, params)
			if err != nil {
				s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setMute",
					Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
				})

				continue
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setMute",
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
	}

	if resp.ExpectedUpdates == 0 {
		return resp
	}

	if len(resp.Actions) > 0 {
		go s.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates)
	}

	return resp
}

type mute struct {
	Muted bool `json:"muted"`
}

func (s *setMuted) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			s.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setMuted",
				Error: err.Error(),
			}

			resp.Updates <- DeviceStateUpdate{}
		}
		received++
		if resp.Error != nil {
			handleErr(fmt.Errorf("unable to make http request: %w", resp.Error))
			continue
		}

		if resp.StatusCode/100 != 2 {
			handleErr(fmt.Errorf("%v response from driver: %s", resp.StatusCode, resp.Body))
			continue
		}

		var state muted
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		s.Logger.Info("Successfully set muted state", zap.Any("device", resp.Action.ID), zap.Bool("muted", state.Muted))
		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Mutes: map[string]bool{
					string(resp.Action.ID): state.Muted,
				},
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
