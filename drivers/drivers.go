package drivers

import (
	"errors"
	"fmt"
	sync "sync"

	avcontrol "github.com/byuoitav/av-control-api"
)

type drivers struct {
	drivers   map[string]avcontrol.Driver
	driversMu sync.RWMutex
}

func New() avcontrol.DriverRegistry {
	return &drivers{
		drivers: make(map[string]avcontrol.Driver),
	}
}

// Register registers a driver with the given name. Name must not be empty.
func (d *drivers) Register(name string, driver avcontrol.Driver) error {
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
func (d *drivers) MustRegister(name string, driver avcontrol.Driver) {
	if err := d.Register(name, driver); err != nil {
		panic(err)
	}
}

// Get returns the driver that was registered with name.
// Returns nil if a matching driver has not been registered.
func (d *drivers) Get(name string) avcontrol.Driver {
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
