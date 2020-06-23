package state

/*

type getMutes struct {
	Logger      api.Logger
	Environment string
}

var (
	ErrMultipleIncoming = errors.New("multiple incoming ports on device")
)

func (g *getMutes) GenerateActions(ctx context.Context, room api.Room) generatedActions {
	var resp generatedActions

	responses := make(chan actionResponse)

	for _, dev := range room.Devices {
		url, order, err := getCommand(dev, "GetMutes", g.Environment)
		switch {
		case errors.Is(err, errCommandNotFound), errors.Is(err, errCommandEnvNotFound):
		case err != nil:
			g.Logger.Warn("unable to get command", zap.String("command", "GetMutes"), zap.Any("device", dev))
			resp.Errors = append(resp.Errors, api.DeviceStateError{
				ID:    dev.ID,
				Field: "muted",
				Error: err.Error(),
			})

			// continue
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

			g.Logger.Info("Successfully built action", zap.Any("device", dev.ID))

			resp.Actions = append(resp.Actions, act)
			resp.ExpectedUpdates++
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

func (g *getMutes) handleResponses(respChan chan actionResponse, expectedResps, expectedUpdates int) {
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
		g.Logger.Info("Successfully got muted state", zap.Any("device", resp.Action.ID), zap.Any("muted", state.Muted))
		resp.Updates <- DeviceStateUpdate{
			ID: resp.Action.ID,
			DeviceState: api.DeviceState{
				Mutes: map[string]bool{
					string(resp.Action.ID): state.Muted,
				},
			},
		}

		if received == expectedResps {
			break
		}
	}

	close(respChan)
}
*/
