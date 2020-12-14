package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/london"
	"go.uber.org/zap"
)

type LondonDriver struct {
	Log *zap.Logger
}

func (l *LondonDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (l *LondonDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return london.New(addr, london.WithLogger(l.Log.Sugar())), nil
}
