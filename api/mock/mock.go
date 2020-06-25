package mock

import (
	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
)

type Config interface {
	Room() api.Room
	Devices() map[string]drivers.Device
}
