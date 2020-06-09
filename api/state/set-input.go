package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/graph"
	"go.uber.org/zap"
	gonum "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
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

	g := graph.NewGraph(room.Devices, "video")
	tmpOutputs := graph.Leaves(g)
	var outputs []graph.Node

	for _, node := range tmpOutputs {
		for _, dev := range devices {
			if node.Device.ID == dev.ID {
				outputs = append(outputs, node)
				break
			}
		}
	}

	t := graph.Transpose(g)
	inputs := graph.Leaves(t)

	paths := path.DijkstraAllPaths(t)

	for _, device := range devices {
		var actsForOutput []action
		var errsForOutput []api.DeviceStateError

		path := graph.PathFromTo(t, &paths, device.ID, *stateReq.Devices[device.ID].Input.Video)
		if len(path) == 0 {
			s.Logger.Warn(fmt.Sprintf("unable to find a path from %s to %v", device.ID, *stateReq.Devices[device.ID].Input), zap.Any("device", device.ID))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    device.ID,
				Field: "setInput",
				Error: fmt.Sprintf("no path from %s to %v", device.ID, *stateReq.Devices[device.ID].Input),
			})

			continue
		}

		acts, errs := s.generateActionsForPath(ctx, path, s.Environment, responses, stateReq)
		actsForOutput = append(actsForOutput, acts...)
		errsForOutput = append(errsForOutput, errs...)

		if len(errsForOutput) == 0 {
			resp.ExpectedUpdates++
			resp.Actions = append(resp.Actions, actsForOutput...)
		}

		resp.Errors = append(resp.Errors, errsForOutput...)
	}

	if resp.ExpectedUpdates == 0 {
		return resp
	}

	resp.Actions = uniqueActions(resp.Actions)

	if len(resp.Actions) > 0 {
		go s.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates, t, &paths, outputs, inputs)
	}

	return resp
}

func (s *setInput) generateActionsForPath(ctx context.Context, path graph.Path, env string, resps chan actionResponse, stateReq api.StateRequest) ([]action, []api.DeviceStateError) {
	var acts []action
	var errs []api.DeviceStateError
	var transmitterAddr string
	for i := range path {
		if strings.Contains(string(path[i].Src.Device.ID), "TX") {
			transmitterAddr = path[i].Src.Device.Address
		}
	}

	for i := range path {
		url, order, err := getCommand(*path[i].Src.Device, "SetInput", s.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			s.Logger.Warn("unable to get command", zap.String("command", "SetInput"), zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "setInput",
				Error: err.Error(),
			})

			return acts, errs

		default:
			params := map[string]string{
				"address":     path[i].Src.Address,
				"port":        path[i].SrcPort.Name,
				"transmitter": transmitterAddr,
			}

			url, err = fillURL(url, params)
			if err != nil {
				s.Logger.Warn("unable to fill url", zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
				errs = append(errs, api.DeviceStateError{
					ID:    path[i].Src.Device.ID,
					Field: "setInput",
					Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
				})

				return acts, errs
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				s.Logger.Warn("unable to build request", zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
				errs = append(errs, api.DeviceStateError{
					ID:    path[i].Src.Device.ID,
					Field: "setInput",
					Error: fmt.Sprintf("unable to build http request: %s", err),
				})

				return acts, errs
			}

			act := action{
				ID:       path[i].Src.Device.ID,
				Req:      req,
				Order:    order,
				Response: resps,
			}

			s.Logger.Info("Successfully built action", zap.Any("device", path[i].Src.Device.ID))

			acts = append(acts, act)

			continue
		}

		url, order, err = getCommand(*path[i].Src.Device, "SetInputByOutput", env)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
			continue
		case err != nil:
			s.Logger.Warn("unable to get command", zap.String("command", "SetInputByOutput"), zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "setInput",
				Error: err.Error(),
			})

			return acts, errs
		}

		params := map[string]string{
			"address": path[i].Src.Address,
			"input":   path[i].SrcPort.Name,
			"output":  path[i-1].DstPort.Name,
		}

		url, err = fillURL(url, params)
		if err != nil {
			s.Logger.Warn("unable to fill url", zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "setInput",
				Error: fmt.Sprintf("%s (url after fill: %s)", err, url),
			})

			return acts, errs
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			s.Logger.Warn("unable to build request", zap.Any("device", path[i].Src.Device.ID), zap.Error(err))
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "setInput",
				Error: fmt.Sprintf("unable to build http request: %s", err),
			})

			return acts, errs
		}

		act := action{
			ID:       path[i].Src.Device.ID,
			Req:      req,
			Order:    order,
			Response: resps,
		}

		s.Logger.Info("Successfully built action", zap.Any("device", path[i].Src.Device.ID))

		acts = append(acts, act)
	}

	return acts, errs
}

type i struct {
	Audio            *string  `json:"audio"`
	Video            *string  `json:"video"`
	CanSetSeparately *bool    `json:"canSetSeparately"`
	AvailableInputs  []string `json:"availableInputs"`
}

func (s *setInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int, t *simple.DirectedGraph, paths *path.AllShortest, devices, inputs []graph.Node) {
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
	status := make(map[api.DeviceID][]i)

	for _, resp := range resps {
		handleErr := func(err error) {
			s.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "setInput",
				Error: err.Error(),
			}

			resp.Updates <- DeviceStateUpdate{}
		}

		if resp.Error != nil {
			handleErr(fmt.Errorf("unable to make http request: %w", resp.Error))
			continue
		}

		if resp.StatusCode/100 != 2 {
			handleErr(fmt.Errorf("%v response from driver: %s", resp.StatusCode, resp.Body))
			continue
		}

		var state i
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		status[resp.Action.ID] = append(status[resp.Action.ID], state)
	}

	for _, device := range devices {
		deepest := device

		var prevEdge graph.Edge
		var prevState i
		search := traverse.DepthFirst{
			Visit: func(node gonum.Node) {
				deepest = node.(graph.Node)
			},
			Traverse: func(edge gonum.Edge) bool {
				e := edge.(graph.Edge)

				states := status[e.Src.Device.ID]

				if _, ok := e.Src.Type.Commands["SetInput"]; ok {
					for _, state := range states {
						if *state.Video == "" {
							continue
						}

						inputStr := *state.Video
						split := strings.Split(inputStr, ":")
						if len(split) > 1 {
							inputStr = split[1]
						}

						if e.SrcPort.Name == inputStr {
							prevState = state
							prevEdge = e
							return true
						}

						// well we took off outgoing so idk how this needs to change yet
						// if len(e.Src.Device.Ports.Outgoing()) == 1 {
						// 	prevState = state
						// 	prevEdge = e
						// 	return true
						// }
					}
					return false
				}

				if _, ok := e.Src.Type.Commands["SetInputByOutput"]; ok {
					for _, state := range states {
						if *state.Video == "" {
							continue
						}

						inputStr := *state.Video
						split := strings.Split(inputStr, ":")
						if len(split) > 1 {
							inputStr = split[1]
						}
						if prevEdge == (graph.Edge{}) {
							if inputStr == e.SrcPort.Name {
								prevState = state
								prevEdge = e
								return true
							}
						} else {
							if len(split) > 1 {
								if split[1] == prevEdge.DstPort.Name && e.SrcPort.Name == split[0] {
									prevState = state
									prevEdge = e
									return true
								}
							} else {
								if inputStr == e.SrcPort.Name {
									prevState = state
									prevEdge = e
									return true
								}
							}
						}
					}
					return false
				}

				// well we took off outgoing so idk how this needs to change yet
				// if len(e.Src.Device.Ports.Outgoing()) == 1 {
				// 	prevEdge = e
				// 	return true
				// }

				if *prevState.Video == e.Dst.Address {
					prevEdge = e
					return true
				}

				return false
			},
		}

		search.Walk(t, device, func(node gonum.Node) bool {
			return t.From(node.ID()).Len() == 0
		})
		// validate that the deepest is an input
		valid := false
		for _, input := range inputs {
			if deepest.Device.ID == input.Device.ID {
				valid = true
				break
			}
		}

		if valid {
			in := api.Input{
				Video: &deepest.Device.ID,
			}
			// yeah i feel like i can't just always give it resps[0]...
			s.Logger.Info("successfully set input", zap.Any("device", resps[0].Action.ID), zap.Any("input", in))
			resps[0].Updates <- DeviceStateUpdate{
				ID: device.Device.ID,
				DeviceState: api.DeviceState{
					Input: &in,
				},
			}
		} else {
			states := status[deepest.Device.ID]
			s.Logger.Warn("unable to traverse back to valid input")
			resps[0].Errors <- api.DeviceStateError{
				ID:    device.Device.ID,
				Field: "input",
				// I tried doing %s with *states[0].Input and that worked sometimes
				// but it fails when the first request fails, so maybe something else...
				Error: fmt.Sprintf("unable to traverse input back to a valid input. only got to %s|%+v", deepest.Device.ID, states),
			}
		}
	}
}
