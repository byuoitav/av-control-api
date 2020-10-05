package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/adcp"
)

// TODO wrap new functions
func registerDrivers(d avcontrol.Drivers) {
	d.MustRegister("sonyADCP", &avcontrol.Driver{
		GetDevice: adcp.NewDevice,
	})
}
