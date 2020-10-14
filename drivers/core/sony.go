package core

import (
	"context"
	"errors"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/sonyrest-driver"
)

type sonyBoy struct {
	PSK string
}

var otherBoy sonyBoy

func ParseSonyConfig(config map[string]interface{}) error {
	if psk, ok := config["psk"].(string); ok {
		if psk == "" {
			return errors.New("given empty psk")
		}

		otherBoy.PSK = psk
	}

	return nil
}

func GetSonyDevice(ctx context.Context, addr string) (drivers.Device, error) {
	return &sonyrest.TV{
		Address: addr,
		PSK:     otherBoy.PSK,
	}, nil
}
