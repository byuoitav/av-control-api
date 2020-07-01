package state

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/drivertest"
	"github.com/byuoitav/av-control-api/drivers/mock"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
)

type contextKey int

const (
	_keyRequestID contextKey = iota
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
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
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
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
					Inputs: map[string]api.Input{
						"out": {
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
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
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
				"ITB-1101-D1": {
					Inputs: map[string]api.Input{
						"1": {
							Audio: stringP("in1"),
							Video: stringP("in4"),
						},
						"2": {
							Audio: stringP("in2"),
							Video: stringP("in3"),
						},
					},
				},
			},
		},
	},
	{
		name: "SimpleSeparateInput",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					On:               boolP(true),
					Blanked:          boolP(false),
					Volumes:          map[string]int{"": 50},
					Mutes:            map[string]bool{"": false},
				},
				"ITB-1101-SW1": &mock.Device{
					AudioVideoInputs: map[string]string{
						"1": "hdmi1",
					},
				},
				"ITB-1101-AMP1": &mock.Device{
					Volumes: map[string]int{"": 30},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-HDMI1": &mock.Device{},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Volumes:   map[string]int{"": 50},
					Mutes:     map[string]bool{"": false},
				},
				"ITB-1101-SW1": {
					Inputs: map[string]api.Input{
						"1": {
							AudioVideo: stringP("hdmi1"),
						},
					},
				},
				"ITB-1101-AMP1": {
					Volumes: map[string]int{"": 30},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-HDMI1": {},
			},
		},
	},
	{
		name: "6x2SeparateInput",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				//Projector
				"ITB-1101-D1": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					Blanked:          boolP(false),
				},
				"ITB-1101-D2": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdbaset"},
					Blanked:          boolP(true),
				},
				"ITB-1101-SW1": &mock.Device{
					AudioInputs: map[string]string{
						"1": "hdmi1",
						"2": "hdbaset",
					},
					VideoInputs: map[string]string{
						"1": "hdbaset",
						"2": "hdmi1",
					},
				},
				"ITB-1101-VIA1": &mock.Device{
					Volumes: map[string]int{"": 69},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-HDMI1": &mock.Device{
					Volumes: map[string]int{"": 50},
					Mutes:   map[string]bool{"": true},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(false),
				},
				"ITB-1101-D2": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdbaset"),
						},
					},
					Blanked: boolP(true),
				},
				"ITB-1101-SW1": {
					Inputs: map[string]api.Input{
						"1": {
							Audio: stringP("hdmi1"),
							Video: stringP("hdbaset"),
						},
						"2": {
							Audio: stringP("hdbaset"),
							Video: stringP("hdmi1"),
						},
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{"": 69},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-HDMI1": {
					Volumes: map[string]int{"": 50},
					Mutes:   map[string]bool{"": true},
				},
			},
		},
	},
	{
		name: "JustAddPower",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					Blanked:          boolP(true),
					Volumes:          map[string]int{"": 50},
					Mutes:            map[string]bool{"": true},
				},
				"ITB-1101-D2": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdmi2"},
					Blanked:          boolP(false),
					Volumes:          map[string]int{"": 80},
					Mutes:            map[string]bool{"": false},
				},
				"ITB-1101-D3": &mock.Device{
					On: boolP(false),
				},
				"ITB-1101-RX1": &mock.Device{
					AudioVideoInputs: map[string]string{"": "10.66.78.155"},
				},
				"ITB-1101-RX2": &mock.Device{
					AudioVideoInputs: map[string]string{"": "10.66.78.156"},
				},
				"ITB-1101-RX3": &mock.Device{
					AudioVideoInputs: map[string]string{"": "10.66.78.157"},
				},
				"ITB-1101-VIA1": &mock.Device{
					Volumes: map[string]int{"": 60},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-DSP1": &mock.Device{
					Volumes: map[string]int{
						"1": 70,
						"2": 30,
						"3": 45,
					},
					Mutes: map[string]bool{
						"1": false,
						"2": true,
						"3": true,
					},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(true),
					Volumes: map[string]int{"": 50},
					Mutes:   map[string]bool{"": true},
				},
				"ITB-1101-D2": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{"": 80},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-D3": {
					PoweredOn: boolP(false),
				},
				"ITB-1101-RX1": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("10.66.78.155"),
						},
					},
				},
				"ITB-1101-RX2": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("10.66.78.156"),
						},
					},
				},
				"ITB-1101-RX3": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("10.66.78.157"),
						},
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{"": 60},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-DSP1": {
					Volumes: map[string]int{
						"1": 70,
						"2": 30,
						"3": 45,
					},
					Mutes: map[string]bool{
						"1": false,
						"2": true,
						"3": true,
					},
				},
			},
		},
	},
	{
		name: "JRCB-205",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"JRCB-205-D1": &mock.Device{
					On:               boolP(true),
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					Blanked:          boolP(true),
					Volumes:          map[string]int{"": 50},
					Mutes:            map[string]bool{"": false},
				},
				"JRCB-205-D2": &mock.Device{
					On:               boolP(false),
					AudioVideoInputs: map[string]string{"": "hdmi1"},
					Blanked:          boolP(false),
					Volumes:          map[string]int{"": 42},
					Mutes:            map[string]bool{"": true},
				},
				"JRCB-205-DSP1": &mock.Device{
					Volumes: map[string]int{
						"1":  10,
						"2":  15,
						"3":  20,
						"4":  25,
						"5":  30,
						"6":  35,
						"7":  40,
						"8":  45,
						"9":  50,
						"10": 55,
						"11": 60,
						"12": 65,
						"13": 70,
						"14": 75,
						"15": 80,
						"16": 85,
						"17": 90,
					},
					Mutes: map[string]bool{
						"1":  true,
						"2":  false,
						"3":  true,
						"4":  false,
						"5":  true,
						"6":  false,
						"7":  true,
						"8":  false,
						"9":  true,
						"10": false,
						"11": false,
						"12": true,
						"13": false,
						"14": true,
						"15": false,
						"16": true,
						"17": false,
					},
				},
				"JRCB-205-SW1": &mock.Device{
					AudioVideoInputs: map[string]string{
						"1": "1",
						"2": "1",
						"3": "4",
						"4": "3",
						"5": "5",
						"6": "2",
						"7": "2",
						"8": "10",
					},
				},
				"JRCB-205-VIA1": &mock.Device{
					Volumes: map[string]int{"": 40},
					Mutes:   map[string]bool{"": true},
				},
				"JRCB-205-VIA2": &mock.Device{
					Volumes: map[string]int{"": 30},
					Mutes:   map[string]bool{"": false},
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"JRCB-205-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(true),
					Volumes: map[string]int{"": 50},
					Mutes:   map[string]bool{"": false},
				},
				"JRCB-205-D2": {
					PoweredOn: boolP(false),
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{"": 42},
					Mutes:   map[string]bool{"": true},
				},
				"JRCB-205-DSP1": {
					Volumes: map[string]int{
						"1":  10,
						"2":  15,
						"3":  20,
						"4":  25,
						"5":  30,
						"6":  35,
						"7":  40,
						"8":  45,
						"9":  50,
						"10": 55,
						"11": 60,
						"12": 65,
						"13": 70,
						"14": 75,
						"15": 80,
						"16": 85,
						"17": 90,
					},
					Mutes: map[string]bool{
						"1":  true,
						"2":  false,
						"3":  true,
						"4":  false,
						"5":  true,
						"6":  false,
						"7":  true,
						"8":  false,
						"9":  true,
						"10": false,
						"11": false,
						"12": true,
						"13": false,
						"14": true,
						"15": false,
						"16": true,
						"17": false,
					},
				},
				"JRCB-205-SW1": {
					Inputs: map[string]api.Input{
						"1": {
							AudioVideo: stringP("1"),
						},
						"2": {
							AudioVideo: stringP("1"),
						},
						"3": {
							AudioVideo: stringP("4"),
						},
						"4": {
							AudioVideo: stringP("3"),
						},
						"5": {
							AudioVideo: stringP("5"),
						},
						"6": {
							AudioVideo: stringP("2"),
						},
						"7": {
							AudioVideo: stringP("2"),
						},
						"8": {
							AudioVideo: stringP("10"),
						},
					},
				},
				"JRCB-205-VIA1": {
					Volumes: map[string]int{"": 40},
					Mutes:   map[string]bool{"": true},
				},
				"JRCB-205-VIA2": {
					Volumes: map[string]int{"": 30},
					Mutes:   map[string]bool{"": false},
				},
			},
		},
	},
	{
		name: "Errors!",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					GetPowerError:   errors.New("power error"),
					GetVolumesError: errors.New("volume error"),
					GetBlankError:   errors.New("blank error"),
					UnknownCapError: errors.New("unknown error"),
				},
				"ITB-1101-D2": &mock.Device{
					GetAudioVideoInputsError: errors.New("av error"),
					GetAudioInputsError:      errors.New("audio error"),
					GetVideoInputsError:      errors.New("video error"),
					GetMutesError:            errors.New("mutes error"),
				},
				"ITB-1101-D3": &mock.Device{
					GetCapsError: errors.New("no capabilities"),
				},
			},
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {},
				"ITB-1101-D2": {},
				"ITB-1101-D3": {
					Inputs:  map[string]api.Input{},
					Volumes: map[string]int{},
					Mutes:   map[string]bool{},
				},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D3",
					Error: "unable to get capabilities: no capabilities",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "power",
					Error: "power error",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "blank",
					Error: "blank error",
				},
				{
					ID:    "ITB-1101-D1",
					Error: "unknown capability unknown",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "volumes",
					Error: "volume error",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "inputs.$.audio",
					Error: "audio error",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "inputs.$.video",
					Error: "video error",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "inputs.$.audioVideo",
					Error: "av error",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "mutes",
					Error: "mutes error",
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

			ctx = api.WithRequestID(ctx, "ID")

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

func TestGetWrongDriver(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	errWanted := errors.New("unknown driver: bad driver")

	t.Run("", func(t *testing.T) {
		room := api.Room{
			Devices: make(map[api.DeviceID]api.Device),
		}

		apiDev := api.Device{
			Address: "ITB-1101-D1",
			Driver:  "bad driver",
		}
		room.Devices[api.DeviceID("ITB-1101-D1")] = apiDev

		fakeDriver := drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D2": &mock.Device{},
			},
		}

		server := drivertest.NewServer(fakeDriver.NewDeviceFunc())
		conn, err := server.GRPCClientConn(ctx)
		if err != nil {
			t.Fatalf("unable to get grpc client connection: %s", err)
		}

		gs := &getSetter{
			logger: zap.NewNop(),
			drivers: map[string]drivers.DriverClient{
				"": drivers.NewDriverClient(conn),
			},
		}

		_, err = gs.Get(ctx, room)
		if err != nil {
			if diff := cmp.Diff(errWanted.Error(), err.Error()); diff != "" {
				t.Fatalf("generated incorrect error (-want, +got):\n%s", diff)
			}
		}
	})
}
