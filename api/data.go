package api

import (
	"context"
	"encoding/json"
	"regexp"
)

type DataService interface {
	Room(context.Context, string) (Room, error)
}

type Room struct {
	ID      string   `json:"id"`
	Devices []Device `json:"devices"`
}

type Device struct {
	ID      DeviceID                  `json:"id"`
	Type    DeviceType                `json:"type"`
	Address string                    `json:"address"`
	Proxy   map[*regexp.Regexp]string `json:"-"`
	Ports   Ports                     `json:"ports,omitempty"`
}

type DeviceType struct {
	ID       string             `json:"id"`
	Commands map[string]Command `json:"commands,omitempty"`
}

type Command struct {
	URLs  map[string]string `json:"urls"`
	Order *int              `json:"order,omitempty"`
}

type Port struct {
	Name      string    `json:"name"`
	Endpoints Endpoints `json:"endpoints"`
	Incoming  bool      `json:"incoming"`
	Type      string    `json:"type"`
}

type Endpoints []DeviceID

func (e Endpoints) Contains(id DeviceID) bool {
	for i := range e {
		if e[i] == id {
			return true
		}
	}

	return false
}

type Ports []Port

func (p Ports) Outgoing() Ports {
	var toReturn Ports
	for _, port := range p {
		if !port.Incoming {
			toReturn = append(toReturn, port)
		}
	}

	return toReturn
}

func (p Ports) Incoming() Ports {
	var toReturn Ports
	for _, port := range p {
		if port.Incoming {
			toReturn = append(toReturn, port)
		}
	}

	return toReturn
}

func (d Device) MarshalJSON() ([]byte, error) {
	type Alias Device

	changed := struct {
		*Alias
		Proxy map[string]string `json:"proxy,omitempty"`
	}{
		Alias: (*Alias)(&d),
		Proxy: make(map[string]string),
	}

	for k, v := range d.Proxy {
		changed.Proxy[k.String()] = v
	}

	return json.Marshal(changed)
}
