package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func (h *Handlers) GetRoomConfiguration(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "must include a room")
	}

	id := c.Get(_cRequestID).(string)
	log := h.Logger.With(zap.String("requestID", id))

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	ctx = api.WithRequestID(ctx, id)

	log.Info("Getting room", zap.String("endpoint", c.Request().URL.String()), zap.String("room", roomID))
	room, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		log.Warn("failed to get room", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Successfully got configuration")
	return c.JSON(http.StatusOK, room)
}

func (h *Handlers) GetDriverMapping(c echo.Context) error {
	id := c.Get(_cRequestID).(string)
	log := h.Logger.With(zap.String("requestID", id))

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	ctx = api.WithRequestID(ctx, id)

	log.Info("Getting driver mapping", zap.String("endpoint", c.Request().URL.String()))
	mapping, err := h.DataService.DriverMapping(ctx)
	if err != nil {
		log.Warn("failed to get driver mapping", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Info("Successfully got driver mapping")
	return c.JSON(http.StatusOK, mapping)
}
