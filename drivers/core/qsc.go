package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/qsc"
	"go.uber.org/zap"
)

type QSCDriver struct {
	Log *zap.Logger
}

func (q *QSCDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (q *QSCDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return qsc.New(addr, qsc.WithLogger(q.Log.Sugar())), nil
}
