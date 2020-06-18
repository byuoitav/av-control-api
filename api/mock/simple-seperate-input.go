package mock

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

// SimpleSeparateInput implements api.DeviceService and represents a simple room with the following devices:
//
// * a tv
// * a 4x1 videoswitcher that can separate video/audio input
// * an amp that we can control volume on
// * a HDMI input
// * a VIA input
//
// In this room, the audio is being split off of the HDMI out of the 4x1
// the video is going to the TV, and the audio is going to a controllable amp and into the speakers.
type SimpleSeparateInput struct {
	Mapping api.DriverMapping
}

func (s *SimpleSeparateInput) DriverMapping(context.Context) (api.DriverMapping, error) {
	s.Mapping = api.DriverMapping{
		"TV No Audio": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8016",
				"k8s":     "http://sony-tv.service",
			},
		},
		"4x1": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8002",
				"k8s":     "http://atlona-4x1.service",
			},
		},
		"amp": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8069",
				"k8s":     "http://amp.service",
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

func (s *SimpleSeparateInput) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type:    "TV No Audio",
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-SW1",
				Address: "ITB-1101-SW1.av",
				Type:    "4x1",
				Ports: []api.Port{
					{
						Name: "1",
						Type: "audiovideo",
					},
					{
						Name: "1",
						Type: "audiovideo",
					},
					{
						Name: "2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-AMP1",
				Address: "ITB-1101-AMP1.av",
				Type:    "amp",
				Ports: []api.Port{
					{
						Name: "",
						Type: "audio",
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
