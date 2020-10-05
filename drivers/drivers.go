package drivers

import (
	"context"
	"errors"
	"fmt"
	sync "sync"
)

type (
	Drivers interface {
		MustRegister(string, *Driver)
		Get(string) *Driver
		List() []string
	}

	GetDeviceFunc func(context.Context, string) (Device, error)
	Driver        struct {
		GetDevice GetDeviceFunc
	}

	drivers struct {
		drivers   map[string]*Driver
		driversMu sync.RWMutex
	}
)

func New() Drivers {
	return &drivers{
		drivers: make(map[string]*Driver),
	}
}

// Register registers a driver with the given name. Name must not be empty.
func (d *drivers) Register(name string, driver *Driver) error {
	if name == "" {
		return errors.New("driver must have a name")
	}

	d.driversMu.Lock()
	defer d.driversMu.Unlock()

	// make sure this isn't a duplicate
	if _, ok := d.drivers[name]; ok {
		return fmt.Errorf("driver %q already registered", name)
	}

	d.drivers[name] = driver
	return nil
}

// MustRegister is like Register but panics if there is an error registering the driver.
func (d *drivers) MustRegister(name string, driver *Driver) {
	if err := d.Register(name, driver); err != nil {
		panic(err)
	}
}

// Get returns the driver that was registered with name.
// Returns nil if a matching driver has not been registered.
func (d *drivers) Get(name string) *Driver {
	d.driversMu.RLock()
	defer d.driversMu.RUnlock()

	return d.drivers[name]
}

// List returns the list of names that have been registered.
func (d *drivers) List() []string {
	d.driversMu.RLock()
	defer d.driversMu.RUnlock()

	var list []string
	for k := range d.drivers {
		list = append(list, k)
	}

	return list
}
