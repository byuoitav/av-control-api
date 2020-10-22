package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer-driver"
	"go.uber.org/zap"
)

type KramerVSDriver struct {
	Log *zap.Logger
}

func (k *KramerVSDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KramerVSDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return kramer.NewVideoSwitcher(addr, kramer.WithLogger4x4(k.Log.Named(addr).Sugar())), nil
}
