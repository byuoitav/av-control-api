package state

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

var (
	commandEvaluators = []commandEvaluator{
		// &setMuted{},
		// &setPower{},
		// &setVolume{},
		&setBlanked{},
		// &setInput{},
	}
)

type DeviceStateUpdate struct {
	ID api.DeviceID
	api.DeviceState
}

type generatedActions struct {
	Actions         []action
	Errors          []api.DeviceStateError
	ExpectedUpdates int
}

type statusEvaluator interface {
	GenerateActions(ctx context.Context, room api.Room) generatedActions
}

type commandEvaluator interface {
	GenerateActions(ctx context.Context, room api.Room, state api.StateRequest) generatedActions
}
