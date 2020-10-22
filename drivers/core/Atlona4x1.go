package core

import (
	"context"
	"errors"

	"github.com/byuoitav/atlona/AT-JUNO-451-HDBT"
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type Atlona4x1Driver struct {
	Username string
	Password string
	Log      *zap.Logger
}

func (a *Atlona4x1Driver) ParseConfig(config map[string]interface{}) error {
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

func (a *Atlona4x1Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &atlona.AtlonaVideoSwitcher4x1{
		Username: a.Username,
		Password: a.Password,
		Address:  addr,
	}, nil
}
