package core

import (
	"context"

	"github.com/byuoitav/atlona-driver"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/wspool"
)

// func GetADCPDevice(ctx context.Context, addr string) (drivers.Device, error) {
// 	return &adcp.Projector{
// 		Address: addr,
// 	}, nil
// }

func NewDevice(ctx context.Context, addr, username, password string, log wspool.Logger) (drivers.Device, error) {
	return atlona.CreateVideoSwitcher(ctx, addr, username, password, log)
}
