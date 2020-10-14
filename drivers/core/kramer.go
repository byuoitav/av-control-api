package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/kramer-driver"
)

func ParseKramerConfig(config map[string]interface{}) error {
	return nil
}

func GetKramerDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return kramer.NewVideoSwitcher(addr), nil
}
