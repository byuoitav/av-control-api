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

type setVolume struct {
	Logger      api.Logger
	Environment string
}

func (s *setVolume) GenerateActions(ctx context.Context, room api.Room, env string, stateReq api.StateRequest) generatedActions {
	var resp generatedActions
	gr := graph.NewGraph(room.Devices, "audio")

	responses := make(chan actionResponse)

	var devices []api.Device

	for k, v := range stateReq.Devices {
		if v.Volume != nil {
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
		if len(path) == 0 {
			url, order, err := getCommand(dev, "SetVolumeByBlock", env)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetVolumeByBlock"), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Error: err.Error(),
				})

				continue

			default:
				params := map[string]string{
					"address": dev.Address,
					"input":   string(dev.ID),
					"volume":  strconv.Itoa(*stateReq.Devices[dev.ID].Volume),
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
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
			// it should always be by block so we should remove this later
			url, order, err = getCommand(dev, "SetVolume", env)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetVolume"), zap.Any("device", dev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setVolume",
					Error: err.Error(),
				})

				continue
			default:
				params := map[string]string{
					"address": dev.Address,
					"level":   strconv.Itoa(*stateReq.Devices[dev.ID].Volume),
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
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
			url, order, err := getCommand(*endDev.Device, "GetVolumeByBlock", env)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				s.Logger.Warn("unable to get command", zap.String("command", "SetVolumeByBlock"), zap.Any("device", endDev.ID), zap.Error(err))
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "setVolume",
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
					"volume":  strconv.Itoa(*stateReq.Devices[dev.ID].Volume),
				}

				url, err = fillURL(url, params)
				if err != nil {
					s.Logger.Warn("unable to fill url", zap.Any("device", endDev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					s.Logger.Warn("unable to build request", zap.Any("device", endDev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setVolume",
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
				continue
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

type v struct {
	Volume int `json:"volume"`
}

func (s *setVolume) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			s.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setVolume",
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

		var state v
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		s.Logger.Info("Successfully set blanked state", zap.Any("device", resp.Action.ID), zap.Int("volume", state.Volume))
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
