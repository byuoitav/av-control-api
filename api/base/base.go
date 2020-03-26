package base

import (
	ei "github.com/byuoitav/common/v2/events"
)

//ActionStructure is the internal struct we use to pass commands around once
//they've been evaluated.
//also contains a list of Events to be published
type ActionStructure struct {
	Action              string             `json:"action"`
	GeneratingEvaluator string             `json:"generatingEvaluator"`
	Device              Device             `json:"device"`
	DestinationDevice   DestinationDevice  `json:"destination_device"`
	Parameters          map[string]string  `json:"parameters"`
	DeviceSpecific      bool               `json:"deviceSpecific,omitempty"`
	Overridden          bool               `json:"overridden"`
	EventLog            []ei.Event         `json:"events"`
	Children            []*ActionStructure `json:"children"`
	Callback            func(StatusPackage, chan<- StatusPackage) error
}

// DestinationDevice represents the device that is being acted upon.
type DestinationDevice struct {
	Device
	AudioDevice bool `json:"audio"`
	Display     bool `json:"video"`
}

// StatusPackage contains the callback information for the action.
type StatusPackage struct {
	Key    string
	Value  interface{}
	Device Device
	Dest   DestinationDevice
}

//Equals checks if the action structures are equal
func (a *ActionStructure) Equals(b ActionStructure) bool {
	return a.Action == b.Action &&
		a.Device.ID == b.Device.ID &&
		a.Device.Address == b.Device.Address &&
		a.DeviceSpecific == b.DeviceSpecific &&
		a.Overridden == b.Overridden && CheckStringMapsEqual(a.Parameters, b.Parameters)
}

//ActionByOrder implements the sort.Interface for []ActionStructure
type ActionByOrder []ActionStructure

func (abp ActionByOrder) Len() int { return len(abp) }

func (abp ActionByOrder) Swap(i, j int) { abp[i], abp[j] = abp[j], abp[i] }

func (abp ActionByOrder) Less(i, j int) bool {
	var ipri int
	var jpri int
	//we've gotta go through and get the priorities
	for id, command := range abp[i].Device.Type.Commands {
		if id == abp[i].Action {
			ipri = command.Order
			break
		}
	}
	for id, command := range abp[j].Device.Type.Commands {
		if id == abp[j].Action {
			jpri = command.Order
			break
		}
	}
	return ipri < jpri
}

//CheckStringMapsEqual just takes two map[string]string and compares them.
func CheckStringMapsEqual(a map[string]string, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if b[k] != v {
			return false
		}
	}

	return true
}

//CheckStringSliceEqual is a simple helper to check if two string slices contain
//the same elements.
func CheckStringSliceEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ContainsAnyTags returns true if the taglist contains any of the specified tags
func ContainsAnyTags(tagList []string, tags ...string) bool {
	for i := range tags {
		for j := range tagList {
			if tagList[j] == tags[i] {
				return true
			}
		}
	}

	return false
}
