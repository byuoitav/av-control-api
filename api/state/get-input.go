package state

/*

type getInput struct {
	Logger      api.Logger
	Environment string
}

type getInputData struct {
	Command string
}

func (g *getInput) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions
	responses := make(chan actionResponse)

	for _, dev := range room.Devices {
		act, err := g.checkCommand(dev, "GetAudioVideoInputs", responses)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
			continue
		}

		act, err = g.checkCommand(dev, "GetVideoInputs", responses)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
			// I don't think we need to continue here...
		}

		act, err = g.checkCommand(dev, "GetAudioInputs", responses)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "input",
				Error: err.Error(),
			})
		default:
			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
		}
	}

	if len(resp.Actions) > 0 {
		go g.handleResponses(responses, len(resp.Actions), resp.ExpectedUpdates, room)
	}

	return resp
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
			Data: getInputData{
				Command: cmd,
			},
		}

		g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))
		return act, nil
	}
}

type inputs struct {
	Inputs map[string]input `json:"inputs"`
}

type input struct {
	AudioVideo *string `json:"audiovideo,omitempty"`
	Audio      *string `json:"audio,omitempty"`
	Video      *string `json:"video,omitempty"`
}

func (g *getInput) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int, room api.Room) {
	if expectedResps == 0 {
		return
	}

	received := 0

	for resp := range respChan {
		handleErr := func(err error) {
			g.Logger.Warn("error handling response", zap.Any("device", resp.Action.ID), zap.Error(err))
			resp.Errors <- api.DeviceStateError{
				ID:    resp.Action.ID,
				Field: "input",
				Error: err.Error(),
			}

			//idk if I have to do this now
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

		var state inputs
		if err := json.Unmarshal(resp.Body, &state); err != nil {
			handleErr(fmt.Errorf("unable to parse response from driver: %w. response:\n%s", err, resp.Body))
			continue
		}

		g.Logger.Info("Successfully got input state", zap.Any("device", resp.Action.ID), zap.Any("input", &state))
		tmpInput := make(map[string]api.Input)
		for k, v := range state.Inputs {
			tmpInput[k] = api.Input{
				AudioVideo: v.AudioVideo,
				Audio:      v.Audio,
				Video:      v.Video,
			}
		}

		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Input: tmpInput,
			},
		}
		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
*/
