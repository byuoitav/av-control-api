package core

import (
	"context"
	"errors"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/sony/bravia"
	"go.uber.org/zap"
)

type SonyDriver struct {
	PSK string
	Log *zap.Logger
}

func (s *SonyDriver) ParseConfig(config map[string]interface{}) error {
	if psk, ok := config["psk"].(string); ok {
		if psk == "" {
			return errors.New("given empty psk")
		}

		s.PSK = psk
	} else {
		return errors.New("no psk given")
	}

	return nil
}

func (s *SonyDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &bravia.TV{
		Address: addr,
		PSK:     s.PSK,
		Log:     s.Log,
	}, nil
}
