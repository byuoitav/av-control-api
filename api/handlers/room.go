package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func (h *Handlers) GetRoomState(c echo.Context) error {
	return nil
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

	return c.JSON(http.StatusOK, devices)
}
