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
	BaseURL string
}

func (s *SimpleRoom) SetBaseURL(baseURL string) {
	s.BaseURL = baseURL
}

func (s *SimpleRoom) Room(context.Context, string) ([]api.Device, error) {
	return []api.Device{
		api.Device{
			ID:      "ITB-1101-D1",
			Address: "ITB-1101-D1.av",
			Type: api.DeviceType{
				ID: "Sony XBR",
				Commands: map[string]api.Command{
					"SetPower": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetPower/{{power}}",
						},
						Order: intP(0),
					},
					"GetPower": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetPower",
						},
					},
					"SetAVInput": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetAVInput/{{port}}",
						},
					},
					"GetAVInput": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetAVInput",
						},
					},
					"GetVolume": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetVolume",
						},
					},
					"GetMuted": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetMuted",
						},
					},
					"SetVolume": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetVolume/{{level}}",
						},
					},
					"SetMuted": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetMuted/{{muted}}",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "hdmi!1",
					Endpoints: api.Endpoints{
						"ITB-1101-VIA1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
				api.Port{
					Name: "hdmi!2",
					Endpoints: api.Endpoints{
						"ITB-1101-HDMI1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
			},
		},
		api.Device{
			ID:      "ITB-1101-VIA1",
			Address: "ITB-1101-VIA1.av",
			Type: api.DeviceType{
				ID: "via-connect-pro",
				Commands: map[string]api.Command{
					"GetVolume": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetVolume",
						},
					},
					"GetMuted": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/GetMuted",
						},
					},
					"SetVolume": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetVolume/{{level}}",
						},
					},
					"SetMuted": api.Command{
						URLs: map[string]string{
							"default": s.BaseURL + "/{{address}}/SetMuted/{{muted}}",
						},
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
