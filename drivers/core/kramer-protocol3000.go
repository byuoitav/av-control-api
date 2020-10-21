package core

import (
	"context"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer/protocol3000"
)

type KramerProtocol3000Driver struct{}

func (k *KramerProtocol3000Driver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (k *KramerProtocol3000Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return protocol3000.New(addr), nil
}
