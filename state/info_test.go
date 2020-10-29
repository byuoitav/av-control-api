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

type getInfoTest struct {
	name   string
	log    bool
	driver *driverstest.Driver
	resp   avcontrol.RoomInfo
}

type resp struct {
	response string
}

var getInfoTests = []getInfoTest{
	{
		name: "BasicTV",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithInfo: mock.WithInfo{
						I: resp{
							response: "hello",
						},
					},
				},
			},
		},
		resp: avcontrol.RoomInfo{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceInfo{
				"ITB-1101-D1": {
					Info: resp{
						response: "hello",
					},
				},
			},
		},
	},
	{
		name: "VideoSwitcher",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithInfo: mock.WithInfo{
						Error: errors.New("failed to get info"),
					},
				},
			},
		},
		resp: avcontrol.RoomInfo{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceInfo{
				"ITB-1101-SW1": {
					Error: stringP("failed to get info"),
				},
			},
		},
	},
}

func TestInfo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range getInfoTests {
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

			_, err = gs.GetInfo(ctx, room)
			is.NoErr(err)
		})
	}
}
