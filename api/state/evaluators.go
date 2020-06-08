package state

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

var (
	commandEvaluators = []commandEvaluator{
		// &setMuted{},
		&setPower{},
		// &setVolume{},
		&setBlanked{},
		// &setInput{},
	}
)

type OutputStateUpdate struct {
	ID api.DeviceID
	api.OutputState
}

type generatedActions struct {
	Actions         []action
	Errors          []api.DeviceStateError
	ExpectedUpdates int
}

type statusEvaluator interface {
	GenerateActions(ctx context.Context, room []api.Device) generatedActions
}

type commandEvaluator interface {
	GenerateActions(ctx context.Context, room []api.Device, env string, state api.StateRequest) generatedActions
}
