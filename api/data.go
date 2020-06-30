package api

import (
	"context"
	"net/url"
	"strings"
)

type DataService interface {
	Room(ctx context.Context, id string) (Room, error)
	DriverMapping(ctx context.Context) (DriverMapping, error)
}

type DriverMapping map[string]DriverConfig

type DriverConfig struct {
	Address string
	SSL     bool
}

type Room struct {
	ID string `json:"id"`
	// TODO proxy API requests to here
	ProxyBaseURL *url.URL            `json:"-"`
	Devices      map[DeviceID]Device `json:"devices"`
}

type Device struct {
	Address string `json:"address"`
	Driver  string `json:"driver"`
	Ports   Ports  `json:"ports,omitempty"`
}

type Ports []Port

// Port I think? this will only be used on DSP 's ??
type Port struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p Ports) Names() []string {
	var names []string

	for i := range p {
		names = append(names, p[i].Name)
	}

	return names
}

func (p Ports) OfType(typ string) Ports {
	var tp Ports

	for i := range p {
		if strings.Contains(p[i].Type, typ) {
			tp = append(tp, p[i])
		}
	}

	return tp
}
