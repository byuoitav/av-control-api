package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type getMuted struct {
	Logger      api.Logger
	Environment string
}

var (
	ErrMultipleIncoming = errors.New("multiple incoming ports on device")
)

func (g *getMuted) checkCommand(dev api.Device, responses chan actionResponse, room api.Room, incoming bool) (action, bool, error) {
	numIncoming := 0
	for _, port := range dev.Ports {
		if port.Incoming && strings.Contains(port.Type, "audio") {
			numIncoming++
		}
	}

	if numIncoming > 1 {
		return action{}, true, ErrMultipleIncoming
	}

	for _, port := range dev.Ports {
		if port.Incoming && strings.Contains(port.Type, "audio") {
			incoming = true
			for _, d := range room.Devices {
				if port.Endpoints.Contains(d.ID) {
					url, order, err := getCommand(d, "GetMutedByBlock", g.Environment)
					switch {
					case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
						act, incoming, err := g.checkCommand(d, responses, room, incoming)
						if err != nil {
							return action{}, incoming, err
						}
						return act, incoming, nil

					case err != nil:
						g.Logger.Warn("unable to get command", zap.String("command", "GetMutedByBlock"), zap.Any("device", d.ID), zap.Error(err))
						return action{}, incoming, err

					default:
						params := map[string]string{
							"address": d.Address,
							"block":   port.Name,
						}

						url, err = fillURL(url, params)
						if err != nil {
							g.Logger.Warn("unable to fill url", zap.Any("device", d.ID), zap.Error(err))
							return action{}, incoming, fmt.Errorf("%s (url after fill: %s)", err, url)
						}

						req, err := http.NewRequest(http.MethodGet, url, nil)
						if err != nil {
							g.Logger.Warn("unable to build request", zap.Any("device", d.ID), zap.Error(err))
							return action{}, incoming, fmt.Errorf("unable to build http request: %s", err)
						}

						g.Logger.Info("Successfully built action", zap.Any("device", d.ID))

						act := action{
							ID:       dev.ID,
							Req:      req,
							Order:    order,
							Response: responses,
						}
						return act, incoming, nil
					}
				}
			}
		}
	}

	return action{}, incoming, errCommandNotFound
}

func (g *getMuted) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	responses := make(chan actionResponse)

	for _, dev := range room.Devices {
		act, incoming, err := g.checkCommand(dev, responses, room, false)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound), errors.Is(err, ErrMultipleIncoming):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "volume",
				Error: err.Error(),
			})

			continue
		default:
			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
			continue
		}
		if !incoming || errors.Is(err, ErrMultipleIncoming) {
			url, order, err := getCommand(dev, "GetMuted", g.Environment)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
				g.Logger.Warn("unable to get command", zap.String("command", "GetMuted"), zap.Any("device", dev.ID), zap.Error(err))
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
					g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "muted",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					g.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
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

				fmt.Printf("action for %s: %s\n", act.ID, act.Req.URL)

				g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
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
