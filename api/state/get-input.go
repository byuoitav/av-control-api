package state

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/byuoitav/av-control-api/api"
)

type getInput struct{}

func (i *getInput) GenerateActions(ctx context.Context, room []api.Device, env string) generateActionsResponse {
	var resp generateActionsResponse

	// g := graph.NewGraph(room, "video")
	// paths := path.DijkstraAllPaths(g)

	responses := make(chan actionResponse)

	for _, dev := range room {
		url, order, err := getCommand(dev, "GetInput", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "input",
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
					Field: "input",
					Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
				})

				continue
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				resp.Errors = append(resp.Errors, api.DeviceStateError{
					ID:    dev.ID,
					Field: "input",
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

		// for now - lets just assume that we will need the input for _all_ outgoing ports?
		url, order, err = getCommand(dev, "GetInputByOutput", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Error: err.Error(),
			})

			continue
		default:
			for _, port := range dev.Ports {
				if !port.Outgoing {
					continue
				}

				params := map[string]string{
					"address": dev.Address,
					"output":  port.Name,
				}

				url, err = fillURL(url, params)
				if err != nil {
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "input",
						Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
					})

					continue
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					resp.Errors = append(resp.Errors, api.DeviceStateError{
						ID:    dev.ID,
						Field: "input",
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
				continue
			}
		}
	}

	if resp.ExpectedUpdates == 0 {
		return generateActionsResponse{}
	}

	if len(resp.Actions) > 0 {
		go i.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates)
	}

	return resp
}

type input struct {
	Input *string `json:"input"`
}

func (i *getInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
	if expectedResps == 0 {
		return
	}

	received := 0
	var resps []actionResponse

	for resp := range respChan {
		received++
		resps = append(resps, resp)

		if received == expectedResps {
			break
		}
	}

	close(respChan)

	for _, resp := range resps {
		resp.Updates <- DeviceStateUpdate{}
	}
}
