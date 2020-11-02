package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

const (
	_cRequestID = "requestID"
	_cRoom      = "room"
)

const (
	_hRequestID    = "X-Request-ID"
	_hForwardedFor = "X-Forwarded-For"
	_hContentType  = "Content-Type"
)

// RequestID requests that the client provides an id to be used in log statements.
func (h *Handlers) RequestID(c *gin.Context) {
	var id string
	if c.GetHeader(_hRequestID) != "" {
		id = c.GetHeader(_hRequestID)
	} else {
		uid, err := ksuid.NewRandom()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		id = uid.String()
	}

	c.Set(_cRequestID, id)
	c.Next()
}

// Log logs the start and end of client requests.
func (h *Handlers) Log(c *gin.Context) {
	id := c.GetString(_cRequestID)
	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	start := time.Now()
	log.Info("Starting request", zap.String("from", c.ClientIP()), zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path))
	c.Next()
	log.Info("Finished request", zap.Int("statusCode", c.Writer.Status()), zap.Duration("took", time.Since(start)))
}

// Room parses the http parameter "room", gets the appropriate room, and sets the parameter _cRoom.
func (h *Handlers) Room(c *gin.Context) {
	roomID := c.Param("room")
	if roomID == "" {
		c.String(http.StatusBadRequest, "must include room")
		c.Abort()
		return
	}

	id := c.GetString(_cRequestID)
	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Debug("Getting room", zap.String("room", roomID))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	room, err := h.DataService.RoomConfig(ctx, roomID)
	switch {
	case err != nil:
		c.String(http.StatusInternalServerError, "unable to get room %s", err)
		c.Abort()
		return
	}

	log.Debug("Got room")

	c.Set(_cRoom, room)
	c.Next()
}
