package drivers

import (
	"errors"
	"fmt"
	sync "sync"

	avcontrol "github.com/byuoitav/av-control-api"
)

type registry struct {
	drivers   map[string]avcontrol.Driver
	driversMu sync.RWMutex
}

func New() avcontrol.DriverRegistry {
	return &registry{
		drivers: make(map[string]avcontrol.Driver),
	}
}

// Register registers a driver with the given name. Name must not be empty.
func (r *registry) Register(name string, driver avcontrol.Driver) error {
	if name == "" {
		return errors.New("driver must have a name")
	}

	r.driversMu.Lock()
	defer r.driversMu.Unlock()

	// make sure this isn't a duplicate
	if _, ok := r.drivers[name]; ok {
		return fmt.Errorf("driver %q already registered", name)
	}

	// wrap this driver with the deviceCache
	r.drivers[name] = &deviceCache{
		Driver: driver,
		cache:  make(map[string]avcontrol.Device),
	}

	return nil
}

// MustRegister is like Register but panics if there is an error registering the driver.
func (r *registry) MustRegister(name string, driver avcontrol.Driver) {
	if err := r.Register(name, driver); err != nil {
		panic(err)
	}
}

// Get returns the driver that was registered with name.
// Returns nil if a matching driver has not been registered.
func (r *registry) Get(name string) avcontrol.Driver {
	r.driversMu.RLock()
	defer r.driversMu.RUnlock()

	return r.drivers[name]
}

// List returns the list of names that have been registered.
func (r *registry) List() []string {
	r.driversMu.RLock()
	defer r.driversMu.RUnlock()

	var list []string
	for k := range r.drivers {
		list = append(list, k)
	}

	return list
}
