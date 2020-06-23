package api

import (
	"context"
	"encoding/json"
	"regexp"
)

type DataService interface {
	Room(ctx context.Context, id string) (Room, error)
	DriverMapping(ctx context.Context) (DriverMapping, error)
}

type DriverMapping map[string]string

type Room struct {
	ID      string              `json:"id"`
	Devices map[DeviceID]Device `json:"devices"`
}

type Device struct {
	Address string                    `json:"address"`
	Driver  string                    `json:"driver"`
	Proxy   map[*regexp.Regexp]string `json:"-"`
	Ports   Ports                     `json:"ports,omitempty"`
}

// Port I think? this will only be used on DSP 's ??
type Port struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Ports []Port

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
