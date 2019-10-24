package drivers

type VideoSwitcherDSP interface {
	Device
	videoSwitcher
	dsp
}

func CreateVideoSwitcherDSPServer(vsdsp VideoSwitcherDSP) Server {
	e := newEchoServer()

	addDeviceRoutes(e, vsdsp)
	addVideoSwitcherRoutes(e, vsdsp)
	addDSPRoutes(e, vsdsp)

	return wrapEchoServer(e)
}
