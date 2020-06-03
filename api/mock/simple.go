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
type SimpleRoom struct{}

func (SimpleRoom) Room(context.Context, string) ([]api.Device, error) {
	return []api.Device{
		api.Device{
			ID: "ITB-1101-D1",
			Type: api.DeviceType{
				ID: "Sony XBR",
				Commands: map[string]api.Command{
					"": api.Command{
						URLs:  map[string]string{},
						Order: intP(1),
					},
				},
			},
			Ports: []api.Port{},
		},
		api.Device{
			ID: "ITB-1101-VIA1",
			Type: api.DeviceType{
				ID: "via-connect-pro",
				Commands: map[string]api.Command{
					"GetVolume": api.Command{
						URLs: map[string]string{
							"default": "get volume for :address",
						},
						Order: intP(1),
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "",
					Endpoints: api.Endpoints{
						"ITB-1101-D1",
					},
					Type: "audiovideo",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-HDMI1",
			Type: api.DeviceType{
				ID: "hdmi-input",
			},
			Ports: []api.Port{
				api.Port{
					Name: "",
					Endpoints: api.Endpoints{
						"ITB-1101-D1",
					},
					Type: "audiovideo",
				},
			},
		},
	}, nil
}
