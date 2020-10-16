package adcp

import (
	"context"
	"errors"

	"github.com/byuoitav/adcp-driver"
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type SonyADCPDriver struct {
	Log      *zap.Logger
	username string
}

func (s *SonyADCPDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}

func (s *SonyADCPDriver) ParseConfig(config map[string]interface{}) error {
	var ok bool

	if s.username, ok = config["username"].(string); !ok {
		return errors.New("username must be a string")
	}

	return nil
}
