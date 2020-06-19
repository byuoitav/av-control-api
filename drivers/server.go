package drivers

import (
	"context"
	"net"
	"sync"
)

type Server interface {
	Serve(lis net.Listener) error
	Stop(ctx context.Context) error
}

func NewServer(newDev NewDeviceFunc) (Server, error) {
	newDev = saveDevicesFunc(newDev)
	return newGrpcServer(newDev), nil
}

func saveDevicesFunc(newDev NewDeviceFunc) NewDeviceFunc {
	m := &sync.Map{}

	return func(ctx context.Context, addr string) (Device, error) {
		if dev, ok := m.Load(addr); ok {
			return dev, nil
		}

		dev, err := newDev(ctx, addr)
		if err != nil {
			return dev, err
		}

		m.Store(addr, dev)
		return dev, nil
	}
}
