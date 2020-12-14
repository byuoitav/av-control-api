package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer-driver"
	"go.uber.org/zap"
)

type KramerDSPDriver struct {
	Log *zap.Logger
}

func (k *KramerDSPDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KramerDSPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return kramer.NewDsp(addr, kramer.WithLoggerDSP(k.Log.Sugar())), nil
}
