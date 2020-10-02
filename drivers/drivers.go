package drivers

import (
	"errors"
	"fmt"
	sync "sync"

	"golang.org/x/mod/semver"
)

type Drivers struct {
	// drivers is a map of name -> version -> driver
	drivers   map[string]map[string]*Driver
	driversMu sync.RWMutex
}

func New() *Drivers {
	return &Drivers{
		drivers: make(map[string]map[string]*Driver),
	}
}

type Driver struct {
	GetDevice GetDeviceFunc
}

// Register registers a driver with the given name/version. Name must not be empty,
// and version must be a valid semver (see https://pkg.go.dev/golang.org/x/mod/semver).
func (d *Drivers) Register(name, version string, driver *Driver) error {
	if name == "" {
		return errors.New("driver must have a name")
	}

	if !semver.IsValid(version) {
		return errors.New("version be a valid semantic version")
	}

	version = semver.Canonical(version)

	d.driversMu.Lock()
	defer d.driversMu.Unlock()

	versions, ok := d.drivers[name]
	if !ok {
		versions = make(map[string]*Driver)
		d.drivers[name] = versions
	}

	// make sure this isn't a duplicate
	if _, ok := versions[version]; ok {
		return fmt.Errorf("driver %s/%s already registered", name, version)
	}

	versions[version] = driver
	return nil
}

// MustRegister is like Register but panics if there is an error registering the driver.
func (d *Drivers) MustRegister(name, version string, driver *Driver) {
	if err := d.Register(name, version, driver); err != nil {
		panic(err)
	}
}

// Get returns the driver that was registered with name/version. If version is
// empty, the latest version is returned. Nil is returned if a matching driver has not been registered.
func (d *Drivers) Get(name, version string) *Driver {
	d.driversMu.RLock()
	defer d.driversMu.RUnlock()

	versions, ok := d.drivers[name]
	if !ok {
		return nil
	}

	if version == "" {
		// get the latest version
		for v := range versions {
			version = semver.Max(v, version)
		}
	}

	driver, ok := versions[version]
	if !ok {
		return nil
	}

	return driver
}
