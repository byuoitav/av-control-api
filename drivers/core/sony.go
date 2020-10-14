package core

import (
	"context"
	"errors"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/sonyrest-driver"
)

type SonyDriver struct {
	PSK string
}

func (s *SonyDriver) ParseConfig(config map[string]interface{}) error {
	if psk, ok := config["psk"].(string); ok {
		if psk == "" {
			return errors.New("given empty psk")
		}

		s.PSK = psk
	}

	return nil
}

func (s *SonyDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &sonyrest.TV{
		Address: addr,
		PSK:     s.PSK,
	}, nil
}
