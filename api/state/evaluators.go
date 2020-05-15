package state

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

var (
	statusEvaluators = []statusEvaluator{
		&getBlanked{},
		&getInput{},
		&getPower{},
		&getVolume{},
		&getMuted{},
	}
)

type DeviceStateUpdate struct {
	ID api.DeviceID
	api.DeviceState
}

type generateActionsResponse struct {
	Actions         []action
	Errors          []api.DeviceStateError
	ExpectedUpdates int
}

type statusEvaluator interface {
	GenerateActions(ctx context.Context, room []api.Device, env string) generateActionsResponse
}
