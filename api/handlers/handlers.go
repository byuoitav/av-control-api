package handlers

import "github.com/byuoitav/av-control-api/api"

type Handlers struct {
	DataService api.DeviceService
	Logger      api.Logger
	State       api.StateGetSetter
}
