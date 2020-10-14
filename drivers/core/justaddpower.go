package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/justaddpower-driver"
)

func ParseJAPConfig(config map[string]interface{}) error {
	return nil
}

func GetJAPDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &justaddpower.JustAddPowerReciever{
		Address: addr,
	}, nil
}
