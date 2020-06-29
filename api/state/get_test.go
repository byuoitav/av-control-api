package state

import (
	"context"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/drivertest"
	"github.com/byuoitav/av-control-api/drivers/mock"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
)

type getStateTest struct {
	name    string
	log     bool
	driver  drivertest.Driver
	apiResp api.StateResponse
}

var getTests = []getStateTest{
	{
		name: "BasicTV",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					Blanked:          boolP(false),
					Volumes: map[string]int{
						"": 77,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": api.Input{
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{
						"": 77,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
	},
	{
		name: "SeparateVolumeMute",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					On: boolP(false),
					AudioVideoInputs: map[string]string{
						"out": "hdmi3",
					},
					Blanked: boolP(true),
					Volumes: map[string]int{
						"headphones": 12,
					},
					Mutes: map[string]bool{
						"aux": true,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(false),
					Inputs: map[string]api.Input{
						"out": api.Input{
							AudioVideo: stringP("hdmi3"),
						},
					},
					Blanked: boolP(true),
					Volumes: map[string]int{
						"headphones": 12,
					},
					Mutes: map[string]bool{
						"aux": true,
					},
				},
			},
		},
	},
	{
		name: "SimpleSeparateInput",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					On:          boolP(true),
					AudioInputs: map[string]string{"": "hdmi2"},
					VideoInputs: map[string]string{"": "hdmi4"},
					Blanked:     boolP(false),
					Volumes: map[string]int{
						"": 77,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": api.Input{
							Audio: stringP("hdmi2"),
							Video: stringP("hdmi4"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{
						"": 77,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
	},
	{
		name: "VideoSwitcherSeparateInputs",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					AudioInputs: map[string]string{
						"1": "in1",
						"2": "in2",
						//"3": "in3",
						//"4": "in4",
					},
					VideoInputs: map[string]string{
						"1": "in4",
						"2": "in3",
						//"3": "in2",
						//"4": "in1",
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Inputs: map[string]api.Input{
						"1": api.Input{
							Audio: stringP("in1"),
							Video: stringP("in4"),
						},
						"2": api.Input{
							Audio: stringP("in2"),
							Video: stringP("in3"),
						},
						//"3": api.Input{
						//	Audio: stringP("in3"),
						//	Video: stringP("in2"),
						//},
						//"4": api.Input{
						//	Audio: stringP("in4"),
						//	Video: stringP("in1"),
						//},
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

			for id, dev := range tt.driver.Devices {
				var apiDev api.Device
				apiDev.Address = id

				if d, ok := dev.(*mock.Device); ok {
					vols := d.VolumeBlocks()
					mutes := d.MuteBlocks()

					for _, block := range vols {
						apiDev.Ports = append(apiDev.Ports, api.Port{
							Name: block,
							Type: "volume",
						})
					}

					for _, block := range mutes {
						apiDev.Ports = append(apiDev.Ports, api.Port{
							Name: block,
							Type: "mute",
						})
					}
				}

				room.Devices[api.DeviceID(id)] = apiDev
			}

			// start a driver server
			server := drivertest.NewServer(tt.driver.NewDeviceFunc())
			conn, err := server.GRPCClientConn(ctx)
			if err != nil {
				t.Fatalf("unable to get grpc client connection: %s", err)
			}

			// build the getSetter
			gs := &getSetter{
				logger: zap.NewNop(),
				drivers: map[string]drivers.DriverClient{
					"": drivers.NewDriverClient(conn),
				},
			}

			if tt.log {
				gs.logger = zap.NewExample()
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
