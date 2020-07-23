package drivers

import "fmt"

var drivers = make(map[string]Driver)

type Driver struct {
	GetDevice GetDeviceFunc
}

func Register(name string, d Driver) {
	if name == "" {
		panic("driver must have a name")
	}

	if _, dup := drivers[name]; dup {
		panic(fmt.Sprintf("driver named %q already registered", name))
	}

	drivers[name] = d
}

func Get(name string) (Driver, bool) {
	d, ok := drivers[name]
	return d, ok
}
