package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/justaddpower-driver"
)

type JAPDriver struct{}

func (j *JAPDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (j *JAPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &justaddpower.JustAddPowerReciever{
		Address: addr,
	}, nil
}
