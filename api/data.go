package api

import "context"

type DeviceService interface {
	Device(context.Context, string) (Device, error)
	Room(context.Context, string) ([]Device, error)
}

type Device struct {
	ID      string            `json:"_id"`
	TypeID  string            `json:"typeID"`
	Address string            `json:"address"`
	Proxy   map[string]string `json:"proxy"`
	Ports   []Port            `json:"ports"`
}

type DeviceType struct {
	ID       string             `json:"_id"`
	Commands map[string]Command `json:"commands"`
}

type Command struct {
	URLs  map[string]string `json:"urls"`
	Order *int              `json:"order,omitempty"`
}

type Port struct {
	Name     string   `json:"name"`
	Endpoint DeviceID `json:"endpoint"`
	Incoming bool     `json:"incoming"`
	Outgoing bool     `json:"outgoing"`
	Type     string   `json:"type"`
}
