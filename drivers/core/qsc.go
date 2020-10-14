package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/qsc-driver"
)

type QSCDriver struct{}

func (q *QSCDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (q *QSCDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &qsc.DSP{
		Address: addr,
	}, nil
}
