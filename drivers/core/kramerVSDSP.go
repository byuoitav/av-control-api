package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer-driver"
	"go.uber.org/zap"
)

type KramerVSDSPDriver struct {
	Log *zap.Logger
}

func (k *KramerVSDSPDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KramerVSDSPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return kramer.NewVideoSwitcherDsp(addr, kramer.WithLoggerVSDSP(k.Log.Sugar())), nil
}
