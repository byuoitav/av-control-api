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
func (gs *GetSetter) Get(ctx context.Context, room api.Room) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
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
		&getVolumes{
			Environment: gs.Environment,
			Logger:      log.With(zap.String("evaluator", "getVolume")),
		},
		&getMutes{
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
