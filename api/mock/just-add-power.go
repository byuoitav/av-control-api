package mock

import (
	"context"

	"github.com/byuoitav/av-control-api/api"
)

// JustAddPowerRoom implements api.DeviceService and represents a room with the following devices:
//
// * 3 tvs
// * a HDMI input
// * a VIA input
// * a PC input
// * 3 Just Add Power receivers
// * 4 Just Add Power transmitters
// * a network switch
// * a Digital Sign... thingy
// * a DSP
// * 3 Shure microphones
//
// The audio for this room is coming directly out of the TVs,
// and the inputs are going directly into hdmi 2 on the TV.
// Each TV goes to a reciever which goes to a transmitter and then a final input.
// the via is also capable of getting/setting it's own volume/mute.
type JustAddPowerRoom struct {
	BaseURL string
}

func (j *JustAddPowerRoom) SetBaseURL(baseURL string) {
	j.BaseURL = baseURL
}

func (j *JustAddPowerRoom) Room(context.Context, string) (api.Room, error) {
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
								"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
							},
							Order: intP(0),
						},
						"GetPower": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetPower",
							},
						},
						"GetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetBlanked",
							},
						},
						"SetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetBlanked",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetAVInput/{{port}}",
							},
						},
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetAVInput",
							},
						},
						"GetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Endpoints: api.Endpoints{
							"ITB-1101-RX1",
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
					ID: "Sony XBR",
					Commands: map[string]api.Command{
						"SetPower": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
							},
							Order: intP(0),
						},
						"GetPower": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetPower",
							},
						},
						"GetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetBlanked",
							},
						},
						"SetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetBlanked",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetAVInput/{{port}}",
							},
						},
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetAVInput",
							},
						},
						"GetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Endpoints: api.Endpoints{
							"ITB-1101-RX2",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-D3",
				Address: "ITB-1101-D3.av",
				Type: api.DeviceType{
					ID: "Sony XBR",
					Commands: map[string]api.Command{
						"SetPower": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
							},
							Order: intP(0),
						},
						"GetPower": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetPower",
							},
						},
						"GetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetBlanked",
							},
						},
						"SetBlanked": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetBlanked",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetAVInput/{{port}}",
							},
						},
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetAVInput",
							},
						},
						"GetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Endpoints: api.Endpoints{
							"ITB-1101-RX3",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-RX1",
				Address: "ITB-1101-RX1.av",
				Type: api.DeviceType{
					ID: "JAP3GRX",
					Commands: map[string]api.Command{
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetStream",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{receiverAddr}}/SetStream/{{transmitterAddr}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "rx",
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Incoming: true,
						Type:     "audiovideo",
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-D1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-RX2",
				Address: "ITB-1101-RX2.av",
				Type: api.DeviceType{
					ID: "JAP3GRX",
					Commands: map[string]api.Command{
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetStream",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{receiverAddr}}/SetStream/{{transmitterAddr}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "rx",
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Incoming: true,
						Type:     "audiovideo",
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-D2",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-RX3",
				Address: "ITB-1101-RX3.av",
				Type: api.DeviceType{
					ID: "JAP3GRX",
					Commands: map[string]api.Command{
						"GetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetStream",
							},
						},
						"SetAVInput": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{receiverAddr}}/SetStream/{{transmitterAddr}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "rx",
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Incoming: true,
						Type:     "audiovideo",
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-D3",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-NS1",
				Address: "0.0.0.0",
				Type: api.DeviceType{
					ID:       "Aruba8PortNetworkSwitch",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Name: "1",
						Endpoints: api.Endpoints{
							"ITB-1101-RX1",
						},
						Type: "audiovideo",
					},
					{
						Name: "2",
						Endpoints: api.Endpoints{
							"ITB-1101-RX2",
						},
						Type: "audiovideo",
					},
					{
						Name: "3",
						Endpoints: api.Endpoints{
							"ITB-1101-RX3",
						},
						Type: "audiovideo",
					},
					{
						Name: "4",
						Endpoints: api.Endpoints{
							"ITB-1101-TX1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Name: "5",
						Endpoints: api.Endpoints{
							"ITB-1101-TX2",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Name: "6",
						Endpoints: api.Endpoints{
							"ITB-1101-TX3",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Name: "7",
						Endpoints: api.Endpoints{
							"ITB-1101-TX4",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
				},
			},
			{
				ID:      "ITB-1101-TX1",
				Address: "10.66.76.185",
				Type: api.DeviceType{
					ID:       "JAP3GTX",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Name: "tx",
						Endpoints: api.Endpoints{
							"ITB-1101-HDMI1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX2",
				Address: "10.66.76.186",
				Type: api.DeviceType{
					ID:       "JAP3GTX",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Name: "tx",
						Endpoints: api.Endpoints{
							"ITB-1101-VIA1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX3",
				Address: "10.66.76.187",
				Type: api.DeviceType{
					ID:       "JAP3GTX",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Name: "tx",
						Endpoints: api.Endpoints{
							"ITB-1101-PC1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX4",
				Address: "10.66.76.188",
				Type: api.DeviceType{
					ID:       "JAP3GTX",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Name: "tx",
						Endpoints: api.Endpoints{
							"ITB-1101-SIGN1",
						},
						Type:     "audiovideo",
						Incoming: true,
					},
					{
						Endpoints: api.Endpoints{
							"ITB-1101-NS1",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID: "ITB-1101-HDMI1",
				Type: api.DeviceType{
					ID:       "non-controllable",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Endpoints: api.Endpoints{
							"ITB-1101-TX1",
						},
						Type: "audiovideo",
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
								"default": j.BaseURL + "/{{address}}/GetVolume",
							},
						},
						"GetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/GetMuted",
							},
						},
						"SetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetVolume/{{level}}",
							},
						},
						"SetMuted": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/SetMuted/{{muted}}",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Endpoints: api.Endpoints{
							"ITB-1101-TX2",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID: "ITB-1101-PC1",
				Type: api.DeviceType{
					ID:       "non-controllable",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Endpoints: api.Endpoints{
							"ITB-1101-TX3",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID: "ITB-1101-SIGN1",
				Type: api.DeviceType{
					ID:       "non-controllable",
					Commands: make(map[string]api.Command),
				},
				Ports: []api.Port{
					{
						Endpoints: api.Endpoints{
							"ITB-1101-TX4",
						},
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-DSP1",
				Address: "ITB-1101-DSP1.av",
				Type: api.DeviceType{
					ID: "QSC-Core-110F",
					Commands: map[string]api.Command{
						"SetVolume": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/{{input}}/volume/set/{{level}}",
							},
						},
						"GetVolumeByBlock": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/{{block}}/volume/level",
							},
						},
						"Mute": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/{{input}}/volume/mute",
							},
						},
						"UnMute": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/{{input}}/volume/unmute",
							},
						},
						"GetMutedByBlock": {
							URLs: map[string]string{
								"default": j.BaseURL + "/{{address}}/{{block}}/mute/status",
							},
						},
					},
				},
				Ports: []api.Port{
					{
						Name: "Mic1",
						Endpoints: api.Endpoints{
							"ITB-1101-MIC1",
						},
						Type: "audio",
						// Incoming: true,
					},
					{
						Name: "Mic2",
						Endpoints: api.Endpoints{
							"ITB-1101-MIC2",
						},
						Type: "audio",
						// Incoming: true,
					},
					{
						Name: "Mic3",
						Endpoints: api.Endpoints{
							"ITB-1101-MIC3",
						},
						Type: "audio",
						// Incoming: true,
					},
				},
			},
			{
				ID: "ITB-1101-MIC1",
				Type: api.DeviceType{
					ID:       "Shure Microphone",
					Commands: map[string]api.Command{
						// "GetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/GetPower",
						// 	},
						// },
						// "SetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
						// 	},
						// },
					},
				},
				Ports: []api.Port{
					{
						Name: "Mic1",
						Endpoints: api.Endpoints{
							"ITB-1101-DSP1",
						},
						Type:     "audio",
						Incoming: true,
					},
				},
			},
			{
				ID: "ITB-1101-MIC2",
				Type: api.DeviceType{
					ID:       "Shure Microphone",
					Commands: map[string]api.Command{
						// "GetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/GetPower",
						// 	},
						// },
						// "SetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
						// 	},
						// },
					},
				},
				Ports: []api.Port{
					{
						Name: "Mic2",
						Endpoints: api.Endpoints{
							"ITB-1101-DSP1",
						},
						Type:     "audio",
						Incoming: true,
					},
				},
			},
			{
				ID: "ITB-1101-MIC3",
				Type: api.DeviceType{
					ID:       "Shure Microphone",
					Commands: map[string]api.Command{
						// "GetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/GetPower",
						// 	},
						// },
						// "SetPower": {
						// 	URLs: map[string]string{
						// 		"default": j.BaseURL + "/{{address}}/SetPower/{{power}}",
						// 	},
						// },
					},
				},
				Ports: []api.Port{
					{
						Name: "Mic3",
						Endpoints: api.Endpoints{
							"ITB-1101-DSP1",
						},
						Type:     "audio",
						Incoming: true,
					},
				},
			},
		},
	}, nil
}
