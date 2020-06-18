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
	Mapping api.DriverMapping
}

func (j *JustAddPowerRoom) DriverMapping(context.Context) (api.DriverMapping, error) {
	j.Mapping = api.DriverMapping{
		"Sony XBR": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8016",
				"k8s":     "http://sony-tv.service",
			},
		},
		"JAP3GRX": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8022",
				"k8s":     "http://just-add-power.service",
			},
		},
		"via-connect-pro": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8016",
				"k8s":     "http://via-service.service",
			},
		},
		"Shure Microphone": {
			BaseURLs: map[string]string{
				"default": "http://localhost:8014",
				"k8s":     "http://shure-service.service",
			},
		},
	}
	return j.Mapping, nil
}

func (j *JustAddPowerRoom) Room(context.Context, string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Devices: []api.Device{
			{
				ID:      "ITB-1101-D1",
				Address: "ITB-1101-D1.av",
				Type:    "Sony XBR",
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-D2",
				Address: "ITB-1101-D2.av",
				Type:    "Sony XBR",
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-D3",
				Address: "ITB-1101-D3.av",
				Type:    "Sony XBR",
				Ports: []api.Port{
					{
						Name: "hdmi!2",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-RX1",
				Address: "ITB-1101-RX1.av",
				Type:    "JAP3GRX",
				Ports: []api.Port{
					{
						Name: "rx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-RX2",
				Address: "ITB-1101-RX2.av",
				Type:    "JAP3GRX",
				Ports: []api.Port{
					{
						Name: "rx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-RX3",
				Address: "ITB-1101-RX3.av",
				Type:    "JAP3GRX",
				Ports: []api.Port{
					{
						Name: "rx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-NS1",
				Address: "0.0.0.0",
				Type:    "Aruba8PortNetworkSwitch",
				Ports: []api.Port{
					{
						Name: "1",
						Type: "audiovideo",
					},
					{
						Name: "2",
						Type: "audiovideo",
					},
					{
						Name: "3",
						Type: "audiovideo",
					},
					{
						Name: "4",
						Type: "audiovideo",
					},
					{
						Name: "5",
						Type: "audiovideo",
					},
					{
						Name: "6",
						Type: "audiovideo",
					},
					{
						Name: "7",
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX1",
				Address: "10.66.76.185",
				Type:    "JAP3GTX",
				Ports: []api.Port{
					{
						Name: "tx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX2",
				Address: "10.66.76.186",
				Type:    "JAP3GTX",
				Ports: []api.Port{
					{
						Name: "tx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX3",
				Address: "10.66.76.187",
				Type:    "JAP3GTX",
				Ports: []api.Port{
					{
						Name: "tx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-TX4",
				Address: "10.66.76.188",
				Type:    "JAP3GTX",
				Ports: []api.Port{
					{
						Name: "tx",
						Type: "audiovideo",
					},
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:   "ITB-1101-HDMI1",
				Type: "non-controllable",
				Ports: []api.Port{
					{
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
						Type: "audiovideo",
					},
				},
			},
			{
				ID:   "ITB-1101-PC1",
				Type: "non-controllable",
				Ports: []api.Port{
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:   "ITB-1101-SIGN1",
				Type: "non-conrtollable",
				Ports: []api.Port{
					{
						Type: "audiovideo",
					},
				},
			},
			{
				ID:      "ITB-1101-DSP1",
				Address: "ITB-1101-DSP1.av",
				Type:    "QSC-Core-110F",
				Ports: []api.Port{
					{
						Name: "Mic1",
						Type: "audio",
					},
					{
						Name: "Mic2",
						Type: "audio",
					},
					{
						Name: "Mic3",
						Type: "audio",
					},
				},
			},
			{
				ID:   "ITB-1101-MIC1",
				Type: "Shure Microphone",
				Ports: []api.Port{
					{
						Name: "Mic1",
						Type: "audio",
					},
				},
			},
			{
				ID:   "ITB-1101-MIC2",
				Type: "Shure Microphone",
				Ports: []api.Port{
					{
						Name: "Mic2",
						Type: "audio",
					},
				},
			},
			{
				ID:   "ITB-1101-MIC3",
				Type: "Shure Microphone",
				Ports: []api.Port{
					{
						Name: "Mic3",
						Type: "audio",
					},
				},
			},
		},
	}, nil
}
