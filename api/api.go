package api

import (
	"strings"
)

type StateRequest struct {
	Devices map[DeviceID]DeviceState `json:"devices,omitempty"`
}

type StateResponse struct {
	Devices map[DeviceID]DeviceState `json:"devices,omitempty"`
	Errors  []DeviceStateError       `json:"errors,omitempty"`
}

/*
{
	"devices": {
		"poweredOn": true,
		"inputs": {
			"output": {
				"audioVideo": "input",
				"audio": "input",
				"video": "input",
			},
			"output2": {
				"audioVideo": "input",
				"audio": "input",
				"video": "input",
			}
		},
		"blanked": true,
		"volumes": {
			"001100": 50,
			"001200": 30,
			"TM CB1170MediaGain": 30
		},
		"mutes": {
			"out1Mute": false
		}
	}
}
*/

type DeviceState struct {
	PoweredOn *bool `json:"poweredOn,omitempty"`
	Blanked   *bool `json:"blanked,omitempty"`

	Input   map[string]Input `json:"inputs,omitempty"`
	Volumes map[string]int   `json:"volumes,omitempty"`
	Mutes   map[string]bool  `json:"mutes,omitempty"`
}

type Input struct {
	AudioVideo *string `json:"audiovideo,omitempty"`
	Audio      *string `json:"audio,omitempty"`
	Video      *string `json:"video,omitempty"`
	// TODO ?
	// AvailableInputs []DeviceID `json:"availableInputs,omitempty"`
}

// for option 2
//type VolumeState struct {
//	Outputs map[string]int
//	Inputs  map[string]int
//}
//
//type MutedState struct {
//	Outputs map[string]bool
//	Inputs  map[string]bool
//}

type DeviceStateError struct {
	ID    DeviceID    `json:"id"`
	Field string      `json:"field"`
	Value interface{} `json:"value,omitempty"`
	Error string      `json:"error"`
}

type DeviceID string

func (id DeviceID) Room() string {
	split := strings.SplitN(string(id), "-", 3)
	if len(split) != 3 {
		return string(id)
	}

	return split[0] + "-" + split[1]
}
