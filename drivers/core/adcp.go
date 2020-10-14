package core

import (
	"context"

	"github.com/byuoitav/adcp-driver"
	avcontrol "github.com/byuoitav/av-control-api"
)

type SonyADCPDriver struct{}

func (s *SonyADCPDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (s *SonyADCPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}
