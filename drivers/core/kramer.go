package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer-driver"
)

type KramerDriver struct{}

func (k *KramerDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KramerDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return kramer.NewVideoSwitcher(addr), nil
}
