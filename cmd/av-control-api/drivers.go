package main

import (
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/adcp"
)

// TODO wrap new functions
func registerDrivers(d *drivers.Drivers) {
	d.Register("sonyADCP", &drivers.Driver{
		GetDevice: adcp.NewDevice,
	})
}
