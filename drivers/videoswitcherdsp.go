package drivers

import (
	"context"
	"sync"
)

type VideoSwitcherDSP interface {
	Device
	videoSwitcher
	dsp
}

type CreateVideoSwitcherDSPFunc func(context.Context, string) (VideoSwitcherDSP, error)

func CreateVideoSwitcherDSPServer(create CreateVideoSwitcherDSPFunc) Server {
	e := newEchoServer()
	m := &sync.Map{}

	vsdsp := func(ctx context.Context, addr string) (VideoSwitcherDSP, error) {
		if vsdsp, ok := m.Load(addr); ok {
			return vsdsp.(VideoSwitcherDSP), nil
		}

		vsdsp, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}
		m.Store(addr, vsdsp)
		return vsdsp, nil
	}

	dev := func(ctx context.Context, addr string) (Device, error) {
		return vsdsp(ctx, addr)
	}

	vs := func(ctx context.Context, addr string) (VideoSwitcher, error) {
		return vsdsp(ctx, addr)
	}

	dsp := func(ctx context.Context, addr string) (DSP, error) {
		return vsdsp(ctx, addr)
	}

	addDeviceRoutes(e, dev)
	addVideoSwitcherRoutes(e, vs)
	addDSPRoutes(e, dsp)

	return wrapEchoServer(e)
}
