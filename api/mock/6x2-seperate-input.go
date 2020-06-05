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
type SixTwoSeparateInput struct{}

func (SixTwoSeparateInput) Room(context.Context, string) ([]api.Device, error) {
	return []api.Device{
		api.Device{
			ID:      "ITB-1101-D1",
			Address: "ITB-1101-D1.av",
			Type: api.DeviceType{
				ID: "Projector",
				Commands: map[string]api.Command{
					"SetPower": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetPower/{{power}}",
						},
						Order: intP(0),
					},
					"GetPower": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetPower",
						},
					},
					"SetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetBlanked/{{blanked}}",
						},
						Order: intP(0),
					},
					"GetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetBlanked",
						},
					},
					"SetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetAVInput/{{port}}",
						},
					},
					"GetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetAVInput",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "hdbaset",
					Endpoints: api.Endpoints{
						"ITB-1101-SW1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
			},
		},
		api.Device{
			ID:      "ITB-1101-D2",
			Address: "ITB-1101-D2.av",
			Type: api.DeviceType{
				ID: "Projector",
				Commands: map[string]api.Command{
					"SetPower": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetPower/{{power}}",
						},
						Order: intP(0),
					},
					"GetPower": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetPower",
						},
					},
					"SetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetBlanked/{{blanked}}",
						},
						Order: intP(0),
					},
					"GetBlanked": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetBlanked",
						},
					},
					"SetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetAVInput/{{port}}",
						},
					},
					"GetAVInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetAVInput",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "hdbaset",
					Endpoints: api.Endpoints{
						"ITB-1101-SW1",
					},
					Type:     "audiovideo",
					Incoming: true,
				},
			},
		},
		api.Device{
			ID: "ITB-1101-AUD1",
			Type: api.DeviceType{
				ID: "non-controllable",
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
			ID: "ITB-1101-AUD2",
			Type: api.DeviceType{
				ID: "non-controllable",
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
			ID:      "ITB-1101-SW1",
			Address: "ITB-1101-SW1.av",
			Type: api.DeviceType{
				ID: "6x1",
				Commands: map[string]api.Command{
					"SetAudioInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetAudioInput/{{port}}",
						},
					},
					"GetAudioInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetAudioInput",
						},
					},
					"SetVideoInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetVideoInput/{{port}}",
						},
					},
					"GetVideoInput": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetVideoInput",
						},
					},
				},
			},
			Ports: []api.Port{
				api.Port{
					Name: "videoOut1",
					Endpoints: api.Endpoints{
						"ITB-1101-D1",
					},
					Type: "audiovideo",
				},
				api.Port{
					Name: "videoOut2",
					Endpoints: api.Endpoints{
						"ITB-1101-D2",
					},
					Type: "audiovideo",
				},
				api.Port{
					Name: "audioOut1",
					Endpoints: api.Endpoints{
						"ITB-1101-AUD1",
					},
					Type: "audio",
				},
				api.Port{
					Name: "audioOut2",
					Endpoints: api.Endpoints{
						"ITB-1101-AUD2",
					},
					Type: "audio",
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
			ID:      "ITB-1101-VIA1",
			Address: "ITB-1101-VIA1.av",
			Type: api.DeviceType{
				ID: "via-connect-pro",
				Commands: map[string]api.Command{
					"GetVolume": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetVolume",
						},
					},
					"GetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/GetMuted",
						},
					},
					"SetVolume": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetVolume/{{level}}",
						},
					},
					"SetMuted": api.Command{
						URLs: map[string]string{
							"default": "http://ITB-1101-CP1.byu.edu/{{address}}/SetMuted/{{muted}}",
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
