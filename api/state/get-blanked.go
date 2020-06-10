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

func (g *getBlanked) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	for _, dev := range room.Devices {
		log := g.Logger.With(zap.String("device", string(dev.ID)))

		url, order, err := getCommand(dev, "GetBlanked", g.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			log.Warn("unable to get command", zap.String("command", "GetBlanked"), zap.Error(err))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: err.Error(),
			})

			continue
		}

		params := map[string]string{
			"address": dev.Address,
		}

		url, err = fillURL(url, params)
		if err != nil {
			log.Warn("unable to fill url", zap.Error(err))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "blanked",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Warn("unable to build request", zap.Error(err))
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

		log.Info("Successfully built action")
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

	log := g.Logger.With(zap.String("device", string(aResp.Action.ID)))

	handleErr := func(err error) {
		log.Warn("error handling response", zap.Error(err))
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
		handleErr(fmt.Errorf("%v response from driver: %s", aResp.StatusCode, aResp.Body))
		return
	}

	var state blanked
	if err := json.Unmarshal(aResp.Body, &state); err != nil {
		handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, aResp.Body))
		return
	}

	log.Info("Successfully got blanked state", zap.Boolp("blanked", state.Blanked))

	aResp.Updates <- DeviceStateUpdate{
		ID: aResp.Action.ID,
		DeviceState: api.DeviceState{
			Blanked: state.Blanked,
		},
	}
}
