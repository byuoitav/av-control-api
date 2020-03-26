package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/byuoitav/av-control-api/api/db"
	"github.com/byuoitav/av-control-api/api/helpers"
	"github.com/byuoitav/av-control-api/api/inputgraph"
	"github.com/byuoitav/av-control-api/api/rest"
	"github.com/byuoitav/av-control-api/api/state"
	"github.com/byuoitav/common/log"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const (
	timeout = 50 * time.Millisecond
)

type RoomHandler struct {
	Environment string
}

// GetRoomResource returns the resourceID for a request
func (r *RoomHandler) GetRoomResource(context echo.Context) string {
	return context.Param("building") + "-" + context.Param("room")
}

//GetRoomState to get the current state of a room
func (r *RoomHandler) GetRoomState(context echo.Context) error {
	building, room := context.Param("building"), context.Param("room")

	status, err := state.GetRoomState(building, room, r.Environment)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, status)
}

//GetRoomByNameAndBuilding is almost identical to GetRoomByName
func (r *RoomHandler) GetRoomByNameAndBuilding(context echo.Context) error {
	building, roomName := context.Param("building"), context.Param("room")

	log.L.Info("Getting room...")
	room, err := db.GetDB().GetRoom(fmt.Sprintf("%s-%s", building, roomName))
	if err != nil {
		return context.JSON(http.StatusBadRequest, helpers.ReturnError(err))
	}

	//we need to add the input reachability stuff
	reachable, err := inputgraph.GetVideoDeviceReachability(room)

	log.L.Info("Done.\n")
	return context.JSON(http.StatusOK, reachable)
}

// SetRoomState to update the state of the room
func (r *RoomHandler) SetRoomState(ctx echo.Context) error {
	building, room := ctx.Param("building"), ctx.Param("room")

	log.L.Infof("%s", color.HiGreenString("[handlers] putting room changes..."))

	var roomInQuestion rest.PublicRoom
	err := ctx.Bind(&roomInQuestion)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, helpers.ReturnError(err))
	}

	roomInQuestion.Room = room
	roomInQuestion.Building = building
	var report rest.PublicRoom

	gctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	rslvr := net.Resolver{}
	hn, err := rslvr.LookupAddr(gctx, ctx.RealIP())

	color.Set(color.FgYellow, color.Bold)
	if err != nil || len(hn) == 0 {
		log.L.Debugf("REQUESTOR: %s", ctx.RealIP())
		color.Unset()
		report, err = state.SetRoomState(roomInQuestion, r.Environment, ctx.RealIP())
	} else if strings.Contains(hn[0], "localhost") {
		log.L.Debugf("REQUESTOR: %s", os.Getenv("SYSTEM_ID"))
		color.Unset()
		report, err = state.SetRoomState(roomInQuestion, r.Environment, os.Getenv("SYSTEM_ID"))
	} else {
		log.L.Debugf("REQUESTOR: %s", hn[0])
		color.Unset()
		report, err = state.SetRoomState(roomInQuestion, r.Environment, hn[0])
	}

	if err != nil {
		log.L.Errorf("Error: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, helpers.ReturnError(err))
	}

	log.L.Info("Done.\n")

	return ctx.JSON(http.StatusOK, report)
}
