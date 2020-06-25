package mock

import (
	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/drivers/mock"
)

// SimpleRoom implements api.DeviceService and represents a simple room with the following devices:
//
// * a tv
// * a HDMI input
// * a VIA input
//
// The audio for this room is coming directly out of the TV,
// and the inputs are going directly into hdmi 1 & 2 on the TV.
// the via is also capable of getting/setting it's own volume/mute.
type SimpleRoom struct {
}

func (s *SimpleRoom) Room() api.Room {
	return api.Room{
		ID:           "ITB-1101",
		ProxyBaseURL: "",
		Devices: map[api.DeviceID]api.Device{
			"ITB-1101-D1": api.Device{
				Address: "ITB-1101-D1.av",
			},
		},
	}
}

func (s *SimpleRoom) Devices() map[string]drivers.Device {
	return map[string]drivers.Device{
		"ITB-1101-D1.av": &mock.TV{},
	}
}
