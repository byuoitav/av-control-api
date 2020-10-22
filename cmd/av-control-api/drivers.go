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
	d.MustRegister("atlona6x2", &core.Atlona6x2Driver{})
	d.MustRegister("atlona5x1", &core.Atlona5x1Driver{
		Log: log.Named("drivers/atlona"),
	})
	d.MustRegister("atlona4x1", &core.Atlona4x1Driver{})
	d.MustRegister("atlona2x1", &core.Atlona2x1Driver{})
	d.MustRegister("JAP", &core.JAPDriver{})
	d.MustRegister("keyDigital", &core.KeyDigitalDriver{
		Log: log.Named("drivers/keyDigital"),
	})
	d.MustRegister("kramerDSP", &core.KramerDSPDriver{
		Log: log.Named("driver/kramerDSP"),
	})
	d.MustRegister("kramerVia", &core.KramerViaDriver{
		Log: log.Named("drivers/kramerVia"),
	})
	d.MustRegister("kramerVS", &core.KramerVSDriver{
		Log: log.Named("drivers/kramerVS"),
	})
	d.MustRegister("kramerVSDSP", &core.KramerVSDSPDriver{
		Log: log.Named("driver/kramerVSDSP"),
	})

	d.MustRegister("london", &core.LondonDriver{
		Log: log.Named("drivers/london"),
	})
	// d.MustRegister("NEC", &core.NECDriver{})
	d.MustRegister("QSC", &core.QSCDriver{})
	d.MustRegister("sony", &core.SonyDriver{})
}
