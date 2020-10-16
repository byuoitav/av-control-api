package core

import (
	"context"
	"errors"

	"github.com/byuoitav/atlona-driver"
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/wspool"
)

// func GetADCPDevice(ctx context.Context, addr string) (drivers.Device, error) {
// 	return &adcp.Projector{
// 		Address: addr,
// 	}, nil
// }

type AtlonaDriver struct {
	Username string
	Password string
	// lol i don't think we can get a logger
	Log wspool.Logger
}

func (a *AtlonaDriver) ParseConfig(config map[string]interface{}) error {
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

	//something something logger

	return nil
}

func (a *AtlonaDriver) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	return atlona.CreateVideoSwitcher(ctx, addr, a.Username, a.Password, a.Log)
}
