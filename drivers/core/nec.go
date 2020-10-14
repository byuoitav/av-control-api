package core

import (
	"context"
	"time"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/nec-driver"
)

func ParseNECConfig(config map[string]interface{}) error {
	return nil
}

func GetNECDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return nec.NewProjector(addr, nec.WithDelay(300*time.Second)), nil
}
