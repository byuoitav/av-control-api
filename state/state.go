package state

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type GetSetter struct {
	Logger         *zap.Logger
	DriverRegistry avcontrol.DriverRegistry
}
