package state

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

// GetSetter is used to get and set the status of the devices in a room.
type GetSetter struct {
	Logger         *zap.Logger
	DriverRegistry avcontrol.DriverRegistry
}
