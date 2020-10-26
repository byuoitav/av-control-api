package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/sony/adcp"
	"go.uber.org/zap"
)

type SonyADCPDriver struct {
	Log *zap.Logger
}

func (s *SonyADCPDriver) ParseConfig(config map[string]interface{}) error {
	s.Log.Info("logging something")
	return nil
}

func (s *SonyADCPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}
