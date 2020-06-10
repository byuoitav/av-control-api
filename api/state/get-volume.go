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

type getVolume struct {
	Logger      api.Logger
	Environment string
}

func (g *getVolume) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	gr := graph.NewGraph(room.Devices, "audio")

	responses := make(chan actionResponse)

	for _, dev := range room.Devices {

		path := graph.PathToEnd(gr, dev.ID)
		if len(path) == 0 {
			url, order, err := getCommand(dev, "GetVolumeByBlock", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetVolumeByBlock"), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Error: err.Error(),
				})

				continue
			default:
				params := map[string]string{
					"address": dev.Address,
					"block":   string(dev.ID),
				}

				url, err = fillURL(url, params)
				if err != nil {
					g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
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
				continue
			}

			url, order, err = getCommand(dev, "GetVolume", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetVolume"), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
					Error: err.Error(),
				})

				continue
			default:
				params := map[string]string{
					"address": dev.Address,
				}

				url, err = fillURL(url, params)
				if err != nil {
					g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
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
		} else {
			endDev := path[len(path)-1].Dst
			url, order, err := getCommand(*endDev.Device, "GetVolumeByBlock", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetVolumeByBlock"), zap.Any("device", endDev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
					Error: err.Error(),
				})

				continue
			default:

				for _, port := range endDev.Ports {
					if port.Endpoints.Contains(dev.ID) {
						continue
					}

					//at this point we have the right port

					params := map[string]string{
						"address": endDev.Address,
						"block":   port.Name,
					}

					url, err = fillURL(url, params)
					if err != nil {
						g.Logger.Warn("unable to fill url", zap.Any("device", endDev.ID), zap.Error(err))
						resp.Errors = append(resp.Errors, api.DeviceStateError{
							ID:    dev.ID,
							Field: "volume",
							Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
						})

						continue
					}

					req, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						g.Logger.Warn("unable to build request", zap.Any("device", endDev.ID), zap.Error(err))
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

					g.Logger.Info("Successfully built action", zap.Any("device", endDev.ID))

					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
			}
			url, order, err = getCommand(*endDev.Device, "GetVolume", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetVolume"), zap.Any("device", endDev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
					Error: err.Error(),
				})

				continue
			default:

				for _, port := range endDev.Ports {
					if port.Endpoints.Contains(dev.ID) {
						continue
					}

					//at this point we have the right port

					params := map[string]string{
						"address": endDev.Address,
					}

					url, err = fillURL(url, params)
					if err != nil {
						g.Logger.Warn("unable to fill url", zap.Any("device", endDev.ID), zap.Error(err))
						resp.Errors = append(resp.Errors, api.DeviceStateError{
							ID:    dev.ID,
							Field: "volume",
							Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
						})

						continue
					}

					req, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						g.Logger.Warn("unable to build request", zap.Any("device", endDev.ID), zap.Error(err))
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

					g.Logger.Info("Successfully built action", zap.Any("device", endDev.ID))

					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
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

type volume struct {
	Volume int `json:"volume"`
}

func (g *getVolume) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
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

		g.Logger.Info("Successfully got volume state", zap.Any("device", resp.Action.ID), zap.Int("volume", state.Volume))
		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Volume: &state.Volume,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
