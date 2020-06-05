package drivers

import (
	"context"
	"sync"
)

type DisplayDSP interface {
	Device
	Display
	DSP
}

type CreateDisplayDSPFunc func(context.Context, string) (DisplayDSP, error)

func CreateDisplayDSPServer(create CreateDisplayDSPFunc) (Server, error) {
	e := newEchoServer()
	m := &sync.Map{}

	ddsp := func(ctx context.Context, addr string) (DisplayDSP, error) {
		if ddsp, ok := m.Load(addr); ok {
			return ddsp.(DisplayDSP), nil
		}

		ddsp, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}

		m.Store(addr, ddsp)
		return ddsp, nil
	}

	dev := func(ctx context.Context, addr string) (Device, error) {
		return ddsp(ctx, addr)
	}

	dis := func(ctx context.Context, addr string) (Display, error) {
		return ddsp(ctx, addr)
	}

	dsp := func(ctx context.Context, addr string) (DSP, error) {
		return ddsp(ctx, addr)
	}

	addDeviceRoutes(e, dev)
	addDisplayRoutes(e, dis)
	addDSPRoutes(e, dsp)

	return wrapEchoServer(e), nil
}
