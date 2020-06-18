package mock

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

// SimpleRoom implements api.DeviceService and represents a simple room with the following devices:
//
// * a tv
// * a HDMI input
// * a VIA input
//
// The audio for this room is coming directly out of the TV,
// and the inputs are going directly into hdmi 1 & 2 on the TV.
// the via is also capable of getting/setting it's own volume/mute.
type SimpleRoom struct {
	Mapping api.DriverMapping
}

func (s *SimpleRoom) DriverMapping(context.Context) (api.DriverMapping, error) {
	s.Mapping = api.DriverMapping{
		"Sony XBR": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8016",
				"k8s":     "http://sony-tv.service",
			},
		},
		"via-connect-pro": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8012",
				"k8s":     "http://via-service.service",
			},
		},
	}
	return s.Mapping, nil
}

func (s *SimpleRoom) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type:    "Sony XBR",
				Ports: []api.Port{
					{
						Name: "hdmi!1",
						Type: "audiovideo",
					},
					{
						Name: "hdmi!2",
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
