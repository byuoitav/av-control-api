package adcp

import (
	"context"

	"github.com/byuoitav/adcp-driver"
	"github.com/byuoitav/av-control-api/drivers"
)

func NewDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}
