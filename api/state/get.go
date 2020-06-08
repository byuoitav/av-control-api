package state

import (
	"context"
	"errors"
	"sort"

	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

var (
	ErrNoStateGettable = errors.New("can't get the state of any devices in this room")
)

// Get .
func (gs *GetSetter) Get(ctx context.Context, room []api.Device) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		OutputGroups: make(map[api.DeviceID]api.OutputGroupState),
	}

	id := api.RequestID(ctx)
	log := gs.Logger.With(zap.String("requestID", id))

	evaluators := []statusEvaluator{
		&getPower{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getPower")),
		},
		&getBlanked{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getBlanked")),
		},
		&getInput{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getInput")),
		},
		&getVolume{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getVolume")),
		},
		&getMuted{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getMuted")),
		},
	}

	var actions []action
	var expectedUpdates int

	log.Info("Generating actions")
	for i := range evaluators {
		resp := evaluators[i].GenerateActions(ctx, room)
		actions = append(actions, resp.Actions...)
		stateResp.Errors = append(stateResp.Errors, resp.Errors...)
		expectedUpdates += resp.ExpectedUpdates
	}

	log.Info("Done generating actions", zap.Int("actions", len(actions)), zap.Int("errors", len(stateResp.Errors)), zap.Int("expectedUpdates", expectedUpdates))

	if expectedUpdates == 0 {
		return stateResp, ErrNoStateGettable
	}

	// split the commands into their lists by id
	actsByID := make(map[api.DeviceID][]action)
	for i := range actions {
		actsByID[actions[i].ID] = append(actsByID[actions[i].ID], actions[i])
	}

	log.Info("Ordering commands for each device")

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

	log.Info("Done ordering commands")

	// execute commands
	updates := make(chan OutputStateUpdate)
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
