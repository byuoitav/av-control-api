package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/adcp"
)

// TODO wrap new functions
func registerDrivers(d avcontrol.DriverRegistry) {
	d.MustRegister("sonyADCP", &adcp.SonyADCPDriver{})
}
