package drivers

import (
	"context"
	sync "sync"

	"golang.org/x/sync/singleflight"
)

func CacheDevices(get GetDeviceFunc) GetDeviceFunc {
	devs := make(map[string]Device)
	devsMu := sync.RWMutex{}
	single := singleflight.Group{}

	check := func(addr string) (Device, bool) {
		devsMu.RLock()
		defer devsMu.RUnlock()

		dev, ok := devs[addr]
		return dev, ok
	}

	return func(ctx context.Context, addr string) (Device, error) {
		if dev, ok := check(addr); ok {
			return dev, nil
		}

		val, err, _ := single.Do(addr, func() (interface{}, error) {
			if dev, ok := check(addr); ok {
				return dev, nil
			}

			dev, err := get(ctx, addr)
			if err != nil {
				return nil, err
			}

			devsMu.Lock()
			defer devsMu.Unlock()
			devs[addr] = dev

			return dev, nil
		})
		if err != nil {
			return nil, err
		}

		return val.(Device), nil
	}
}
