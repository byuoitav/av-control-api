package main

import (
	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers/core"
	"go.uber.org/zap"
)

// TODO wrap new functions
func registerDrivers(d avcontrol.DriverRegistry, log *zap.Logger) {
	d.MustRegister("sony/adcp", &core.SonyADCPDriver{
		Log: log.Named("drivers/sony/adcp"),
	})

	d.MustRegister("atlona6x2", &core.Atlona6x2Driver{})

	d.MustRegister("atlona5x1", &core.Atlona5x1Driver{
		Log: log.Named("drivers/atlona"),
	})

	d.MustRegister("atlona4x1", &core.Atlona4x1Driver{})

	d.MustRegister("atlona2x1", &core.Atlona2x1Driver{})

	d.MustRegister("JAP", &core.JAPDriver{
		Log: log.Named("drivers/JAP"),
	})

	d.MustRegister("keyDigital", &core.KeyDigitalDriver{
		Log: log.Named("drivers/keyDigital"),
	})

	d.MustRegister("kramerDSP", &core.KramerDSPDriver{
		Log: log.Named("driver/kramerDSP"),
	})

	d.MustRegister("kramer/via", &core.KramerViaDriver{
		Log: log.Named("drivers/kramerVia"),
	})

	d.MustRegister("kramerVSDSP", &core.KramerVSDSPDriver{
		Log: log.Named("driver/kramerVSDSP"),
	})

	d.MustRegister("kramer/protocol3000", &core.KramerProtocol3000Driver{
		Log: log.Named("drivers/protocol3000"),
	})

	d.MustRegister("london", &core.LondonDriver{
		Log: log.Named("drivers/london"),
	})
	// d.MustRegister("NEC", &core.NECDriver{})

	d.MustRegister("QSC", &core.QSCDriver{
		Log: log.Named("drivers/qsc"),
	})

	d.MustRegister("sony/bravia", &core.SonyDriver{
		Log: log.Named("drivers/sony/bravia"),
	})
}
