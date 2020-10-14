package core

import (
	"context"

	"github.com/byuoitav/adcp-driver"
	"github.com/byuoitav/av-control-api/drivers"
)

func ParseADCPConfig(config map[string]interface{}) error {
	return nil
}

func GetADCPDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}
