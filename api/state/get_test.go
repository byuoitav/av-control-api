package state

import (
	"context"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/log"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/drivertest"
	"github.com/byuoitav/av-control-api/drivers/mock"
	"github.com/google/go-cmp/cmp"
)

type getStateTest struct {
	name    string
	driver  drivertest.Driver
	apiResp api.StateResponse
}

var (
	via  = "ITB-1101-VIA1"
	hdmi = "ITB-1101-HDMI1"
	sign = "ITB-1101-SIGN1"
	pc   = "ITB-1101-PC1"
)

var getTests = []getStateTest{
	{
		name: "Simple",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.TV{},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(false),
					Blanked:   boolP(false),
					Inputs:    map[string]api.Input{},
					Volumes: map[string]int{
						"": 0,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
	},
}

func TestGetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			// build the room from the driver config
			room := api.Room{
				Devices: make(map[api.DeviceID]api.Device),
			}

			for id := range tt.driver.Devices {
				room.Devices[api.DeviceID(id)] = api.Device{
					Address: id,
				}
			}

			// start a driver server
			server := drivertest.NewServer(tt.driver.NewDeviceFunc())
			conn, err := server.GRPCClientConn(ctx)
			if err != nil {
				t.Fatalf("unable to get grpc client connection: %s", err)
			}

			// build the getSetter
			gs := &getSetter{
				log: log.Logger{},
				drivers: map[string]drivers.DriverClient{
					"": drivers.NewDriverClient(conn),
				},
			}

			// get the state of this room
			resp, err := gs.Get(ctx, room)
			if err != nil {
				t.Fatalf("unable to get room state: %s", err)
			}

			// compare the expected response to what we got
			if diff := cmp.Diff(tt.apiResp, resp); diff != "" {
				t.Errorf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}
