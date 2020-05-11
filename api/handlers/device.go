package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func (h *Handlers) GetDeviceState(c echo.Context) error {
	return nil
}

func (h *Handlers) GetDeviceConfiguration(c echo.Context) error {
	deviceID := c.Param("device")
	if len(deviceID) == 0 {
		return c.String(http.StatusBadRequest, "device must be in the format BLDG-ROOM-CP1")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	device, err := h.DataService.Device(ctx, deviceID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, device)
}
