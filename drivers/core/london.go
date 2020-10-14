package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/london-driver"
)

type LondonDriver struct{}

func (l *LondonDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (l *LondonDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return london.New(addr, london.WithLogger(drivers.Log.Named(addr))), nil
}
