package drivers

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

var (
	ErrFuncNotSupported = errors.New("device does not support this function")
	ErrMissingAddress   = errors.New("must include address of the device")
	ErrMissingInput     = errors.New("missing input")
)

// CreateDeviceFunc is passed to CreateDeviceServer and is called to create a new Device struct whenever the Server needs to communicate  with a new Device.
type CreateDeviceFunc func(context.Context, string) (Device, error)

func CreateDeviceServer(create CreateDeviceFunc) (Server, error) {
	m := &sync.Map{}
	single := &singleflight.Group{}

	newDev := func(ctx context.Context, addr string) (Device, error) {
		if dev, ok := m.Load(addr); ok {
			return dev, nil
		}

		dev, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}

		m.Store(addr, dev)
		return dev, nil
	}

	dev, err := newDev(context.TODO(), "")
	if err != nil {
		return nil, err
	}

	// build a DeviceServer

	server := grpc.NewServer()
	RegisterDeviceServer(server, nil)

	return server, nil
}
