package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/core"
)

// TODO wrap new functions
func registerDrivers(d avcontrol.DriverRegistry) {
	d.MustRegister("sonyADCP", &core.SonyADCPDriver{})
	// d.MustRegister("Atlona", &core.AtlonaDriver{})
	d.MustRegister("JAP", &core.JAPDriver{})
	// d.MustRegister("KeyDigital", &core.KeyDigitalDriver{})
	d.MustRegister("kramer", &core.KramerDriver{})
	// d.MustRegister("London", &core.LondonDriver{})
	// d.MustRegister("NEC", &core.NECDriver{})
	d.MustRegister("QSC", &core.QSCDriver{})
	d.MustRegister("sony", &core.SonyDriver{})
}
