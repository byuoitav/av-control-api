package state

import "github.com/byuoitav/av-control-api/api"

type GetSetter struct {
	Environment string
	Logger      api.Logger
}

func boolP(b bool) *bool {
	return &b
}
