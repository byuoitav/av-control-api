package state

import (
	"context"
	"sort"

	"github.com/byuoitav/av-control-api/api"
)

func GetDevices(ctx context.Context, room []api.Device, env string) (api.StateResponse, error) {
	// devState := make(map[api.DeviceID]api.DeviceState)

	var acts []action
	var errs []api.DeviceStateError

	for i := range statusEvaluators {
		actions, errors, _ := statusEvaluators[i].GenerateActions(ctx, room, env)
		acts = append(acts, actions...)
		errs = append(errs, errors...)
	}

	// TODO combine identical actions

	// split the commands into their lists by id
	actsByID := make(map[api.DeviceID][]action)
	for i := range acts {
		actsByID[acts[i].ID] = append(actsByID[acts[i].ID], acts[i])
	}

	// order every id's commands
	for id, _ := range actsByID {
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

	return api.StateResponse{}, nil
}
