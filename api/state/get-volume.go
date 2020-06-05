package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/graph"
)

type getVolume struct{}

func (g *getVolume) GenerateActions(ctx context.Context, room []api.Device, env string) generateActionsResponse {
	var resp generateActionsResponse

	gr := graph.NewGraph(room, "audio")

	responses := make(chan actionResponse)

	for _, dev := range room {
		if dev.ID == "ITB-1108B-DSP1" || dev.ID == "ITB-1108B-DSP2" || dev.ID == "ITB-1108A-DSP1" {
			continue
		}
		path := graph.PathToEnd(gr, dev.ID)
		if len(path) == 0 {
			url, order, err := getCommand(dev, "GetVolumeByBlock", env)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			case err != nil:
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
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "volume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
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

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
				continue
			}

			url, order, err = getCommand(dev, "GetVolume", env)
			switch {
			case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
				continue
			case err != nil:
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
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "volume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
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
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "volume",
					Error: err.Error(),
				})

				continue
			}

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
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "volume",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
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

				resp.Actions = append(resp.Actions, act)
				resp.ExpectedUpdates++
				continue
			}
		}
	}

	if resp.ExpectedUpdates == 0 {
		return generateActionsResponse{}
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
	var resps []actionResponse

	for resp := range respChan {
		received++
		resps = append(resps, resp)

		var state volume
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "volume",
				Error: fmt.Sprintf("unable to parse response from driver: %v. response:\n%s", err, resp.Body),
			}

			resp.Updates <- OutputStateUpdate{}
			continue
		}

		resp.Updates <- OutputStateUpdate{
			ID: resp.Action.ID,
			OutputState: api.OutputState{
				Volume: &state.Volume,
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
