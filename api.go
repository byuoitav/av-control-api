package avcontrol

import (
	"strings"
)

// StateRequest is the JSON object that a consumer of the av-control-api sends in a PUT request
// to set state
type StateRequest struct {
	Devices map[DeviceID]DeviceState `json:"devices,omitempty"`
}

// StateResponse is the JSON object that the API responds with when getting or setting state
type StateResponse struct {
	Devices map[DeviceID]DeviceState `json:"devices,omitempty"`
	Errors  []DeviceStateError       `json:"errors,omitempty"`
}

// DeviceState represents all of the possible fields that can can be set for a device.
type DeviceState struct {
	PoweredOn *bool `json:"poweredOn,omitempty"`
	Blanked   *bool `json:"blanked,omitempty"`

	Inputs  map[string]Input `json:"inputs,omitempty"`
	Volumes map[string]int   `json:"volumes,omitempty"`
	Mutes   map[string]bool  `json:"mutes,omitempty"`
}

// Input represents the current input state for a specific output on a device.
// Logically, Audio/Video will not be set if AudioVideo is set.
type Input struct {
	AudioVideo *string `json:"audioVideo,omitempty"`
	Audio      *string `json:"audio,omitempty"`
	Video      *string `json:"video,omitempty"`
}

// DeviceStateError is included in StateResponse whenever there is an error
// getting or setting a specific DeviceState field.
type DeviceStateError struct {
	// ID is the device that this error is associated with.
	ID DeviceID `json:"id"`

	// Field is the field (on DeviceState) that caused an error while the API was trying to get or set it.
	Field string `json:"field,omitempty"`

	// Value is the value the API was trying to set this field to when the error occurred.
	Value interface{} `json:"value,omitempty"`

	// Error is the error that happened.
	Error string `json:"error"`
}

// RoomHealth maps device id to device health status for each device in the room.
type RoomHealth struct {
	Devices map[DeviceID]DeviceHealth `json:"devices,omitempty"`
}

// DeviceHealth contains information about the device's health status
type DeviceHealth struct {
	Healthy *bool   `json:"healthy,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// RoomInfo maps device id to device info for each device in the room.
type RoomInfo struct {
	Devices map[DeviceID]DeviceInfo `json:"devices,omitempty"`
}

type DeviceInfo struct {
	Info  interface{} `json:"info,omitempty"`
	Error *string     `json:"error,omitempty"`
}

// DeviceID is a string in the format of Building-Room-DeviceName
type DeviceID string

// Room returns the Building-Room portion of a DeviceID
func (id DeviceID) Room() string {
	split := strings.SplitN(string(id), "-", 3)
	if len(split) != 3 {
		return string(id)
	}

	return split[0] + "-" + split[1]
}
