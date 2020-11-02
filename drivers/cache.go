package drivers

import (
	"context"
	sync "sync"

	avcontrol "github.com/byuoitav/av-control-api"
	"golang.org/x/sync/singleflight"
)

var _ avcontrol.Driver = &deviceCache{}

type deviceCache struct {
	avcontrol.Driver

	single  singleflight.Group
	cache   map[string]avcontrol.Device
	cacheMu sync.Mutex
}

// CreateDevice gets and returns the device from the cache.
// If the device is not found in the cache, then a new one is created, added to the cache, and returned.
func (c *deviceCache) CreateDevice(ctx context.Context, addr string) (avcontrol.Device, error) {
	if dev, ok := c.get(addr); ok {
		return dev, nil
	}

	val, err, _ := c.single.Do(addr, func() (interface{}, error) {
		if dev, ok := c.get(addr); ok {
			return dev, nil
		}

		dev, err := c.Driver.CreateDevice(ctx, addr)
		if err != nil {
			return nil, err
		}

		c.cacheMu.Lock()
		defer c.cacheMu.Unlock()
		c.cache[addr] = dev

		return dev, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(avcontrol.Device), nil
}

func (c *deviceCache) get(addr string) (avcontrol.Device, bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	dev, ok := c.cache[addr]
	return dev, ok
}
