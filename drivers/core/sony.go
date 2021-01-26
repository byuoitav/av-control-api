package core

import (
	"context"
	"errors"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/sony/bravia"
	"go.uber.org/zap"
)

type SonyDriver struct {
	PreSharedKey string
	Log          *zap.Logger
}

func (s *SonyDriver) ParseConfig(config map[string]interface{}) error {
	if psk, ok := config["psk"].(string); ok {
		if psk == "" {
			return errors.New("given empty psk")
		}

		s.PreSharedKey = psk
	} else {
		return errors.New("no psk given")
	}

	return nil
}

func (s *SonyDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &bravia.Display{
		Address:      addr,
		PreSharedKey: s.PreSharedKey,
		Log:          s.Log,
		RequestDelay: 250 * time.Millisecond,
	}, nil
}
