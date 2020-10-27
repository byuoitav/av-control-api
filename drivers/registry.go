package drivers

import (
	"errors"
	"fmt"
	"io/ioutil"
	sync "sync"

	avcontrol "github.com/byuoitav/av-control-api"
	"gopkg.in/yaml.v2"
)

type registry struct {
	configs map[string]map[string]interface{}

	drivers   map[string]avcontrol.Driver
	driversMu sync.RWMutex
}

func New(configPath string) (avcontrol.DriverRegistry, error) {
	r := &registry{
		configs: make(map[string]map[string]interface{}),
		drivers: make(map[string]avcontrol.Driver),
	}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	if err := yaml.Unmarshal(buf, &r.configs); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	return r, nil
}

func NewWithConfig(configs map[string]map[string]interface{}) (avcontrol.DriverRegistry, error) {
	if configs == nil {
		configs = make(map[string]map[string]interface{})
	}

	r := &registry{
		configs: configs,
		drivers: make(map[string]avcontrol.Driver),
	}

	return r, nil
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
		return fmt.Errorf("registry/%s: already registered", name)
	}

	if err := driver.ParseConfig(r.configs[name]); err != nil {
		return fmt.Errorf("registry/%s: unable to parse config: %w", name, err)
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
