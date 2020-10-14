package core

import (
	"context"
	"errors"

	"github.com/byuoitav/atlona-driver"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/wspool"
)

// func GetADCPDevice(ctx context.Context, addr string) (drivers.Device, error) {
// 	return &adcp.Projector{
// 		Address: addr,
// 	}, nil
// }

type atlonaBoy struct {
	Username string
	Password string
	// lol i don't think we can get a logger
	Log wspool.Logger
}

var bestBoy atlonaBoy

func ParseConfig(config map[string]interface{}) error {
	if username, ok := config["username"].(string); ok {
		if username == "" {
			return errors.New("given empty username")
		}

		bestBoy.Username = username
	}

	if password, ok := config["password"].(string); ok {
		if password == "" {
			return errors.New("given empty password")
		}

		bestBoy.Password = password
	}

	//something something logger

	return nil
}

func GetAtlonaDevice(ctx context.Context, addr, username, password string, log wspool.Logger) (drivers.Device, error) {
	return atlona.CreateVideoSwitcher(ctx, addr, bestBoy.Username, bestBoy.Password, log)
}
