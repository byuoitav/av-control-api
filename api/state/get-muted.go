package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/graph"
	"go.uber.org/zap"
)

type getMuted struct {
	Logger      api.Logger
	Environment string
}

func (g *getMuted) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	gr := graph.NewGraph(room.Devices, "audio")

	responses := make(chan actionResponse)

	for _, dev := range room.Devices {
		path := graph.PathToEnd(gr, dev.ID)
		if len(path) == 0 {
			url, order, err := getCommand(dev, "GetMutedByBlock", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetMutedByBlock"), zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Error: err.Error(),
				})
				continue
			default:
				params := map[string]string{
					"address": dev.Address,
					"input":   string(dev.ID),
				}

				url, err = fillURL(url, params)
				if err != nil {
					g.Logger.Warn("unable to fill url", zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					g.Logger.Warn("unable to build request", zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
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

				g.Logger.Info("Successfully built action", zap.Any("device", path[0].Src.Device.ID))

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
				continue
			}

			// it should always be by block
			url, order, err = getCommand(dev, "GetMuted", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetMuted"), zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "muted",
					Error: err.Error(),
				})

				continue
			default:
				params := map[string]string{
					"address": dev.Address,
				}

				url, err = fillURL(url, params)
				if err != nil {
					g.Logger.Warn("unable to fill url", zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					g.Logger.Warn("unable to build request", zap.Any("device", path[0].Src.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
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

				g.Logger.Info("Successfully built action", zap.Any("device", path[0].Src.Device.ID))

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
			}
		} else {
			endDev := path[len(path)-1].Dst
			url, order, err := getCommand(*endDev.Device, "GetMutedByBlock", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetMutedByBlock"), zap.Any("device", endDev.Device.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "muted",
					Error: err.Error(),
				})

				continue
			}

			for _, port := range endDev.Ports {
				if !port.Endpoints.Contains(dev.ID) {
					continue
				}

				params := map[string]string{
					"address": endDev.Address,
					"input":   port.Name,
				}

				url, err = fillURL(url, params)
				if err != nil {
					g.Logger.Warn("unable to fill url", zap.Any("device", endDev.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					g.Logger.Warn("unable to build request", zap.Any("device", endDev.Device.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
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

				g.Logger.Info("Successfully built action", zap.Any("device", endDev.Device.ID))

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
				continue
			}
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

type muted struct {
	Muted bool `json:"muted"`
}

func (g *getMuted) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			g.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "muted",
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

		g.Logger.Info("Successfully got muted state", zap.Any("device", resp.Action.ID), zap.Boolp("muted", &state.Muted))
		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Muted: &state.Muted,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
