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
	BaseURL string
}

func (s *SixTwoSeparateInput) SetBaseURL(baseURL string) {
	s.BaseURL = baseURL
}

func (s *SixTwoSeparateInput) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type: api.DeviceType{
					ID: "Projector",
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
						Name: "hdbaset",
						Endpoints: api.Endpoints{
							"ITB-1101-SW1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-D2",
				Address: "ITB-1101-D2.av",
				Type: api.DeviceType{
					ID: "Projector",
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
						Name: "hdbaset",
						Endpoints: api.Endpoints{
							"ITB-1101-SW1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID: "ITB-1101-AUD1",
				Type: api.DeviceType{
					ID: "non-controllable",
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
				ID: "ITB-1101-AUD2",
				Type: api.DeviceType{
					ID: "non-controllable",
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
				ID:      "ITB-1101-SW1",
				Address: "ITB-1101-SW1.av",
				Type: api.DeviceType{
					ID: "6x1",
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
						Name: "videoOut1",
						Endpoints: api.Endpoints{
							"ITB-1101-D1",
						},
						Type: "audiovideo",
					},
					{
						Name: "videoOut2",
						Endpoints: api.Endpoints{
							"ITB-1101-D2",
						},
						Type: "audiovideo",
					},
					{
						Name: "audioOut1",
						Endpoints: api.Endpoints{
							"ITB-1101-AUD1",
						},
						Type: "audio",
					},
					{
						Name: "audioOut2",
						Endpoints: api.Endpoints{
							"ITB-1101-AUD2",
						},
						Type: "audio",
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
