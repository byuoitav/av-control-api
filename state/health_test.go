package state

import (
	"context"
	"errors"
	"testing"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/driverstest"
	"github.com/byuoitav/av-control-api/mock"
	"github.com/matryer/is"
	"go.uber.org/zap"
)

type getHealthTest struct {
	name   string
	log    bool
	driver *driverstest.Driver
	resp   avcontrol.RoomHealth
}

var getHealthTests = []getHealthTest{
	{
		name: "BasicTV",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithHealth: mock.WithHealth{
						Error: nil,
					},
				},
			},
		},
		resp: avcontrol.RoomHealth{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceHealth{
				"ITB-1101-D1": {
					Healthy: boolP(true),
				},
			},
		},
	},
	{
		name: "VideoSwitcher",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithHealth: mock.WithHealth{
						Error: errors.New("failed health check"),
					},
				},
			},
		},
		resp: avcontrol.RoomHealth{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceHealth{
				"ITB-1101-SW1": {
					Healthy: boolP(false),
					Error:   stringP("failed health check"),
				},
			},
		},
	},
}

func TestHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range getHealthTests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			room := avcontrol.RoomConfig{
				Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceConfig),
			}

			for id := range tt.driver.Devices {
				var c avcontrol.DeviceConfig
				c.Address = id
				c.Driver = "driverstest/driver"

				room.Devices[avcontrol.DeviceID(id)] = c
			}

			registry, err := drivers.NewWithConfig(nil)
			is.NoErr(err)

			err = registry.Register("driverstest/driver", tt.driver)
			is.NoErr(err)

			// build the getSetter
			gs := &GetSetter{
				Logger:         zap.NewNop(),
				DriverRegistry: registry,
			}

			if tt.log {
				gs.Logger = zap.NewExample()
			}

			ctx = avcontrol.WithRequestID(ctx, "ID")

			resp, err := gs.GetHealth(ctx, room)
			is.NoErr(err)
			is.Equal(resp, tt.resp)
		})
	}
}
