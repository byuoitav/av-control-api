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
type SimpleSeparateInput struct{}

func (SimpleSeparateInput) Room(context.Context, string) ([]api.Device, error) {
	return []api.Device{
		api.Device{
			ID:      "ITB-1101-D1",
			Address: "ITB-1101-D1.av",
			Type: api.DeviceType{
				ID: "Sony XBR",
				Commands: map[string]api.Command{
					"SetPower": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetPower/{{power}}",
						},
						Order: intP(0),
					},
					"GetPower": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetPower",
						},
					},
					"SetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetBlanked/{{blanked}}",
						},
						Order: intP(0),
					},
					"GetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetBlanked",
						},
					},
					"SetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetAVInput/{{port}}",
						},
					},
					"GetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetAVInput",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "hdmi!2",
					Endpoints: api.Endpoints{
						"ITB-1101-SW1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
			},
		},
		api.Device{
			ID:      "ITB-1101-SW1",
			Address: "ITB-1101-SW1.av",
			Type: api.DeviceType{
				ID: "4x1",
				Commands: map[string]api.Command{
					"SetAudioInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetAudioInput/{{port}}",
						},
					},
					"GetAudioInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetAudioInput",
						},
					},
					"SetVideoInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetVideoInput/{{port}}",
						},
					},
					"GetVideoInput": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetVideoInput",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "1",
					Endpoints: api.Endpoints{
						"ITB-1101-D1",
						"ITB-1101-AMP1",
					},
					Type: "audiovideo",
				},
				api.Port{
					Name: "1",
					Endpoints: api.Endpoints{
						"ITB-1101-VIA1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
				api.Port{
					Name: "2",
					Endpoints: api.Endpoints{
						"ITB-1101-HDMI1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
			},
		},
		api.Device{
			ID:      "ITB-1101-AMP1",
			Address: "ITB-1101-AMP1.av",
			Type: api.DeviceType{
				ID: "amp",
				Commands: map[string]api.Command{
					"GetVolume": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetVolume",
						},
					},
					"GetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetMuted",
						},
					},
					"SetVolume": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetVolume/{{level}}",
						},
					},
					"SetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetMuted/{{muted}}",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "",
					Endpoints: api.Endpoints{
						"ITB-1101-SW1",
					},
					Type:     "audio",
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
							"default": "http://{{address}}/GetVolume",
						},
					},
					"GetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/GetMuted",
						},
					},
					"SetVolume": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetVolume/{{level}}",
						},
					},
					"SetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://{{address}}/SetMuted/{{muted}}",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "",
					Endpoints: api.Endpoints{
						"ITB-1101-SW1",
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
						"ITB-1101-SW1",
					},
					Type: "audiovideo",
				},
			},
		},
	}, nil
}
