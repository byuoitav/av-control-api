package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
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
func (g *getInput) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions
	responses := make(chan actionResponse)
	gr := graph.NewGraph(room.Devices, "video")

	outputs := graph.Leaves(gr)

	t := graph.Transpose(gr)
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

			acts, errs := g.generateActionsForPath(ctx, path, responses)
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
		go g.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates, t, &paths, outputs, inputs, room)
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
	Audio            *string        `json:"audio"`
	Video            *string        `json:"video"`
	CanSetSeparately *bool          `json:"canSetSeparately"`
	AvailableInputs  []api.DeviceID `json:"availableInputs"`
}

type respInput struct {
	Input *string `json:"input"`
}

func (g *getInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int, t *simple.DirectedGraph, paths *path.AllShortest, outputs, inputs []graph.Node, room api.Room) {
	if expectedResps == 0 {
		return
	}

	var resps []actionResponse
	var received int

	separatable := make(map[api.DeviceID]*bool)
	availableInputs := make(map[api.DeviceID][]api.DeviceID)

	for _, dev := range room.Devices {
		path := graph.PathToEnd(t, dev.ID)

		for i := range path {
			if len(path[i].Dst.Device.Ports.Incoming()) == 0 {
				availableInputs[dev.ID] = append(availableInputs[dev.ID], path[i].Dst.Device.ID)
			}
			if _, ok := path[i].Src.Device.Type.Commands["GetVideoInput"]; ok {
				separatable[dev.ID] = boolP(true)
			}
			if _, ok := path[i].Src.Device.Type.Commands["GetVideoInputForOutput"]; ok {
				separatable[dev.ID] = boolP(true)
			}
			if _, ok := path[i].Src.Device.Type.Commands["GetAudioInput"]; ok {
				separatable[dev.ID] = boolP(true)
			}
			if _, ok := path[i].Src.Device.Type.Commands["GetAudioInputForOutput"]; ok {
				separatable[dev.ID] = boolP(true)
			}
		}
		if _, ok := separatable[dev.ID]; !ok {
			separatable[dev.ID] = boolP(false)
		}

	}

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
		var tmpInput respInput
		var switcherBackup map[string]string
		if err := json.Unmarshal(resp.Body, &tmpInput); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}
		if tmpInput == (respInput{}) {
			if err := json.Unmarshal(resp.Body, &switcherBackup); err != nil {
				handleErr(fmt.Errorf("unable to parse response from driver for switcher: %w. response:\n%s", err, resp.Body))
			}
		}

		switch {
		case strings.Contains(resp.Action.Req.URL.String(), "GetAVInputForOutput"):
			if tmpInput != (respInput{}) {
				state.Video = tmpInput.Input
				state.Audio = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Video = &tmpInput
					state.Audio = &tmpInput
				}
			}

		case strings.Contains(resp.Action.Req.URL.String(), "GetAVInput"):
			if tmpInput != (respInput{}) {
				state.Video = tmpInput.Input
				state.Audio = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Video = &tmpInput
					state.Audio = &tmpInput
				}
			}
		case strings.Contains(resp.Action.Req.URL.String(), "GetStream"):
			if tmpInput != (respInput{}) {
				state.Video = tmpInput.Input
				state.Audio = tmpInput.Input
			}

		case strings.Contains(resp.Action.Req.URL.String(), "GetVideoInputForOutput"):
			if tmpInput != (respInput{}) {
				state.Video = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Video = &tmpInput
				}
			}

		case strings.Contains(resp.Action.Req.URL.String(), "GetVideoInput"):
			if tmpInput != (respInput{}) {
				state.Video = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Video = &tmpInput
				}
			}

		case strings.Contains(resp.Action.Req.URL.String(), "GetAudioInputForOutput"):
			if tmpInput != (respInput{}) {
				state.Audio = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Audio = &tmpInput
				}
			}

		case strings.Contains(resp.Action.Req.URL.String(), "GetAudioInput"):
			if tmpInput != (respInput{}) {
				state.Audio = tmpInput.Input
			} else {
				for k, v := range switcherBackup {
					tmpInput := v + ":" + k
					state.Audio = &tmpInput
				}
			}

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
		deepestVideo := g.getDeepest(output, "video", t, status)
		deepestAudio := g.getDeepest(output, "audio", t, status)

		// validate that the deepest is an input
		valid := false
		for _, input := range inputs {
			if deepestVideo.Device.ID == input.Device.ID {
				valid = true
				break
			}
		}
		if valid {
			// it works with these print statements but not without(?)
			// for k, v := range availableInputs[output.Device.ID] {
			// 	fmt.Printf("Before: %v, %v\n", k, v)
			// }
			sort.Slice(availableInputs[output.Device.ID], func(i, j int) bool { return i < j })
			// for k, v := range availableInputs[output.Device.ID] {
			// 	fmt.Printf("after: %v, %v\n", k, v)
			// }
			i := api.Input{
				Video:            &deepestVideo.Device.ID,
				Audio:            &deepestAudio.Device.ID,
				CanSetSeparately: separatable[output.Device.ID],
				AvailableInputs:  availableInputs[output.Device.ID],
			}
			resps[0].Updates <- DeviceStateUpdate{
				ID: output.Device.ID,
				DeviceState: api.DeviceState{
					Input: &i,
				},
			}

		} else {
			states := status[deepestVideo.Device.ID]
			g.Logger.Warn("unable to traverse input tree back to a valid video input", zap.Any("device", deepestVideo.Device.ID))

			resps[0].Errors <- api.DeviceStateError{
				ID:    output.Device.ID,
				Field: "input",
				Error: fmt.Sprintf("unable to traverse input tree back to a valid video input. only got to %s|%+v", deepestVideo.Device.ID, states),
			}

			// I don't know if it will get both of these but we'll see
			states = status[deepestAudio.Device.ID]
			g.Logger.Warn("unable to traverse input tree back to a valid audio input", zap.Any("device", deepestAudio.Device.ID))

			resps[0].Errors <- api.DeviceStateError{
				ID:    output.Device.ID,
				Field: "input",
				Error: fmt.Sprintf("unable to traverse input tree back to a valid audio input. only got to %s|%+v", deepestAudio.Device.ID, states),
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

func (g *getInput) getDeepest(output graph.Node, inputType string, t *simple.DirectedGraph, status map[api.DeviceID][]input) graph.Node {
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

			if inputType == "video" {
				if prevEdge != (graph.Edge{}) && prevState.Video != nil {
					if len(states) == 0 && prevState.Video == &e.Dst.Device.Address {
						prevEdge = e
						return true
					}
				}

				fmt.Println()
				if _, ok := e.Src.Type.Commands["GetAVInput"]; ok {
					for _, state := range states {
						if state.Video == nil {
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
						return false
					}
				}

				if _, ok := e.Src.Type.Commands["GetAVInputForOutput"]; ok {
					for _, state := range states {
						if state.Video == nil {
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

				if _, ok := e.Src.Type.Commands["GetVideoInput"]; ok {
					for _, state := range states {
						if state.Video == nil {
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

						// not sure if we want this outgoing thing

						// if len(e.Src.Device.Ports.Outgoing()) == 1 {
						// 	prevState = state
						// 	prevEdge = e
						// 	return true
						// }
					}

					return false
				}

				if _, ok := e.Src.Type.Commands["GetVideoInputForOutput"]; ok {
					for _, state := range states {
						if state.Video == nil {
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
			}
			if _, ok := e.Src.Type.Commands["GetAVInput"]; ok {
				for _, state := range states {
					if state.Audio == nil {
						continue
					}

					inputStr := *state.Audio
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

			if _, ok := e.Src.Type.Commands["GetAVInputForOutput"]; ok {
				for _, state := range states {
					if state.Audio == nil {
						continue
					}

					inputStr := *state.Audio
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

			if _, ok := e.Src.Type.Commands["GetAudioInput"]; ok {
				for _, state := range states {
					if state.Audio == nil {
						continue
					}

					inputStr := *state.Audio
					split := strings.Split(inputStr, ":")
					if len(split) > 1 {
						inputStr = split[1]
					}

					if e.SrcPort.Name == inputStr {
						prevEdge = e
						prevState = state
						return true
					}

					// if len(e.Src.Device.Ports.Outgoing()) == 1 {
					// 	prevState = state
					// 	prevEdge = e
					// 	return true
					// }
				}
				return false
			}

			if _, ok := e.Src.Type.Commands["GetAudioInputForOutput"]; ok {
				for _, state := range states {
					if state.Audio == nil {
						continue
					}

					inputStr := *state.Audio
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
			if prevState.Audio == nil {
				return false
			}

			if *prevState.Audio == e.Dst.Address {
				prevEdge = e
				return true
			}

			return false
		},
	}

	search.Walk(t, output, func(node gonum.Node) bool {
		return t.From(node.ID()).Len() == 0
	})

	return deepest
}
