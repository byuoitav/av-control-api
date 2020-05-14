package handlers

import (
	"context"
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

func test(room []api.Device) {
	g := graph.NewGraph(room, "audio")

	path := graph.PathToEnd(g, "ITB-1108B-MIC1")
	if len(path) == 0 {
		// no path down
	}

	end := path[len(path)-1]
	fmt.Printf("%s %s %s\n%v\n", end.Dst.Device.ID, end.DstPort.Name, end.Dst.Address, end.Dst.Type.Commands)
}
