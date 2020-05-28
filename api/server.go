package main

import (
	"net/http"
	"os"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/av-control-api/api/handlers"
	"github.com/byuoitav/av-control-api/api/health"
	avapi "github.com/byuoitav/av-control-api/api/init"
	hub "github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/status/databasestatus"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/common/v2/events"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
)

func main() {
	// Parse flags

	var env string

	pflag.StringVarP(&env, "env", "e", "fallback", "The deployment environment for the API")
	pflag.Parse()

	var nerr *nerr.E

	base.Messenger, nerr = messenger.BuildMessenger(os.Getenv("HUB_ADDRESS"), hub.Messenger, 1000)
	if nerr != nil {
		log.L.Errorf("unable to connect to the hub: %s", nerr.String())
	}

	go func() {
		err := avapi.CheckRoomInitialization()
		if err != nil {
			base.PublishError("Fail to run init script. Terminating. ERROR:"+err.Error(), events.Error, os.Getenv("SYSTEM_ID"))
			log.L.Errorf("Could not initialize room. Error: %v\n", err.Error())
		}
	}()

	h := handlers.RoomHandler{
		Environment: env,
	}

	port := ":8000"
	router := common.NewRouter()

	router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "Up an running!")
	})
	router.GET("/mstatus", databasestatus.Handler)
	router.GET("/status", databasestatus.Handler)

	// PUT requests
	router.PUT("/buildings/:building/rooms/:room", h.SetRoomState, auth.AuthorizeRequest("write-state", "room", h.GetRoomResource))

	// room status
	router.GET("/buildings/:building/rooms/:room", h.GetRoomState, auth.AuthorizeRequest("read-state", "room", h.GetRoomResource))
	router.GET("/buildings/:building/rooms/:room/configuration", h.GetRoomByNameAndBuilding, auth.AuthorizeRequest("read-config", "room", h.GetRoomResource))

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	go health.StartupCheckAndReport()

	router.StartServer(&server)
}
