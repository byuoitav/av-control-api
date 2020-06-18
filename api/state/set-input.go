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

type setInput struct {
	Logger      api.Logger
	Environment string
}

func (s *setInput) GenerateActions(ctx context.Context, room api.Room, stateReq api.StateRequest) generatedActions {
	var resp generatedActions

	var devices []api.Device
	for k, v := range stateReq.Devices {
		if v.Input != nil {
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
		for _, input := range stateReq.Devices[dev.ID].Input {
			if input.AudioVideo != nil {
				act, err := s.checkCommand(dev, "SetAudioVideoInput", responses, *input.AudioVideo)
				switch {
				case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
					// maybe still return an error here?
				case err != nil:
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setInput",
						Error: err.Error(),
					})

					continue
				default:
					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
			}
			if input.Audio != nil {
				act, err := s.checkCommand(dev, "SetAudioInput", responses, *input.Audio)
				switch {
				case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
					// maybe still return an error here?
				case err != nil:
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setInput",
						Error: err.Error(),
					})

					continue
				default:
					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
			}
			if input.Video != nil {
				act, err := s.checkCommand(dev, "SetVideoInput", responses, *input.Video)
				switch {
				case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
					// maybe still return an error here?
				case err != nil:
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "setInput",
						Error: err.Error(),
					})

					continue
				default:
					resp.Actions = append(resp.Actions, act)
					resp.ExpectedUpdates++
				}
			}
		}
	}

	if resp.ExpectedUpdates == 0 {
		return resp
	}

	// probably unnecessary
	resp.Actions = uniqueActions(resp.Actions)

	if len(resp.Actions) > 0 {
		go s.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates)
	}

	return resp
}

func (s *setInput) checkCommand(dev api.Device, cmd string, resps chan (actionResponse), block string) (action, error) {
	url, order, err := getCommand(dev, cmd, s.Environment)
	switch {
	case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		return action{}, err
	case err != nil:
		s.Logger.Warn("unable to get command", zap.String("command", cmd), zap.Any("device", dev.ID), zap.Error(err))
		return action{}, err
	default:
		params := map[string]string{
			"address": dev.Address,
			// TODO: figure out block and transmitter for JAPs
			"block":       block,
			"transmitter": "alsoidk",
		}

		url, err = fillURL(url, params)
		if err != nil {
			s.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
			return action{}, fmt.Errorf("%s (url after fill: %s)", err, url)
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			s.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
			return action{}, fmt.Errorf("unable to build http request: %s", err)
		}

		act := action{
			ID:       dev.ID,
			Req:      req,
			Order:    order,
			Response: resps,
		}

		s.Logger.Info("Successfully built action", zap.Any("device", dev.ID))

		return act, nil
	}
}

// we'll receive
// "inputs": {
// 	"": {
// 		"audioVideo": "hdmi2"
// 	}
// }

type in struct {
	AudioVideo string
	Video      string
	Audio      string
}

func (s *setInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			s.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setInput",
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

		var state in
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		s.Logger.Info("Successfully set input state", zap.Any("device", resp.Action.ID), zap.Any("input", &state))

		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Input: map[string]api.Input{
					string(resp.Action.ID): api.Input{
						AudioVideo: &state.AudioVideo,
						Video:      &state.Video,
						Audio:      &state.Audio,
					},
				},
			},
		}
		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
