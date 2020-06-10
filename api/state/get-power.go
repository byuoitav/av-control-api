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

type getPower struct {
	Logger      api.Logger
	Environment string
}

func (g *getPower) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	for _, dev := range room.Devices {
		url, order, err := getCommand(dev, "GetPower", g.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			g.Logger.Warn("unable to get command", zap.String("command", "GetPower"), zap.Any("device", dev.ID), zap.Error(err))
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
			g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
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
			g.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
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

		g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))
		go g.handleResponse(act.Response)

		resp.Actions = append(resp.Actions, act)
		resp.ExpectedUpdates++
	}

	return resp
}

type poweredOn struct {
	PoweredOn bool `json:"poweredOn"`

	// TODO we want to get rid of this once drivers support it
	Power string `json:"power"`
}

func (g *getPower) handleResponse(respChan chan actionResponse) {
	aResp := <-respChan
	close(respChan)

	handleErr := func(err error) {
		g.Logger.Warn("error handling response", zap.Any("device", aResp.Action.ID), zap.Error(err))
		aResp.Errors <- api.DeviceStateError{
			ID:    aResp.Action.ID,
			Field: "poweredOn",
			Error: err.Error(),
		}

		aResp.Updates <- DeviceStateUpdate{}
	}

	if aResp.Error != nil {
		handleErr(fmt.Errorf("unable to make http request: %w", aResp.Error))
		return
	}

	var state poweredOn
	if err := json.Unmarshal(aResp.Body, &state); err != nil {
		handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body))
		return
	}

	if state.Power != "" {
		state.PoweredOn = state.Power == "on"
	}

	g.Logger.Info("Successfully got power state", zap.Any("device", aResp.Action.ID), zap.Boolp("blanked", &state.PoweredOn))
	aResp.Updates <- DeviceStateUpdate{
		ID: aResp.Action.ID,
		DeviceState: api.DeviceState{
			PoweredOn: &state.PoweredOn,
		},
	}
}
