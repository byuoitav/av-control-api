package avcontrol

import (
	"context"
)

type DriverRegistry interface {
	// Register registers a driver with the given name.
	Register(string, Driver) error

	// MustRegister is like Register but panics if there is an error registering the driver.
	MustRegister(string, Driver)

	// Get returns the driver that was registered with name.
	Get(string) Driver

	// List returns the list of names that have been registered.
	List() []string
}

type Driver interface {
	// CreateDevice is called whenever the API needs to get or set state on
	// a device. Addr is the address (IP or hostname) of the device.
	// See Device for interfaces the returned struct should implement.
	CreateDevice(ctx context.Context, addr string) (Device, error)

	// ParseConfig is called by the DriverRegistry when a driver is registered.
	ParseConfig(map[string]interface{}) error
}
