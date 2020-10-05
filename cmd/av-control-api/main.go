package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"net/http"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/av-control-api/handlers"
	"github.com/byuoitav/av-control-api/state"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type dataServiceConfig struct {
	Addr     string
	Username string
	Password string
	Insecure bool
}

func main() {
	var (
		port     int
		logLevel string

		host string

		dataServiceConfig dataServiceConfig
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "L", "", "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVarP(&host, "host", "h", "", "host of this server. necessary to proxy requests")
	pflag.StringVar(&dataServiceConfig.Addr, "db-address", "", "database address")
	pflag.StringVar(&dataServiceConfig.Username, "db-username", "", "database username")
	pflag.StringVar(&dataServiceConfig.Password, "db-password", "", "database password")
	pflag.BoolVar(&dataServiceConfig.Insecure, "db-insecure", false, "don't use SSL in database connection")
	pflag.Parse()

	fmt.Printf("dataServiceConfig: %+v\n", dataServiceConfig)

	// build a logger
	config, log := logger(logLevel)
	defer log.Sync() // nolint:errcheck

	// validate flags
	if host == "" {
		log.Fatal("--host is required. use --help for more details")
	}

	// register all of the drivers
	drivers := drivers.New()
	registerDrivers(drivers)

	// ctx for setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build the data service
	ds := dataService(ctx, dataServiceConfig)

	// build the getsetter
	gs := &state.GetSetter{
		Logger:  log,
		Drivers: drivers,
	}

	// build http stuff
	handlers := handlers.Handlers{
		Host:        host,
		DataService: ds,
		Logger:      log,
		State:       gs,
	}

	// TODO add auth
	r := gin.New()
	r.Use(gin.Recovery())

	debug := r.Group("/debug")
	debug.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})
	debug.GET("/statz", handlers.Stats)
	debug.GET("/infoz", handlers.Info)
	debug.GET("/logz", func(c *gin.Context) {
		c.String(http.StatusOK, config.Level.String())
	})
	debug.GET("/logz/:level", func(c *gin.Context) {
		var level zapcore.Level
		if err := level.Set(c.Param("level")); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Printf("***\n\tSetting log level to %s\n***\n", level.String())
		config.Level.SetLevel(level)
		c.String(http.StatusOK, config.Level.String())
	})

	api := r.Group("/api/v1", handlers.RequestID, handlers.Log)

	room := api.Group("/room", handlers.Room, handlers.Proxy)
	room.GET("/:room", handlers.GetRoomConfiguration)
	room.GET("/:room/state", handlers.GetRoomState)
	room.PUT("/:room/state", handlers.SetRoomState)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("unable to bind listener", zap.Error(err))
	}

	log.Info("Starting server", zap.String("on", lis.Addr().String()))
	err = r.RunListener(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		log.Fatal("failed to serve", zap.Error(err))
	}
}
