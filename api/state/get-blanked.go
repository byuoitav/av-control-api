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

type getBlanked struct {
	Logger      api.Logger
	Environment string
}

func (g *getBlanked) GenerateActions(ctx context.Context, room []api.Device) generatedActions {
	var resp generatedActions

	for _, dev := range room {
		url, order, err := getCommand(dev, "GetBlanked", g.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			g.Logger.Warn("unable to get command", zap.String("command", "GetBlanked"), zap.Any("device", dev.ID), zap.Error(err))
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
			g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
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
			g.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
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

		g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))
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
		g.Logger.Warn("error handling response", zap.Any("device", aResp.Action.ID), zap.Error(err))
		aResp.Errors <- api.DeviceStateError{
			ID:    aResp.Action.ID,
			Field: "blanked",
			Error: err.Error(),
		}

		aResp.Updates <- OutputStateUpdate{}
	}

	if aResp.Error != nil {
		handleErr(fmt.Errorf("unable to make http request: %w", aResp.Error))
		return
	}

	if aResp.StatusCode/100 != 2 {
		handleErr(fmt.Errorf("%v response from driver: %s", aResp.StatusCode, aResp.Body))
		return
	}

	var state blanked
	if err := json.Unmarshal(aResp.Body, &state); err != nil {
		handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body))
		return
	}

	g.Logger.Info("Successfully got blanked state", zap.Any("device", aResp.Action.ID), zap.Boolp("blanked", state.Blanked))
	aResp.Updates <- OutputStateUpdate{
		ID: aResp.Action.ID,
		OutputState: api.OutputState{
			Blanked: state.Blanked,
		},
	}
}
