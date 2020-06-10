package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/graph"
	"go.uber.org/zap"
)

type setMuted struct {
	Logger      api.Logger
	Environment string
}

func (s *setMuted) GenerateActions(ctx context.Context, room api.Room, stateReq api.StateRequest) generatedActions {
	var resp generatedActions

	gr := graph.NewGraph(room.Devices, "audio")

	responses := make(chan actionResponse)

	var devices []api.Device
	for k, v := range stateReq.Devices {
		if v.Muted != nil {
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

	for _, dev := range devices {
		path := graph.PathToEnd(gr, dev.ID)
		var cmd string
		if *stateReq.Devices[dev.ID].Muted == true {
			cmd = "Mute"
		} else {
			cmd = "UnMute"
		}

		if len(path) == 0 {
			url, order, err := getCommand(dev, "SetMutedByBlock", s.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetBlanked"), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Error: err.Error(),
				})
				continue
			default:
				params := map[string]string{
					"address": dev.Address,
					"input":   string(dev.ID),
					"muted":   strconv.FormatBool(*stateReq.Devices[dev.ID].Muted),
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
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
				continue
			}

			url, order, err = getCommand(dev, cmd, s.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", cmd), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setMuted",
					Error: err.Error(),
				})

				continue
			default:
				params := map[string]string{
					"address": dev.Address,
					// "muted":   strconv.FormatBool(*stateReq.Devices[dev.ID].Muted),
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
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
		} else {
			endDev := path[len(path)-1].Dst
			url, order, err := getCommand(*endDev.Device, "SetMutedByBlock", s.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetMutedByBlock"), zap.Any("device", endDev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setMuted",
					Error: err.Error(),
				})

				continue

			default:
				for _, port := range endDev.Ports {
					if !port.Endpoints.Contains(dev.ID) {
						continue
					}

					params := map[string]string{
						"address": endDev.Address,
						"input":   port.Name,
						"muted":   strconv.FormatBool(*stateReq.Devices[dev.ID].Muted),
					}

					url, err = fillURL(url, params)
					if err != nil {
						s.Logger.Warn("unable to fill url", zap.Any("device", endDev.ID), zap.Error(err))
						resp.Errors = append(resp.Errors, api.DeviceStateError{
							ID:    dev.ID,
							Field: "setMuted",
							Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
						})

						continue
					}

					req, err := http.NewRequest(http.MethodGet, url, nil)
					if err != nil {
						s.Logger.Warn("unable to build request", zap.Any("device", endDev.ID), zap.Error(err))
						resp.Errors = append(resp.Errors, api.DeviceStateError{
							ID:    dev.ID,
							Field: "setMuted",
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

					s.Logger.Info("Successfully built action", zap.Any("device", endDev.ID))

					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
			}

			url, order, err = getCommand(*endDev.Device, cmd, s.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetMutedByBlock"), zap.Any("device", endDev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setMuted",
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
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", endDev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", endDev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setMuted",
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

				s.Logger.Info("Successfully built action", zap.Any("device", endDev.ID))

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
			}
		}
	}

	if resp.ExpectedUpdates == 0 {
		return generatedActions{}
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
				Muted: &state.Muted,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
