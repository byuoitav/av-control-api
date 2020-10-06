package core

import (
	"context"
	"time"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/nec-driver"
)

func GetNECDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return nec.NewProjector(addr, nec.WithDelay(300*time.Second)), nil
}
