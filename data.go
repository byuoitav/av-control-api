package avcontrol

import (
	"context"
	"net/url"
	"strings"
)

// DataService is used by the API to get information about rooms and drivers that it should use.
type DataService interface {
	Room(ctx context.Context, id string) (Room, error)
	DriverMapping(ctx context.Context) (DriverMapping, error)
}

// DriverMapping is a map of a driver name -> DriverConfig
type DriverMapping map[string]DriverConfig

// DriverConfig contains all the information necessary to connect to a driver
type DriverConfig struct {
	Address string
	SSL     bool
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

	// Driver should match a driver in in the DriverMapping. That driver will be used to
	// communicate with this device.
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
