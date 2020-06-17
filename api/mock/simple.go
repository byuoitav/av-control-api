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

func (s *SimpleRoom) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type: api.DeviceType{
					ID: "Sony XBR",
					Commands: map[string]api.Command{
						"SetPower": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetPower/{{power}}",
							},
							Order: intP(0),
						},
						"GetPower": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetPower",
							},
						},
						"GetBlanked": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetBlanked",
							},
						},
						"SetBlanked": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetBlanked/{{blanked}}",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetAVInput/{{port}}",
							},
						},
						"GetAVInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetAVInput",
							},
						},
						"GetVolume": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "hdmi!1",
						Endpoints: api.Endpoints{
							"ITB-1101-VIA1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Name: "hdmi!2",
						Endpoints: api.Endpoints{
							"ITB-1101-HDMI1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-VIA1",
				Address: "ITB-1101-VIA1.av",
				Type: api.DeviceType{
					ID: "via-connect-pro",
					Commands: map[string]api.Command{
						"GetVolume": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "",
						Endpoints: api.Endpoints{
							"ITB-1101-D1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID: "ITB-1101-HDMI1",
				Type: api.DeviceType{
					ID: "hdmi-input",
				},
				Ports: []api.Port{
					{
						Name: "",
						Endpoints: api.Endpoints{
							"ITB-1101-D1",
						},
						Type: "audiovideo",
					},
				},
			},
		},
	}, nil
}
