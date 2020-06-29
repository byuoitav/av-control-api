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

type setStateTest struct {
	name    string
	log     bool
	driver  drivertest.Driver
	apiReq  api.StateRequest
	apiResp api.StateResponse
}

var setTests = []setStateTest{
	{
		name: "BasicTV/PowerOff",
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
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(false),
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(false),
				},
			},
		},
	},
	{
		name: "BasicTV/ChangeInput",
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
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Inputs: map[string]api.Input{
						"": api.Input{
							AudioVideo: stringP("hdmi3"),
						},
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Inputs: map[string]api.Input{
						"": api.Input{
							AudioVideo: stringP("hdmi3"),
						},
					},
				},
			},
		},
	},
	{
		name: "BasicTV/Blank",
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
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Blanked: boolP(true),
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Blanked: boolP(true),
				},
			},
		},
	},
	{
		name: "BasicTV/ChangeVolume",
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
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Volumes: map[string]int{
						"": 15,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Volumes: map[string]int{
						"": 15,
					},
				},
			},
		},
	},
	{
		name: "BasicTV/Mute",
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
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Mutes: map[string]bool{
						"": true,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					Mutes: map[string]bool{
						"": true,
					},
				},
			},
		},
	},
	{
		name: "BasicTV/PowerOn",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					On:               boolP(false),
					AudioVideoInputs: map[string]string{"": ""},
					Blanked:          boolP(false),
					Volumes: map[string]int{
						"": 0,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
		apiReq: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": api.DeviceState{
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Inputs: map[string]api.Input{
						"": api.Input{
							AudioVideo: stringP("hdmi2"),
						},
					},
					Volumes: map[string]int{
						"": 30,
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
					Blanked:   boolP(false),
					Inputs: map[string]api.Input{
						"": api.Input{
							AudioVideo: stringP("hdmi2"),
						},
					},
					Volumes: map[string]int{
						"": 30,
					},
					Mutes: map[string]bool{
						"": false,
					},
				},
			},
		},
	},
}

func TestSetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range setTests {
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
			resp, err := gs.Set(ctx, room, tt.apiReq)
			if err != nil {
				t.Fatalf("unable to set room state: %s", err)
			}

			// compare the expected response to what we got
			if diff := cmp.Diff(tt.apiResp, resp); diff != "" {
				t.Errorf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}
