package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api/graph"
	"github.com/labstack/echo"
)

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

func (h *Handlers) GetRoomGraph(c echo.Context) error {
	roomID := c.Param("room")
	portType := c.Param("type")
	switch {
	case len(roomID) == 0:
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	case len(portType) == 0:
		return c.String(http.StatusBadRequest, "invalid port type")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	room, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	g := graph.NewGraph(room, portType)
	svg, err := graph.GraphToSvg(g)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "image/svg+xml", svg)
}

func (h *Handlers) GetRoomGraphTranspose(c echo.Context) error {
	roomID := c.Param("room")
	portType := c.Param("type")
	switch {
	case len(roomID) == 0:
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	case len(portType) == 0:
		return c.String(http.StatusBadRequest, "invalid port type")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	room, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	g := graph.NewGraph(room, portType)
	t := graph.Transpose(g)

	svg, err := graph.GraphToSvg(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "image/svg+xml", svg)
}
