package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/keydigital-driver"
	"go.uber.org/zap"
)

type KeyDigitalDriver struct {
	Log *zap.Logger
}

func (k *KeyDigitalDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KeyDigitalDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return keydigital.CreateVideoSwitcher(ctx, addr, k.Log.Sugar())
}
