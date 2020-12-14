package driverstest

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
)

type Driver struct {
	Devices map[string]avcontrol.Device
}

func (d *Driver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (d *Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return d.Devices[addr], nil
}
