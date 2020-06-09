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

type getInput struct {
	Logger      api.Logger
	Environment string
}

// GenerateActions makes an assumption that GetInput and GetInputByBlock will not ever be on the same device
func (i *getInput) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions
	responses := make(chan actionResponse)
	g := graph.NewGraph(room.Devices, "video")

	outputs := graph.Leaves(g)

	t := graph.Transpose(g)
	inputs := graph.Leaves(t)

	paths := path.DijkstraAllPaths(t)

	for _, output := range outputs {
		var actsForOutput []action
		var errsForOutput []api.DeviceStateError

		for _, input := range inputs {
			path := graph.PathFromTo(t, &paths, output.Device.ID, input.Device.ID)
			if len(path) == 0 {
				continue
			}

			acts, errs := i.generateActionsForPath(ctx, path, responses)
			actsForOutput = append(actsForOutput, acts...)
			errsForOutput = append(errsForOutput, errs...)
		}

		if len(errsForOutput) == 0 {
			resp.ExpectedUpdates++
			resp.Actions = append(resp.Actions, actsForOutput...)
		}

		resp.Errors = append(resp.Errors, errsForOutput...)
	}

	if resp.ExpectedUpdates == 0 {
		return generatedActions{}
	}

	// combine identical actions
	resp.Actions = uniqueActions(resp.Actions)

	if len(resp.Actions) > 0 {
		go i.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates, t, &paths, outputs, inputs)
	}

	return resp
}

func (g *getInput) generateActionsForPath(ctx context.Context, path graph.Path, resps chan actionResponse) ([]action, []api.DeviceStateError) {
	var acts []action
	var errs []api.DeviceStateError
	for i := range path {
		act, err := g.checkCommand(*path[i].Src.Device, "GetAVInput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}

		act, err = g.checkCommand(*path[i].Src.Device, "GetVideoInput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}

		act, err = g.checkCommand(*path[i].Src.Device, "GetAudioInput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}

		// also need to pass in output
		act, err = g.checkCommand(*path[i].Src.Device, "GetAVInputByOutput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}

		// also need to pass in output
		act, err = g.checkCommand(*path[i].Src.Device, "GetVideoInputByOutput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}

		// also need to pass in output
		act, err = g.checkCommand(*path[i].Src.Device, "GetAudioInputByOutput", resps)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			errs = append(errs, api.DeviceStateError{
				ID:    path[i].Src.Device.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			acts = append(acts, act)
		}
	}

	return acts, errs
}

type input struct {
	Audio            *string  `json:"audio"`
	Video            *string  `json:"video"`
	CanSetSeparately *bool    `json:"canSetSeparately"`
	AvailableInputs  []string `json:"availableInputs"`
}

func (g *getInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int, t *simple.DirectedGraph, paths *path.AllShortest, outputs, inputs []graph.Node) {
	if expectedResps == 0 {
		return
	}

	var resps []actionResponse
	var received int

	for resp := range respChan {
		received++
		resps = append(resps, resp)

		if received == expectedResps {
			break
		}
	}
	close(respChan)
	status := make(map[api.DeviceID][]input)
	var emptyChecker []api.DeviceID
	for _, resp := range resps {
		handleErr := func(err error) {
			g.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "input",
				Error: err.Error(),
			}
		}

		if resp.Error != nil {
			handleErr(fmt.Errorf("unable to make http request: %w", resp.Error))
			continue
		}

		if resp.StatusCode/100 != 2 {
			handleErr(fmt.Errorf("%v response from driver: %s", resp.StatusCode, resp.Body))
			continue
		}

		var state input
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		fmt.Printf("%s input: %s\n", resp.Action.ID, resp.Body)

		// If all devices are off and we can't get any inputs for them we just need to return out
		if string(resp.Body) == "{}" {
			emptyChecker = append(emptyChecker, resp.Action.ID)
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "input",
				Error: fmt.Sprintf("unable to get input for %s (probably powered off)", resp.Action.ID),
			}
			resp.Updates <- DeviceStateUpdate{}
			continue
		}

		status[resp.Action.ID] = append(status[resp.Action.ID], state)
	}

	if len(emptyChecker) == len(outputs) {
		return
	}

	// now calculate the state of the outputs
	for _, output := range outputs {
		skip := false
		for i := range emptyChecker {
			if output.Device.ID == emptyChecker[i] {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		deepest := output

		var prevEdge graph.Edge
		var prevState input
		search := traverse.DepthFirst{
			Visit: func(node gonum.Node) {
				deepest = node.(graph.Node)
			},
			Traverse: func(edge gonum.Edge) bool {
				e := edge.(graph.Edge)

				states := status[e.Src.Device.ID]

				if prevEdge != (graph.Edge{}) && *prevState.Video != "" {
					if len(states) == 0 && *prevState.Video == e.Dst.Device.Address {
						prevEdge = e
						return true
					}
				}

				if _, ok := e.Src.Type.Commands["GetInput"]; ok {
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
							prevEdge = e
							prevState = state
							return true
						}

						if len(e.Src.Device.Ports.Outgoing()) == 1 {
							prevState = state
							prevEdge = e
							return true
						}
					}

					return false
				}

				if _, ok := e.Src.Type.Commands["GetInputByOutput"]; ok {
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

				if len(e.Src.Device.Ports.Outgoing()) == 1 {
					prevEdge = e
					return true
				}

				// dannyrandall: TODO idk how to handle prevState.Video being nil, i just put something
				if prevState.Video == nil {
					return false
				}

				if *prevState.Video == e.Dst.Address {
					prevEdge = e
					return true
				}

				return false
			},
		}

		search.Walk(t, output, func(node gonum.Node) bool {
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
			i := api.Input{
				Video: &deepest.Device.ID,
			}
			resps[0].Updates <- DeviceStateUpdate{
				ID: output.Device.ID,
				DeviceState: api.DeviceState{
					Input: &i,
				},
			}

		} else {
			states := status[deepest.Device.ID]
			g.Logger.Warn("unable to traverse input tree back to a valid input", zap.Any("device", deepest.Device.ID))

			resps[0].Errors <- api.DeviceStateError{
				ID:    output.Device.ID,
				Field: "input",
				Error: fmt.Sprintf("unable to traverse input tree back to a valid input. only got to %s|%+v", deepest.Device.ID, states),
			}
			resps[0].Updates <- DeviceStateUpdate{}
		}
	}
}

func (g *getInput) checkCommand(dev api.Device, cmd string, resps chan actionResponse, port ...string) (action, error) {
	url, order, err := getCommand(dev, cmd, g.Environment)
	switch {
	case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		return action{}, err
	case err != nil:
		g.Logger.Warn("unable to get command", zap.String("command", cmd), zap.Any("device", dev.ID), zap.Error(err))
		return action{}, err
	default:
		params := make(map[string]string)
		params["address"] = dev.Address
		if len(port) > 0 {
			// this might actually be output, idk
			params["port"] = port[0]
		}

		url, err = fillURL(url, params)
		if err != nil {
			g.Logger.Warn("unable to fill url", zap.Any("device", dev.ID), zap.Error(err))
			return action{}, fmt.Errorf("%s (url after fill: %s)", err, url)
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			g.Logger.Warn("unable to build request", zap.Any("device", dev.ID), zap.Error(err))
			return action{}, fmt.Errorf("unable to build http request: %s", err)
		}

		act := action{
			ID:       dev.ID,
			Req:      req,
			Order:    order,
			Response: resps,
		}

		g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))
		return act, nil
	}
}
