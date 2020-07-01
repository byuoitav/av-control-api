package state

import (
	"context"
	"errors"
	"fmt"
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
	name   string
	log    bool
	driver drivertest.Driver
	req    api.StateRequest
	err    error
	resp   api.StateResponse
}

var setTests = []setStateTest{
	{
		name: "EmptyRequest",
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
		req:  api.StateRequest{},
		resp: api.StateResponse{},
	},
	{
		name: "InvalidDevices",
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-SW1": {
					PoweredOn: boolP(false),
				},
			},
		},
		err: fmt.Errorf("ITB-1101-SW1: %s", ErrInvalidDevice),
	},
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi3"),
						},
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Inputs: map[string]api.Input{
						"": {
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Blanked: boolP(true),
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Volumes: map[string]int{
						"": 15,
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Mutes: map[string]bool{
						"": true,
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
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
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Inputs: map[string]api.Input{
						"": {
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
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Inputs: map[string]api.Input{
						"": {
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
	{
		name: "VideoSwitcher/ChangeInput1",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-SW1": &mock.Device{
					AudioInputs: map[string]string{
						"1": "1",
						"2": "2",
						"3": "3",
						"4": "4",
					},
					VideoInputs: map[string]string{
						"1": "4",
						"2": "3",
						"3": "2",
						"4": "1",
					},
				},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-SW1": {
					Inputs: map[string]api.Input{
						"1": {
							Audio: stringP("4"),
							Video: stringP("1"),
						},
						"2": {
							Audio: stringP("3"),
							Video: stringP("2"),
						},
						"3": {
							Audio: stringP("2"),
							Video: stringP("3"),
						},
						"4": {
							Audio: stringP("1"),
							Video: stringP("4"),
						},
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-SW1": {
					Inputs: map[string]api.Input{
						"1": {
							Audio: stringP("4"),
							Video: stringP("1"),
						},
						"2": {
							Audio: stringP("3"),
							Video: stringP("2"),
						},
						"3": {
							Audio: stringP("2"),
							Video: stringP("3"),
						},
						"4": {
							Audio: stringP("1"),
							Video: stringP("4"),
						},
					},
				},
			},
		},
	},
	{
		name: "Errors!1",
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D1": &mock.Device{
					SetPowerError:   errors.New("power error"),
					SetVolumeError:  errors.New("volume error"),
					SetBlankError:   errors.New("blank error"),
					UnknownCapError: errors.New("unknown"),
				},
				"ITB-1101-D2": &mock.Device{
					SetAudioVideoInputError: errors.New("av error"),
					SetAudioInputError:      errors.New("audio error"),
					SetVideoInputError:      errors.New("video error"),
					SetMuteError:            errors.New("mute error"),
				},
				"ITB-1101-D5": &mock.Device{},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Volumes:   map[string]int{"": 30},
				},
				"ITB-1101-D2": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
							Audio:      stringP("hdmi2"),
							Video:      stringP("hdmi2"),
						},
					},
					Mutes: map[string]bool{"": true},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {},
				"ITB-1101-D2": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D1",
					Field: "poweredOn",
					Value: true,
					Error: "power error",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "blanked",
					Value: false,
					Error: "blank error",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "volumes.",
					Value: int32(30),
					Error: "this output does not exist",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "input..audio",
					Value: "hdmi2",
					Error: "this output does not exist",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "input..video",
					Value: "hdmi2",
					Error: "this output does not exist",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "input..audioVideo",
					Value: "hdmi2",
					Error: "this output does not exist",
				},
				{
					ID:    "ITB-1101-D2",
					Field: "mutes.",
					Value: true,
					Error: "this output does not exist",
				},
			},
		},
	},
	{
		name: "NoCapabilitiesError",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D3": &mock.Device{
					GetCapsError: errors.New("no capabilities"),
				},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D3": {},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D3": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D3",
					Error: "unable to get capabilities: no capabilities",
				},
			},
		},
	},
	{
		name: "CantSetPower",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					PoweredOn: boolP(true),
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "poweredOn",
					Value: true,
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetBlank",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Blanked: boolP(true),
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "blanked",
					Value: true,
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetVolumes",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Volumes: map[string]int{"": 10},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "volumes",
					Value: map[string]int{"": 10},
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetMutes",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Mutes: map[string]bool{"": true},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "mutes",
					Value: map[string]bool{"": true},
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetAudioInputs",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{VideoInputs: map[string]string{"": "hdmi2"}},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							Audio: stringP("hdmi2"),
							Video: stringP("hdmi2"),
						},
						"other": {
							Audio: stringP("hdmi2"),
						},
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							Video: stringP("hdmi2"),
						},
					},
				},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "input.$.audio",
					Value: map[string]api.Input{
						"": {
							Audio: stringP("hdmi2"),
							Video: stringP("hdmi2"),
						},
						"other": {
							Audio: stringP("hdmi2"),
						},
					},
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetVideoInputs",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{AudioVideoInputs: map[string]string{"": "hdmi2"}},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
							Video:      stringP("hdmi2"),
						},
						"other": {
							Video: stringP("hdmi2"),
						},
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
						},
					},
				},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "input.$.video",
					Value: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
							Video:      stringP("hdmi2"),
						},
						"other": {
							Video: stringP("hdmi2"),
						},
					},
					Error: "can't set this field on this device",
				},
			},
		},
	},
	{
		name: "CantSetAudioVideoInputs",
		log:  true,
		driver: drivertest.Driver{
			Devices: map[string]drivers.Device{
				"ITB-1101-D4": &mock.Device{AudioInputs: map[string]string{"": "hdmi2"}},
			},
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
							Audio:      stringP("hdmi2"),
						},
						"other": {
							AudioVideo: stringP("hdmi2"),
						},
					},
				},
			},
		},
		resp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D4": {
					Inputs: map[string]api.Input{
						"": {
							Audio: stringP("hdmi2"),
						},
					},
				},
			},
			Errors: []api.DeviceStateError{
				{
					ID:    "ITB-1101-D4",
					Field: "input.$.audioVideo",
					Value: map[string]api.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
							Audio:      stringP("hdmi2"),
						},
						"other": {
							AudioVideo: stringP("hdmi2"),
						},
					},
					Error: "can't set this field on this device",
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

			ctx = api.WithRequestID(ctx, "ID")

			// get the state of this room
			resp, err := gs.Set(ctx, room, tt.req)
			if tt.err != nil {
				if diff := cmp.Diff(tt.err.Error(), err.Error()); diff != "" {
					t.Fatalf("generated incorrect error (-want, +got):\n%s", diff)
				}
			}

			// compare the expected response to what we got
			if diff := cmp.Diff(tt.resp, resp); diff != "" {
				t.Fatalf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestSetWrongDriver(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := api.StateRequest{
		Devices: map[api.DeviceID]api.DeviceState{
			"ITB-1101-D1": {},
		},
	}
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

		_, err = gs.Set(ctx, room, req)
		if err != nil {
			if diff := cmp.Diff(errWanted.Error(), err.Error()); diff != "" {
				t.Fatalf("generated incorrect error (-want, +got):\n%s", diff)
			}
		}
	})
}
