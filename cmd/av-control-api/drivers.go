package main

import (
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/adcp"
)

// TODO wrap new functions
func registerDrivers(d drivers.Drivers) {
	d.MustRegister("sonyADCP", &drivers.Driver{
		GetDevice: drivers.CacheDevices(adcp.NewDevice),
	})
}
