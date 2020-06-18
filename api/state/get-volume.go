package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type getVolumes struct {
	Logger      api.Logger
	Environment string
}

func (g *getVolumes) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	responses := make(chan actionResponse)

	for _, dev := range room.Devices {
		url, order, err := getCommand(dev, "GetVolumes", g.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			g.Logger.Warn("unable to get command", zap.String("command", "GetVolumes"), zap.Any("device", dev))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "volume",
				Error: err.Error(),
			})

		default:
			params := map[string]string{
				"address": dev.Address,
			}

			url, err = fillURL(url, params)
			if err != nil {
				g.Logger.Warn("uanble to fill url", zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
					Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
				})

				continue
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				g.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
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

			g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))

			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
		}
	}

	if resp.ExpectedUpdates == 0 {
		return generatedActions{}
	}

	if len(resp.Actions) > 0 {
		go g.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates)
	}

	return resp
}

type volume struct {
	Volume int `json:"volume"`
}

func (g *getVolumes) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			g.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "volume",
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

		var state volume
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		g.Logger.Info("Successfully got volume state", zap.Any("device", resp.Action.ID), zap.Any("volume", state.Volume))
		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Volumes: map[string]int{
					string(resp.Action.ID): state.Volume,
				},
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
