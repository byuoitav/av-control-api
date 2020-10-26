package state

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/driverstest"
	"github.com/byuoitav/av-control-api/mock"
	"github.com/matryer/is"
	"go.uber.org/zap"
)

type getStateTest struct {
	name   string
	log    bool
	driver *driverstest.Driver
	resp   avcontrol.StateResponse
}

var getTests = []getStateTest{
	{
		name: "BasicTV",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi3",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: true,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 69,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": true,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi3"),
						},
					},
					Blanked: boolP(true),
					Volumes: map[string]int{
						"": 69,
					},
					Mutes: map[string]bool{
						"": true,
					},
				},
			},
		},
	},
	{
		name: "SeparateVolumeMute",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: false,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"out": "hdmi1",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"headphones": 42,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"aux": false,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
					Inputs: map[string]avcontrol.Input{
						"out": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{
						"headphones": 42,
					},
					Mutes: map[string]bool{
						"aux": false,
					},
				},
			},
		},
	},
	{
		name: "SimpleSeparateInput",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TVSeparateInput{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioInput: mock.WithAudioInput{
						Inputs: map[string]string{
							"": "hdmi2",
						},
					},
					WithVideoInput: mock.WithVideoInput{
						Inputs: map[string]string{
							"": "hdmi4",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 77,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
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
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithAudioInput: mock.WithAudioInput{
						Inputs: map[string]string{
							"1": "in1",
							"2": "in2",
							"3": "in3",
							"4": "in4",
						},
					},
					WithVideoInput: mock.WithVideoInput{
						Inputs: map[string]string{
							"1": "in4",
							"2": "in3",
							"3": "in2",
							"4": "in1",
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-SW1": {
					Inputs: map[string]avcontrol.Input{
						"1": {
							Audio: stringP("in1"),
							Video: stringP("in4"),
						},
						"2": {
							Audio: stringP("in2"),
							Video: stringP("in3"),
						},
						"3": {
							Audio: stringP("in3"),
							Video: stringP("in2"),
						},
						"4": {
							Audio: stringP("in4"),
							Video: stringP("in1"),
						},
					},
				},
			},
		},
	},
	{
		name: "BasicVideoSwitcherRoom",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi1",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 50,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
				"ITB-1101-SW1": mock.BasicVideoSwitcher{
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"1": "hdmi1",
						},
					},
				},
				"ITB-1101-AMP1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 30,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Inputs: map[string]avcontrol.Input{
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
					Inputs: map[string]avcontrol.Input{
						"1": {
							AudioVideo: stringP("hdmi1"),
						},
					},
				},
				"ITB-1101-AMP1": {
					Volumes: map[string]int{"": 30},
					Mutes:   map[string]bool{"": false},
				},
			},
		},
	},
	{
		name: "6x2SeparateInput",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.Projector{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdbaset",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
				},
				"ITB-1101-D2": mock.Projector{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi1",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: true,
					},
				},
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithAudioInput: mock.WithAudioInput{
						Inputs: map[string]string{
							"hdmiOutA": "1",
							"hdmiOutB": "2",
						},
					},
					WithVideoInput: mock.WithVideoInput{
						Inputs: map[string]string{
							"hdmiOutA": "2",
							"hdmiOutB": "1",
						},
					},
				},
				"ITB-1101-VIA1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 69,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdbaset"),
						},
					},
					Blanked: boolP(false),
				},
				"ITB-1101-D2": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(true),
				},
				"ITB-1101-SW1": {
					Inputs: map[string]avcontrol.Input{
						"hdmiOutA": {
							Audio: stringP("1"),
							Video: stringP("2"),
						},
						"hdmiOutB": {
							Audio: stringP("2"),
							Video: stringP("1"),
						},
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{"": 69},
					Mutes:   map[string]bool{"": false},
				},
			},
		},
	},
	{
		name: "JustAddPower",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi1",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: true,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 50,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": true,
						},
					},
				},
				"ITB-1101-D2": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi2",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 80,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
				"ITB-1101-D3": struct {
					mock.WithPower
				}{
					mock.WithPower{
						PoweredOn: true,
					},
				},
				"ITB-1101-RX1": mock.AVOverIPReceiver{
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "10.66.78.155",
						},
					},
				},
				"ITB-1101-RX2": mock.AVOverIPReceiver{
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "10.66.78.156",
						},
					},
				},
				"ITB-1101-RX3": mock.AVOverIPReceiver{
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "10.66.78.157",
						},
					},
				},
				"ITB-1101-VIA1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 60,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
				"ITB-1101-DSP1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"1": 70,
							"2": 30,
							"3": 45,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"1": false,
							"2": true,
							"3": true,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
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
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
						},
					},
					Blanked: boolP(false),
					Volumes: map[string]int{"": 80},
					Mutes:   map[string]bool{"": false},
				},
				"ITB-1101-D3": {
					PoweredOn: boolP(true),
				},
				"ITB-1101-RX1": {
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("10.66.78.155"),
						},
					},
				},
				"ITB-1101-RX2": {
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("10.66.78.156"),
						},
					},
				},
				"ITB-1101-RX3": {
					Inputs: map[string]avcontrol.Input{
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
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"JRCB-205-D1": mock.Projector{
					WithPower: mock.WithPower{
						PoweredOn: true,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdmi1",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: true,
					},
				},
				"JRCB-205-D2": mock.Projector{
					WithPower: mock.WithPower{
						PoweredOn: false,
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
							"": "hdbaset",
						},
					},
					WithBlank: mock.WithBlank{
						Blanked: false,
					},
				},
				"JRCB-205-DSP1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
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
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
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
				},
				"JRCB-205-SW1": mock.BasicVideoSwitcher{
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Inputs: map[string]string{
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
				},
				"JRCB-205-VIA1": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 40,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": true,
						},
					},
				},
				"JRCB-205-VIA2": mock.DSP{
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"": 30,
						},
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"": false,
						},
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"JRCB-205-D1": {
					PoweredOn: boolP(true),
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi1"),
						},
					},
					Blanked: boolP(true),
				},
				"JRCB-205-D2": {
					PoweredOn: boolP(false),
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdbaset"),
						},
					},
					Blanked: boolP(false),
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
					Inputs: map[string]avcontrol.Input{
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
		name: "Errors",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						Error: errors.New("can't get power"),
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						Error: errors.New("can't get audio video inputs"),
					},
					WithBlank: mock.WithBlank{
						Error: errors.New("can't get blank"),
					},
					WithVolume: mock.WithVolume{
						Error: errors.New("can't get volumes"),
					},
					WithMute: mock.WithMute{
						Error: errors.New("can't get mutes"),
					},
				},
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithAudioInput: mock.WithAudioInput{
						Error: errors.New("can't get audio inputs"),
					},
					WithVideoInput: mock.WithVideoInput{
						Error: errors.New("can't get video inputs"),
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1":  {},
				"ITB-1101-SW1": {},
			},
			Errors: []avcontrol.DeviceStateError{
				{
					ID:    "ITB-1101-D1",
					Field: "blank",
					Error: "can't get blank",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "inputs.$.audioVideo",
					Error: "can't get audio video inputs",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "mutes",
					Error: "can't get mutes",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "power",
					Error: "can't get power",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "volumes",
					Error: "can't get volumes",
				},
				{
					ID:    "ITB-1101-SW1",
					Field: "inputs.$.audio",
					Error: "can't get audio inputs",
				},
				{
					ID:    "ITB-1101-SW1",
					Field: "inputs.$.video",
					Error: "can't get video inputs",
				},
			},
		},
	},
	{
		name:   "EmptyRoom",
		driver: &driverstest.Driver{},
		resp:   avcontrol.StateResponse{},
	},
}

func TestGetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			// build the room from the driver config
			room := avcontrol.RoomConfig{
				Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceConfig),
			}

			// create each of the devices
			for id, dev := range tt.driver.Devices {
				var c avcontrol.DeviceConfig
				c.Address = id
				c.Driver = "driverstest/driver"

				val := reflect.ValueOf(dev)
				for i := 0; i < val.NumField(); i++ {
					if !val.Field(i).CanInterface() {
						continue
					}

					field := val.Field(i).Interface()

					if d, ok := field.(mock.WithVolume); ok {
						for block := range d.Vols {
							c.Ports = append(c.Ports, avcontrol.PortConfig{
								Name: block,
								Type: "volume",
							})
						}
					}

					if d, ok := field.(mock.WithMute); ok {
						for block := range d.Ms {
							c.Ports = append(c.Ports, avcontrol.PortConfig{
								Name: block,
								Type: "mute",
							})
						}
					}
				}

				room.Devices[avcontrol.DeviceID(id)] = c
			}

			// need a way to not pass a file
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

			// get the state of this room
			resp, err := gs.Get(ctx, room)
			is.NoErr(err)
			is.Equal(resp, tt.resp)
		})
	}
}

/*
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
*/
