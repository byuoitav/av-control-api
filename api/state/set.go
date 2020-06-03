package state

import (
	"context"
	"errors"
	"sort"

	"github.com/byuoitav/av-control-api/api"
)

var (
	ErrNoPowerSettable = errors.New("can't set power of given device")
	ErrNoStateSettable = errors.New("can't set the state of any devices in this room")
)

func SetDevices(ctx context.Context, req api.StateRequest, room []api.Device, env string) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		OutputGroups: make(map[api.DeviceID]api.OutputGroupState),
	}

	var actions []action
	var expectedUpdates int

	for i := range commandEvaluators {
		resp := commandEvaluators[i].GenerateActions(ctx, room, env, req)
		actions = append(actions, resp.Actions...)
		stateResp.Errors = append(stateResp.Errors, resp.Errors...)
		expectedUpdates += resp.ExpectedUpdates
	}

	if expectedUpdates == 0 {
		if len(stateResp.Errors) == 0 {
			return api.StateResponse{}, ErrNoStateSettable
		}
		return stateResp, nil
	}

	// split the commands into their lists by id
	actsByID := make(map[api.DeviceID][]action)
	for i := range actions {
		actsByID[actions[i].ID] = append(actsByID[actions[i].ID], actions[i])
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
	updates := make(chan OutputStateUpdate)
	errors := make(chan api.DeviceStateError)

	for id := range actsByID {
		executeActions(ctx, actsByID[id], updates, errors)
	}

	updatesReceived := 0

	for {
		select {
		case update := <-updates:

			updatesReceived++

			if len(update.ID) == 0 {
				break
			}

			curState := stateResp.OutputGroups[update.ID]

			if update.PoweredOn != nil {
				curState.PoweredOn = update.PoweredOn
			}

			if update.Input != nil {
				curState.Input = update.Input
			}

			if update.Blanked != nil {
				curState.Blanked = update.Blanked
			}

			if update.Volume != nil {
				curState.Volume = update.Volume
			}

			if update.Muted != nil {
				curState.Muted = update.Muted
			}

			// this is on a group state so don't need it?
			// if update.Outputs != nil {
			// 	curState.Outputs = update.Outputs
			// }

			stateResp.OutputGroups[update.ID] = curState
		case err := <-errors:
			stateResp.Errors = append(stateResp.Errors, err)
		}

		if updatesReceived == expectedUpdates {
			break
		}
	}

	return stateResp, nil
}
