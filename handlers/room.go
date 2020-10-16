package handlers

import (
	"context"
	"net/http"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) GetRoomConfiguration(c *gin.Context) {
	room := c.MustGet(_cRoom).(avcontrol.RoomConfig)
	c.JSON(http.StatusOK, room)
}

func (h *Handlers) GetRoomState(c *gin.Context) {
	room := c.MustGet(_cRoom).(avcontrol.RoomConfig)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		ctx = avcontrol.WithRequestID(ctx, id)
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Getting room state", zap.String("room", room.ID))

	resp, err := h.State.Get(ctx, room)
	if err != nil {
		log.Warn("failed to get room state", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(resp.Errors) > 0 {
		log.Info("Got room state", zap.Int("numErrors", len(resp.Errors)))
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Info("Got room state")
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetRoomHealth(c *gin.Context) {
	room := c.MustGet(_cRoom).(avcontrol.RoomConfig)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		ctx = avcontrol.WithRequestID(ctx, id)
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Getting room health", zap.String("room", room.ID))

	resp, err := h.State.GetHealth(ctx, room)
	if err != nil {
		log.Warn("failed to get room health", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Got room health")
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) SetRoomState(c *gin.Context) {
	var stateReq avcontrol.StateRequest
	if err := c.Bind(&stateReq); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	room := c.MustGet(_cRoom).(avcontrol.RoomConfig)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		ctx = avcontrol.WithRequestID(ctx, id)
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Setting room state", zap.String("room", room.ID))

	resp, err := h.State.Set(ctx, room, stateReq)
	if err != nil {
		log.Warn("failed to set room state", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if len(resp.Errors) > 0 {
		log.Info("Set room state", zap.Int("numErrors", len(resp.Errors)))
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Info("Set room state")
	c.JSON(http.StatusOK, resp)
}
