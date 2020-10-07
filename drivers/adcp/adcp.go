package adcp

import (
	"context"
	"errors"

	"github.com/byuoitav/adcp-driver"
	avcontrol "github.com/byuoitav/av-control-api"
)

type SonyADCPDriver struct {
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
