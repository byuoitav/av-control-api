package core

import (
	"context"
	"errors"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/kramer/via"
	"go.uber.org/zap"
)

type KramerViaDriver struct {
	Log      *zap.Logger
	Username string
	Password string
}

func (k *KramerViaDriver) ParseConfig(config map[string]interface{}) error {
	if username, ok := config["username"].(string); ok {
		if username == "" {
			k.Log.Info("we should get here!!!!\n\n\n\n")
			return errors.New("given empty username")
		}

		k.Username = username
	} else {
		return errors.New("no username given")
	}

	if password, ok := config["password"].(string); ok {
		if password == "" {
			return errors.New("given empty password")
		}

		k.Password = password
	} else {
		return errors.New("no password given")
	}

	return nil
}

func (k *KramerViaDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return &via.Via{
		Address:  addr,
		Username: k.Username,
		Password: k.Password,
		Log:      k.Log,
	}, nil
}
