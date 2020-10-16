package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/adcp"
	"go.uber.org/zap"
)

func registerDrivers(d avcontrol.DriverRegistry, log *zap.Logger) {
	d.MustRegister("sonyADCP", &adcp.SonyADCPDriver{
		Log: log.Named("drivers/sonyADCP"),
	})
}
