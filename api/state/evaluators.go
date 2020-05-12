package state

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

var (
	statusEvaluators = []statusEvaluator{
		&getBlanked{},
	}
)

type DeviceStateUpdate struct {
	ID api.DeviceID
	api.DeviceState
}

type statusEvaluator interface {
	GenerateActions(ctx context.Context, room []api.Device, env string) ([]action, []api.DeviceStateError, []DeviceStateUpdate)
}
