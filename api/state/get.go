package state

import (
	"context"
	"errors"
	"sort"

	"github.com/byuoitav/av-control-api/api"
)

var (
	ErrNoStateGettable = errors.New("can't get the state of any devices in this room")
)

// GetDevices .
// TODO combine identical actions (?)
func GetDevices(ctx context.Context, room []api.Device, env string) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
	}

	var resp generateActionsResponse

	for i := range statusEvaluators {
		r := statusEvaluators[i].GenerateActions(ctx, room, env)
		resp.Actions = append(resp.Actions, r.Actions...)
		resp.ExpectedUpdates = append(resp.ExpectedUpdates, r.ExpectedUpdates...)
		stateResp.Errors = append(stateResp.Errors, r.Errors...)
	}

	// split the commands into their lists by id
	actsByID := make(map[api.DeviceID][]action)
	for i := range resp.Actions {
		actsByID[resp.Actions[i].ID] = append(actsByID[resp.Actions[i].ID], resp.Actions[i])
	}

	// order every id's commands
	for id := range actsByID {
		sort.Slice(actsByID[id], func(i, j int) bool {
			switch {
			case actsByID[id][i].Order == nil && actsByID[id][j].Order == nil:
				return false
			case actsByID[id][i].Order == nil:
				return false
			case actsByID[id][j].Order == nil:
				return true
			default:
				return *actsByID[id][i].Order < *actsByID[id][j].Order
			}
		})
	}

	// execute commands
	updates := make(chan DeviceStateUpdate)
	errors := make(chan api.DeviceStateError)

	for id := range actsByID {
		executeActions(actsByID[id], updates, errors)
	}

	updatesReceived := 0

	if len(resp.ExpectedUpdates) == 0 {
		return stateResp, ErrNoStateGettable
	}

	for {
		// TODO ctx.Done()?
		select {
		case update := <-updates:
			updatesReceived++
			if len(update.ID) == 0 {
				break
			}

			curState := stateResp.Devices[update.ID]

			if update.Blanked != nil {
				curState.Blanked = update.Blanked
			}

			stateResp.Devices[update.ID] = curState

			// TODO add the other ones
		case err := <-errors:
			stateResp.Errors = append(stateResp.Errors, err)
		}

		if len(resp.ExpectedUpdates) == updatesReceived {
			break
		}
	}

	return stateResp, nil
}
