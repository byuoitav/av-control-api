package avcontrol

import "context"

// TODO gotta figure out where to put these functions!!!! this _should_ return drivers.Device
type GetDeviceFunc func(context.Context, string) (interface{}, error)

type Driver struct {
	GetDevice GetDeviceFunc
}

type Drivers interface {
	Get(string) *Driver
}
