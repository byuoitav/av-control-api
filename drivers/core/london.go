package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/london-driver"
)

func GetLondonDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return london.New(addr, london.WithLogger(drivers.Log.Named(addr))), nil
}
