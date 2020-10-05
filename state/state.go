package state

import (
	"github.com/byuoitav/av-control-api/drivers"
	"go.uber.org/zap"
)

type GetSetter struct {
	Logger  *zap.Logger
	Drivers drivers.Drivers
}
