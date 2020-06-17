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

// Port I think? this will only be used on DSP's ??
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
