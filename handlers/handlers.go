package handlers

import (
	"net/http"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers acts as the http handler interface for the av-control-api.
type Handlers struct {
	Host           string
	DataService    avcontrol.DataService
	Logger         *zap.Logger
	State          avcontrol.StateGetSetter
	DriverRegistry avcontrol.DriverRegistry
}

// Stats returns the status of the http server.
func (h *Handlers) Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// Info returns a list of the registered drivers to the user in the body of an http response.
func (h *Handlers) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"drivers": h.DriverRegistry.List(),
	})
}
