package drivers

import (
	"context"
	"errors"
	"fmt"
	sync "sync"
)

type Drivers struct {
	drivers   map[string]*Driver
	driversMu sync.RWMutex
}

func New() *Drivers {
	return &Drivers{
		drivers: make(map[string]*Driver),
	}
}

type GetDeviceFunc func(context.Context, string) (Device, error)

type Driver struct {
	GetDevice GetDeviceFunc
}

// Register registers a driver with the given name. Name must not be empty.
func (d *Drivers) Register(name string, driver *Driver) error {
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
func (d *Drivers) MustRegister(name string, driver *Driver) {
	if err := d.Register(name, driver); err != nil {
		panic(err)
	}
}

// Get returns the driver that was registered with name.
// Returns nil if a matching driver has not been registered.
func (d *Drivers) Get(name string) *Driver {
	d.driversMu.RLock()
	defer d.driversMu.RUnlock()

	return d.drivers[name]
}
