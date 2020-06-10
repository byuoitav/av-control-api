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
	BaseURL string
}

func (s *SimpleSeparateInput) SetBaseURL(baseURL string) {
	s.BaseURL = baseURL
}

func (s *SimpleSeparateInput) Room(context.Context, string) (api.Room, error) {
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
						"SetBlanked": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetBlanked/{{blanked}}",
							},
							Order: intP(0),
						},
						"GetBlanked": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetBlanked",
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
					},
				},
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Endpoints: api.Endpoints{
							"ITB-1101-SW1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-SW1",
				Address: "ITB-1101-SW1.av",
				Type: api.DeviceType{
					ID: "4x1",
					Commands: map[string]api.Command{
						"SetAudioInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetAudioInput/{{port}}",
							},
						},
						"GetAudioInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetAudioInput",
							},
						},
						"SetVideoInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/SetVideoInput/{{port}}",
							},
						},
						"GetVideoInput": {
							URLs: map[string]string{
								"default": s.BaseURL + "/{{address}}/GetVideoInput",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "1",
						Endpoints: api.Endpoints{
							"ITB-1101-D1",
							"ITB-1101-AMP1",
						},
						Type:     "audiovideo",
						Incoming: false,
					},
					{
						Name: "1",
						Endpoints: api.Endpoints{
							"ITB-1101-VIA1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Name: "2",
						Endpoints: api.Endpoints{
							"ITB-1101-HDMI1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-AMP1",
				Address: "ITB-1101-AMP1.av",
				Type: api.DeviceType{
					ID: "amp",
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
							"ITB-1101-SW1",
						},
						Type:     "audio",
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
							"ITB-1101-SW1",
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
							"ITB-1101-SW1",
						},
						Type: "audiovideo",
					},
				},
			},
		},
	}, nil
}
