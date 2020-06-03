package api

import (
	"encoding/json"
	"strings"
)

type StateRequest struct {
	OutputGroups map[DeviceID]OutputGroupState `json:"outputGroups,omitempty"`
}

type StateResponse struct {
	OutputGroups map[DeviceID]OutputGroupState `json:"outputGroups,omitempty"`
	Errors       []DeviceStateError            `json:"errors,omitempty"`
}

type OutputGroupState struct {
	PoweredOn *bool  `json:"poweredOn,omitempty"`
	Input     *Input `json:"input,omitempty"`
	Blanked   *bool  `json:"blanked,omitempty"`
	Volume    *int   `json:"volume,omitempty"`
	Muted     *bool  `json:"muted,omitempty"`

	Outputs map[DeviceID]OutputState `json:"outputs,omitempty"`
}

type OutputState struct {
	PoweredOn *bool  `json:"poweredOn,omitempty"`
	Input     *Input `json:"input,omitempty"`
	Blanked   *bool  `json:"blanked,omitempty"`
	Volume    *int   `json:"volume,omitempty"`
	Muted     *bool  `json:"muted,omitempty"`
}

type Input struct {
	Audio            *DeviceID  `json:"audio,omitempty"`
	Video            *DeviceID  `json:"video,omitempty"`
	CanSetSeparately *bool      `json:"canSetSeparately,omitempty"`
	AvailableInputs  []DeviceID `json:"availableInputs,omitempty"`
}

func (i Input) JSONMarshal() ([]byte, error) {
	if i.Audio != nil && i.Video != nil {
		return json.Marshal(i)
	}

	i.CanSetSeparately = nil
	return json.Marshal(i)
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
