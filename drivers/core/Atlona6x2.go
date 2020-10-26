package core

import (
	"context"
	"errors"

	atomeps62 "github.com/byuoitav/atlona/AT-OME-PS62"
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type Atlona6x2Driver struct {
	Username string
	Password string
	Log      *zap.Logger
}

func (a *Atlona6x2Driver) ParseConfig(config map[string]interface{}) error {
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

func (a *Atlona6x2Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &atomeps62.AtlonaVideoSwitcher6x2{
		Username: a.Username,
		Password: a.Password,
		Address:  addr,
	}, nil
}
