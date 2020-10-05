package handlers

import (
	"net/http"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	Host           string
	DataService    avcontrol.DataService
	Logger         *zap.Logger
	State          avcontrol.StateGetSetter
	DriverRegistry avcontrol.DriverRegistry
}

func (h *Handlers) Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handlers) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"drivers": h.DriverRegistry.List(),
	})
}
