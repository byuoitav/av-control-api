package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	devices, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var stateReq api.StateRequest
	err = json.NewDecoder(c.Request().Body).Decode(&stateReq)
	if err != nil {
		er := fmt.Sprintf("error decoding state request: %s", err.Error())
		return c.String(http.StatusBadRequest, er)
	}

	if len(stateReq.OutputGroups) == 0 {
		return c.String(http.StatusBadRequest, "no devices found in request")
	}

	resp, err := state.SetDevices(ctx, stateReq, devices, h.Environment)
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
