package api

import (
	"context"
	"encoding/json"
	"regexp"
)

type DeviceService interface {
	Device(context.Context, string) (Device, error)
	Room(context.Context, string) ([]Device, error)
}

type Device struct {
	ID      DeviceID                  `json:"id"`
	Type    DeviceType                `json:"type"`
	Address string                    `json:"address"`
	Proxy   map[*regexp.Regexp]string `json:"-"`
	Ports   []Port                    `json:"ports,omitempty"`
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
	Outgoing  bool      `json:"outgoing"`
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
