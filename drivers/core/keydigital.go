package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/keydigital-driver"
)

func ParseKeyDigitalConfig(config map[string]interface{}) error {
	return nil
}

func GetKeyDigitalDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return keydigital.CreateVideoSwitcher(ctx, addr, drivers.Log.Named(addr))
}
