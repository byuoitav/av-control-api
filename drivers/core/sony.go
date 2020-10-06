package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/sonyrest-driver"
)

func GetSonyDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &sonyrest.TV{
		Address: addr,
		// PSK: psk,
	}, nil
}
