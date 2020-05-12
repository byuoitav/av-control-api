package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api/state"
	"github.com/labstack/echo"
)

func (h *Handlers) GetRoomState(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	devices, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp, err := state.GetDevices(ctx, devices, h.Environment)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetRoomConfiguration(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	devices, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	_, err = json.Marshal(devices)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, devices)
}
