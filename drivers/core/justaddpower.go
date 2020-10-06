package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/justaddpower-driver"
)

func GetJustAddPowerDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &justaddpower.JustAddPowerReciever{
		Address: addr,
	}, nil
}
