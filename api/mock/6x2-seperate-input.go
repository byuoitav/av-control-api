package mock

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

// SixTwoSeparateInput implements api.DeviceService and represents a room with the following devices:
//
// * 2 projectors
// * a 6x2 videoswitcher that can separate video/audio input
// * overhead speakers
// * a HDMI input
// * a VIA input
// * two mic's
//
// In this room, the audio is coming off of the 6x2 audio outputs -> speakers
// mics are going into audioIn1 and audioIn2 on the 6x2
type SixTwoSeparateInput struct {
	Mapping api.DriverMapping
}

func (s *SixTwoSeparateInput) DriverMapping(context.Context) (api.DriverMapping, error) {
	s.Mapping = api.DriverMapping{
		"Projector": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8017",
				"k8s":     "http://projector-service.service",
			},
		},
		"6x2": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8020",
				"k8s":     "http://6x2-service.service",
			},
		},
		"via-connect-pro": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8012",
				"k8s":     "http://via-service.service",
			},
		},
		"hdmi-input": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8420",
				"k8s":     "http://hdmi-service.service",
			},
		},
	}
	return s.Mapping, nil
}

func (s *SixTwoSeparateInput) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type:    "Projector",
				Ports: []api.Port{
					{
						Name: "hdbaset",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-D2",
				Address: "ITB-1101-D2.av",
				Type:    "Projector",
				Ports: []api.Port{
					{
						Name: "hdbaset",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:   "ITB-1101-AUD1",
				Type: "non-controllable",
				Ports: []api.Port{
					{
						Name: "",
						Type: "audio",
					},
				},
			},
			{
				ID:   "ITB-1101-AUD2",
				Type: "non-controllable",
				Ports: []api.Port{
					{
						Name: "",
						Type: "audio",
					},
				},
			},
			{
				ID:      "ITB-1101-SW1",
				Address: "ITB-1101-SW1.av",
				Type:    "6x2",
				Ports: []api.Port{
					{
						Name: "videoOut1",
						Type: "audiovideo",
					},
					{
						Name: "videoOut2",
						Type: "audiovideo",
					},
					{
						Name: "audioOut1",
						Type: "audio",
					},
					{
						Name: "audioOut2",
						Type: "audio",
					},
					{
						Name: "2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-VIA1",
				Address: "ITB-1101-VIA1.av",
				Type:    "via-connect-pro",
				Ports: []api.Port{
					{
						Name: "",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:   "ITB-1101-HDMI1",
				Type: "hdmi-input",
				Ports: []api.Port{
					{
						Name: "",
						Type: "audiovideo",
					},
				},
			},
		},
	}, nil
}
