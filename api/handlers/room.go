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

	ctx, cancel := context.WithTimeout(c.Request().Context(), 20*time.Second)
	defer cancel()

	devices, err := h.DataService.Room(ctx, roomID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var stateReq api.StateRequest
	err = json.NewDecoder(c.Request().Body).Decode(&stateReq)
	if err != nil {
		return c.String(http.StatusBadRequest, "error decoding state request")
	}

	if len(stateReq.Devices) == 0 {
		return c.String(http.StatusBadRequest, "no devices found in request")
	}

	// jk we need everything so we can do stuff like change volume on DSPs
	// // we only need the devices in the room affected by the request
	// var devices []api.Device
	// // this should just be all of the keys which are device IDs
	// for k := range stateReq {
	// 	for d := range room {
	// 		if string(k) == room[d].ID {
	// 			devices = append(devices, room[d])
	// 		}
	// 	}
	// }

	// if len(devices) == 0 {
	// 	return c.String(http.StatusBadRequest, "given devices were not found in given room")
	// }

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
