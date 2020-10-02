package avcontrol

import (
	"context"
	"net/url"
	"strings"
)

// DataService is used by to get information about rooms.
type DataService interface {
	Room(ctx context.Context, id string) (Room, error)
}

// Room is the configuration a room as perceived by the av-control-api.
type Room struct {
	// ID is the ID of this room, in the format of Building-Room
	ID string `json:"id"`

	// Proxy is used to decide if this instance of the API should handle this request,
	// or if it should proxy the request to another instance of the API. See (handlers.Handlers{}).Proxy()
	// for more details.
	Proxy *url.URL `json:"-"`

	// Devices is map of devices that exist in this room.
	Devices map[DeviceID]Device `json:"devices"`
}

// Device contains information about a given device.
type Device struct {
	// Address is the Hostname or IP address of the device
	Address string `json:"address"`

	// Driver should match with a driver that has been registered with the API. The matching driver will be used to
	// communicate with this device. If no drivers with this name have been registered, requests will fail.
	Driver string `json:"driver"`

	// Ports are logical ports that the API must know about to be able to control. For a DSP,
	// these are control block names.
	Ports Ports `json:"ports,omitempty"`
}

// Ports is a slice of Port
type Ports []Port

// Port is the configuration for a specific port
type Port struct {
	// Name is the name of the port that a driver will understand and use. For example,
	// in a QSC DSP, this is the NamedControl
	Name string `json:"name"`

	// Type is the type of port this port is. There are two used types right now:
	//  - volume
	//  - mute
	// Each port is used when getting or setting that field on a device.
	Type string `json:"type"`
}

// Names returns the list of names of these ports.
func (p Ports) Names() []string {
	var names []string

	for i := range p {
		names = append(names, p[i].Name)
	}

	return names
}

// OfType returns the Ports of a specific type.
func (p Ports) OfType(typ string) Ports {
	var tp Ports

	for i := range p {
		if strings.Contains(p[i].Type, typ) {
			tp = append(tp, p[i])
		}
	}

	return tp
}
