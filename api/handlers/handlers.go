package handlers

import (
	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type Handlers struct {
	Environment string
	DataService api.DataService
	Logger      *zap.Logger
	State       api.StateGetSetter
}
