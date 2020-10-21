package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/core"
	"go.uber.org/zap"
)

// TODO wrap new functions
func registerDrivers(d avcontrol.DriverRegistry, log *zap.Logger) {
	d.MustRegister("sonyADCP", &core.SonyADCPDriver{
		Log: log.Named("drivers/sonyADCP"),
	})
	// d.MustRegister("Atlona", &core.AtlonaDriver{})
	d.MustRegister("JAP", &core.JAPDriver{})
	// d.MustRegister("KeyDigital", &core.KeyDigitalDriver{})
	d.MustRegister("kramer/protocol3000", &core.KramerProtocol3000Driver{})
	// d.MustRegister("London", &core.LondonDriver{})
	// d.MustRegister("NEC", &core.NECDriver{})
	d.MustRegister("QSC", &core.QSCDriver{})
	d.MustRegister("sony", &core.SonyDriver{})
}
