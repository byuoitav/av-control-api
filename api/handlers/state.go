package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func (h *Handlers) GetRoomState(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	}

	id := c.Get(_cRequestID).(string)
	log := h.Logger.With(zap.String("requestID", id), zap.String("endpoint", c.Request().URL.String()))

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	ctx = api.WithRequestID(ctx, id)
	log.Info("Getting room", zap.String("room", roomID))

	room, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		log.Warn("failed to get devices", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Got room. Getting state", zap.Int("numDevices", len(room.Devices)))

	resp, err := h.State.Get(ctx, room)
	if err != nil {
		log.Warn("failed to get state", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Successfully got state")
	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) SetRoomState(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	var stateReq api.StateRequest
	err := c.Bind(&stateReq)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(stateReq.Devices) == 0 {
		return c.String(http.StatusBadRequest, "no devices found in request")
	}

	id := c.Get(_cRequestID).(string)
	log := h.Logger.With(zap.String("requestID", id), zap.String("endpoint", c.Request().URL.String()))

	ctx = api.WithRequestID(ctx, id)
	log.Info("Getting room", zap.String("room", roomID))

	room, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		log.Warn("failed to get room", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Got room. Setting state", zap.Int("numDevices", len(room.Devices)))

	resp, err := h.State.Set(ctx, room, stateReq)
	if err != nil {
		log.Warn("failed to set state", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Successfully set state")
	return c.JSON(http.StatusOK, resp)
}
