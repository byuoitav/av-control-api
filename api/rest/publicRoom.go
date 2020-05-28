package rest

//PublicRoom is the struct that is returned (or put) as part of the public API
type PublicRoom struct {
	Building          string        `json:"-"`
	Room              string        `json:"-"`
	CurrentVideoInput string        `json:"currentVideoInput,omitempty"`
	CurrentAudioInput string        `json:"currentAudioInput,omitempty"`
	Power             string        `json:"power,omitempty"`
	Blanked           *bool         `json:"blanked,omitempty"`
	Muted             *bool         `json:"muted,omitempty"`
	Volume            *int          `json:"volume,omitempty"`
	Displays          []Display     `json:"displays,omitempty"`
	AudioDevices      []AudioDevice `json:"audioDevices,omitempty"`
}

//Device is a struct for inheriting
type Device struct {
	Name  string `json:"name,omitempty"`
	Power string `json:"power,omitempty"`
	Input string `json:"input,omitempty"`
}

//AudioDevice represents an audio device
type AudioDevice struct {
	Device
	Muted  *bool `json:"muted,omitempty"`
	Volume *int  `json:"volume,omitempty"`
}

//Display represents a display
type Display struct {
	Device
	Blanked *bool `json:"blanked,omitempty"`
}
