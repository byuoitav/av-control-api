package core

import (
	"context"
	"errors"

	"github.com/byuoitav/atlona/AT-HDVS-S10U"
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type Atlona2x1Driver struct {
	Username string
	Password string
	Log      *zap.Logger
}

func (a *Atlona2x1Driver) ParseConfig(config map[string]interface{}) error {
	if username, ok := config["username"].(string); ok {
		if username == "" {
			return errors.New("given empty username")
		}

		a.Username = username
	}

	if password, ok := config["password"].(string); ok {
		if password == "" {
			return errors.New("given empty password")
		}

		a.Password = password
	}

	return nil
}

func (a *Atlona2x1Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &atlona.AtlonaVideoSwitcher2x1{
		Username: a.Username,
		Password: a.Password,
		Address:  addr,
	}, nil
}
