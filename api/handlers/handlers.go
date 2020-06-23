package handlers

import "github.com/byuoitav/av-control-api/api"

type Handlers struct {
	Environment string
	DataService api.DataService
	Logger      api.Logger
	State       api.StateGetSetter
}
