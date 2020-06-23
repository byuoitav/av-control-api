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

type DeviceStateError struct {
	ID    DeviceID    `json:"id"`
	Field string      `json:"field,omitempty"`
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
