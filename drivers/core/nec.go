package core

import (
	"context"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/nec-driver"
)

type NECDriver struct{}

func (n *NECDriver) ParseConfig(config map[string]interface{}) error {
	return nil
}

func (n *NECDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return nec.NewProjector(addr, nec.WithDelay(300*time.Second)), nil
}
