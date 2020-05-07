package api

import "context"

type DeviceService interface {
	Device(context.Context, string) (Device, error)
	Room(context.Context, string) ([]Device, error)
}

type Device struct {
	ID      string
	TypeID  string
	Address string
	Proxy   map[string]string
	Ports   []Port
}

type DeviceType struct {
	ID       string
	Commands map[string]Command
}

type Command struct {
	URLs  map[string]string
	Order *int
}

type Port struct {
	Name     string
	Endpoint DeviceID
	Incoming bool
	Outgoing bool
	Type     string
}
