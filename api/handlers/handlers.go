package handlers

import "github.com/byuoitav/av-control-api/api"

type Handlers struct {
	DataService api.DataService
	Logger      api.Logger
	State       api.StateGetSetter
}
