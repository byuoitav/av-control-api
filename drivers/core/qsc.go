package core

import (
	"context"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/qsc-driver"
)

func ParseQSCConfig(config map[string]interface{}) error {
	return nil
}

func GetQSCDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &qsc.DSP{
		Address: addr,
	}, nil
}
