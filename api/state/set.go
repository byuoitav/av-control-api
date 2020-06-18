package state

import (
	"context"
	"errors"
	"sort"

	"github.com/byuoitav/av-control-api/api"
)

var (
	// ErrNoPowerSettable = errors.New("can't set power of given device")
	ErrNoStateSettable = errors.New("nothing to do for the given request and room")
)

func (gs *GetSetter) Set(ctx context.Context, room api.Room, req api.StateRequest) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
	}

	var actions []action
	var expectedUpdates int

	for i := range commandEvaluators {
		resp := commandEvaluators[i].GenerateActions(ctx, room, req)
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
	updates := make(chan DeviceStateUpdate)
	errors := make(chan api.DeviceStateError)

	for id := range actsByID {
		gs.executeActions(ctx, actsByID[id], updates, errors)
	}

	updatesReceived := 0

	for {
		select {
		case update := <-updates:

			updatesReceived++

			if len(update.ID) == 0 {
				break
			}

			curState := stateResp.Devices[update.ID]

			if update.PoweredOn != nil {
				curState.PoweredOn = update.PoweredOn
			}

			if update.Input != nil {
				for k, v := range update.Input {
					tmpInput := curState.Input[k]
					if v.AudioVideo != nil {
						tmpInput.AudioVideo = v.AudioVideo
					}
					if v.Video != nil {
						tmpInput.Video = v.Video
					}
					if v.Audio != nil {
						tmpInput.Audio = v.Audio
					}
					curState.Input[k] = tmpInput
				}
			}

			if update.Blanked != nil {
				curState.Blanked = update.Blanked
			}

			if update.Volumes != nil {
				for k, v := range update.Volumes {
					curState.Volumes[k] = v
				}
			}

			if update.Mutes != nil {
				for k, v := range update.Mutes {
					curState.Mutes[k] = v
				}
			}

			stateResp.Devices[update.ID] = curState
		case err := <-errors:
			stateResp.Errors = append(stateResp.Errors, err)
		}

		if updatesReceived == expectedUpdates {
			break
		}
	}

	return stateResp, nil
}
