package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/graph"
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

func (h *Handlers) SetRoomState(c echo.Context) error {
	roomID := c.Param("room")
	if len(roomID) == 0 {
		return c.String(http.StatusBadRequest, "room must be in the format BLDG-ROOM")
	}
	var stateReq api.StateRequest
	err := json.Unmarshal(c.Request().Body, &stateReq)
	if err != nil {
		return c.String(http.StatusBadRequest, "error unmarshaling state request")
	}

	// gotta get the current room state to compare
	err = h.GetRoomState(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	devices, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp, err := state.SetDevices(ctx, devices, env)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
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
