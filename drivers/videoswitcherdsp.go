package drivers

import "sync"

type VideoSwitcherDSP interface {
	Device
	videoSwitcher
	dsp
}

type CreateVideoSwitcherDSPFunc func(string) VideoSwitcherDSP

func CreateVideoSwitcherDSPServer(create CreateVideoSwitcherDSPFunc) Server {
	e := newEchoServer()
	m := &sync.Map{}

	vsdsp := func(addr string) VideoSwitcherDSP {
		if vsdsp, ok := m.Load(addr); ok {
			return vsdsp.(VideoSwitcherDSP)
		}

		vsdsp := create(addr)
		m.Store(addr, vsdsp)
		return vsdsp
	}

	dev := func(addr string) Device {
		return vsdsp(addr)
	}

	vs := func(addr string) VideoSwitcher {
		return vsdsp(addr)
	}

	dsp := func(addr string) DSP {
		return vsdsp(addr)
	}

	addDeviceRoutes(e, dev)
	addVideoSwitcherRoutes(e, vs)
	addDSPRoutes(e, dsp)

	return wrapEchoServer(e)
}
