package state

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/driverstest"
	"github.com/byuoitav/av-control-api/mock"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type setStateTest struct {
	name   string
	log    bool
	driver *driverstest.Driver
	req    avcontrol.StateRequest
	err    error
	resp   avcontrol.StateResponse
}

var setTests = []setStateTest{
	{
		name: "EmptyRequest",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.Projector{
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
				},
			},
		},
		req:  avcontrol.StateRequest{},
		resp: avcontrol.StateResponse{},
	},
	{
		name: "InvalidDevices",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.Projector{
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
				},
			},
		},
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-SW1": {
					PoweredOn: boolP(false),
				},
			},
		},
		err: fmt.Errorf("ITB-1101-SW1: %s", ErrInvalidDevice),
	},
	{
		name: "BasicTV/PowerOff",
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
				},
			},
		},
	},
	{
		name: "BasicTV/ChangeInput",
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi3"),
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
							AudioVideo: stringP("hdmi3"),
						},
					},
				},
			},
		},
	},
	{
		name: "BasicTV/Blank",
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Blanked: boolP(true),
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Blanked: boolP(true),
				},
			},
		},
	},
	{
		name: "BasicTV/ChangeVolume",
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Volumes: map[string]int{
						"": 15,
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					Mutes: map[string]bool{
						"": true,
					},
				},
			},
		},
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
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
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						PoweredOn: false,
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
							"": 0,
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
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Inputs: map[string]avcontrol.Input{
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
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Inputs: map[string]avcontrol.Input{
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
		name: "VideoSwitcher/ChangeInput",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithAudioInput: mock.WithAudioInput{
						Inputs: map[string]string{
							"1": "1",
							"2": "2",
							"3": "3",
							"4": "4",
						},
					},
					WithVideoInput: mock.WithVideoInput{
						Inputs: map[string]string{
							"1": "4",
							"2": "3",
							"3": "2",
							"4": "1",
						},
					},
				},
			},
		},
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-SW1": {
					Inputs: map[string]avcontrol.Input{
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
		resp: avcontrol.StateResponse{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-SW1": {
					Inputs: map[string]avcontrol.Input{
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
		name: "Errors",
		driver: &driverstest.Driver{
			Devices: map[string]avcontrol.Device{
				"ITB-1101-D1": mock.TV{
					WithPower: mock.WithPower{
						SetError: errors.New("can't set power"),
					},
					WithAudioVideoInput: mock.WithAudioVideoInput{
						SetError: errors.New("can't set audio video inputs"),
					},
					WithBlank: mock.WithBlank{
						SetError: errors.New("can't set blank"),
					},
					WithVolume: mock.WithVolume{
						Vols: map[string]int{
							"headphones": 100,
						},
						SetError: errors.New("can't set volumes"),
					},
					WithMute: mock.WithMute{
						Ms: map[string]bool{
							"aux": false,
						},
						SetError: errors.New("can't set mutes"),
					},
				},
				"ITB-1101-SW1": mock.VideoSwitcher{
					WithAudioInput: mock.WithAudioInput{
						SetError: errors.New("can't set audio inputs"),
					},
					WithVideoInput: mock.WithVideoInput{
						SetError: errors.New("can't set video inputs"),
					},
				},
			},
		},
		req: avcontrol.StateRequest{
			Devices: map[avcontrol.DeviceID]avcontrol.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Volumes:   map[string]int{"headphones": 30},
					Mutes:     map[string]bool{"aux": true},
					Inputs: map[string]avcontrol.Input{
						"": {
							AudioVideo: stringP("hdmi2"),
						},
					},
				},
				"ITB-1101-SW1": {
					Inputs: map[string]avcontrol.Input{
						"hdmiOutA": {
							Audio: stringP("1"),
							Video: stringP("2"),
						},
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
					Field: "blanked",
					Value: false,
					Error: "can't set blank",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "input..audioVideo",
					Value: "hdmi2",
					Error: "can't set audio video inputs",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "mutes.aux",
					Value: true,
					Error: "can't set mutes",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "poweredOn",
					Value: true,
					Error: "can't set power",
				},
				{
					ID:    "ITB-1101-D1",
					Field: "volumes.headphones",
					Value: 30,
					Error: "can't set volumes",
				},
				{
					ID:    "ITB-1101-SW1",
					Field: "input.hdmiOutA.audio",
					Value: "1",
					Error: "can't set audio inputs",
				},
				{
					ID:    "ITB-1101-SW1",
					Field: "input.hdmiOutA.video",
					Value: "2",
					Error: "can't set video inputs",
				},
			},
		},
	},
	//	{
	//		name: "CantSetPower",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					PoweredOn: boolP(true),
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "poweredOn",
	//					Value: true,
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetBlank",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Blanked: boolP(true),
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "blanked",
	//					Value: true,
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetVolumes",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Volumes: map[string]int{"": 10},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "volumes",
	//					Value: map[string]int{"": 10},
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetMutes",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Mutes: map[string]bool{"": true},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "mutes",
	//					Value: map[string]bool{"": true},
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetAudioInputs",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{VideoInputs: map[string]string{"": "hdmi2"}},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							Audio: stringP("hdmi2"),
	//							Video: stringP("hdmi2"),
	//						},
	//						"other": {
	//							Audio: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							Video: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "input.$.audio",
	//					Value: map[string]api.Input{
	//						"": {
	//							Audio: stringP("hdmi2"),
	//							Video: stringP("hdmi2"),
	//						},
	//						"other": {
	//							Audio: stringP("hdmi2"),
	//						},
	//					},
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetVideoInputs",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{AudioVideoInputs: map[string]string{"": "hdmi2"}},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							AudioVideo: stringP("hdmi2"),
	//							Video:      stringP("hdmi2"),
	//						},
	//						"other": {
	//							Video: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							AudioVideo: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "input.$.video",
	//					Value: map[string]api.Input{
	//						"": {
	//							AudioVideo: stringP("hdmi2"),
	//							Video:      stringP("hdmi2"),
	//						},
	//						"other": {
	//							Video: stringP("hdmi2"),
	//						},
	//					},
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "CantSetAudioVideoInputs",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D4": &mock.Device{AudioInputs: map[string]string{"": "hdmi2"}},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							AudioVideo: stringP("hdmi2"),
	//							Audio:      stringP("hdmi2"),
	//						},
	//						"other": {
	//							AudioVideo: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D4": {
	//					Inputs: map[string]api.Input{
	//						"": {
	//							Audio: stringP("hdmi2"),
	//						},
	//					},
	//				},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D4",
	//					Field: "input.$.audioVideo",
	//					Value: map[string]api.Input{
	//						"": {
	//							AudioVideo: stringP("hdmi2"),
	//							Audio:      stringP("hdmi2"),
	//						},
	//						"other": {
	//							AudioVideo: stringP("hdmi2"),
	//						},
	//					},
	//					Error: "can't set this field on this device",
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "AudioInvalidBlockError",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D1": &mock.Device{
	//					Volumes: map[string]int{
	//						"": 30,
	//					},
	//					Mutes: map[string]bool{
	//						"": false,
	//					},
	//				},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D1": {
	//					Volumes: map[string]int{"invalid": 77},
	//					Mutes:   map[string]bool{"invalid": false},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D1": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D1",
	//					Field: "mutes.invalid",
	//					Value: false,
	//					Error: ErrInvalidBlock.Error(),
	//				},
	//				{
	//					ID:    "ITB-1101-D1",
	//					Field: "volumes.invalid",
	//					Value: int32(77),
	//					Error: ErrInvalidBlock.Error(),
	//				},
	//			},
	//		},
	//	},
	//	{
	//		name: "AudioSetError",
	//		driver: drivertest.Driver{
	//			Devices: map[string]drivers.Device{
	//				"ITB-1101-D1": &mock.Device{
	//					SetVolumeError: errors.New("no"),
	//					SetMuteError:   errors.New("i won't do it"),
	//					Volumes: map[string]int{
	//						"headphones": 30,
	//					},
	//					Mutes: map[string]bool{
	//						"headphones": false,
	//					},
	//				},
	//			},
	//		},
	//		req: api.StateRequest{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D1": {
	//					Volumes: map[string]int{"headphones": 77},
	//					Mutes:   map[string]bool{"headphones": false},
	//				},
	//			},
	//		},
	//		resp: api.StateResponse{
	//			Devices: map[api.DeviceID]api.DeviceState{
	//				"ITB-1101-D1": {},
	//			},
	//			Errors: []api.DeviceStateError{
	//				{
	//					ID:    "ITB-1101-D1",
	//					Field: "mutes.headphones",
	//					Value: false,
	//					Error: "i won't do it",
	//				},
	//				{
	//					ID:    "ITB-1101-D1",
	//					Field: "volumes.headphones",
	//					Value: int32(77),
	//					Error: "no",
	//				},
	//			},
	//		},
	//	},
}

func TestSetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range setTests {
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
			registry, err := drivers.New("../cmd/av-control-api/driver-config.yaml")
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

			// set the state of this room
			resp, err := gs.Set(ctx, room, tt.req)
			if tt.err != nil {
				is.True(err != nil)
				is.Equal(err.Error(), tt.err.Error())
			} else {
				is.NoErr(err)
				require.Equal(t, tt.resp, resp)
				is.Equal(resp, tt.resp)
			}
		})
	}
}

//func TestSetWrongDriver(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
//	defer cancel()
//
//	req := api.StateRequest{
//		Devices: map[api.DeviceID]api.DeviceState{
//			"ITB-1101-D1": {},
//		},
//	}
//	errWanted := errors.New("unknown driver: bad driver")
//
//	t.Run("", func(t *testing.T) {
//		room := api.Room{
//			Devices: make(map[api.DeviceID]api.Device),
//		}
//
//		apiDev := api.Device{
//			Address: "ITB-1101-D1",
//			Driver:  "bad driver",
//		}
//		room.Devices[api.DeviceID("ITB-1101-D1")] = apiDev
//
//		fakeDriver := drivertest.Driver{
//			Devices: map[string]drivers.Device{
//				"ITB-1101-D2": &mock.Device{},
//			},
//		}
//
//		server := drivertest.NewServer(fakeDriver.NewDeviceFunc())
//		conn, err := server.GRPCClientConn(ctx)
//		if err != nil {
//			t.Fatalf("unable to get grpc client connection: %s", err)
//		}
//
//		gs := &getSetter{
//			logger: zap.NewNop(),
//			drivers: map[string]drivers.DriverClient{
//				"": drivers.NewDriverClient(conn),
//			},
//		}
//
//		_, err = gs.Set(ctx, room, req)
//		if err != nil {
//			if diff := cmp.Diff(errWanted.Error(), err.Error()); diff != "" {
//				t.Fatalf("generated incorrect error (-want, +got):\n%s", diff)
//			}
//		}
//	})
//}
