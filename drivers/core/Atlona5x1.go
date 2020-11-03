package core

import (
	"context"
	"errors"

	atuhdsw52ed "github.com/byuoitav/atlona/AT-UHD-SW-52ED"
	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type Atlona5x1Driver struct {
	Username string
	Password string
	Log      *zap.Logger
}

func (a *Atlona5x1Driver) ParseConfig(config map[string]interface{}) error {
	if username, ok := config["username"].(string); ok {
		if username == "" {
			return errors.New("given empty username")
		}

		a.Username = username
	} else {
		// driver still works so don't need username/password?
		// return errors.New("no username given")
	}

	if password, ok := config["password"].(string); ok {
		if password == "" {
			return errors.New("given empty password")
		}

		a.Password = password
	} else {
		// return errors.New("no password given")
	}

	return nil
}

func (a *Atlona5x1Driver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &atuhdsw52ed.AtlonaVideoSwitcher5x1{
		Username: a.Username,
		Password: a.Password,
		Address:  addr,
		Logger:   a.Log.Sugar(),
	}, nil
}
