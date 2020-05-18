package api

import "strings"

type StateRequest struct {
	Devices map[DeviceID]DeviceState `json:"devices"`
}

type StateResponse struct {
	Devices map[DeviceID]DeviceState `json:"devices,omitempty"`
	Errors  []DeviceStateError       `json:"errors,omitempty"`
}

type DeviceState struct {
	PoweredOn *bool     `json:"poweredOn,omitempty"`
	Input     *DeviceID `json:"input,omitempty"`
	Blanked   *bool     `json:"blanked,omitempty"`
	Volume    *int      `json:"volume,omitempty"`
	Muted     *bool     `json:"muted,omitempty"`
}

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
