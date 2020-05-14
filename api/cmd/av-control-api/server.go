package main

import (
	"errors"
	"fmt"
	"net"

	"net/http"
	"os"

	"github.com/byuoitav/av-control-api/api/couch"
	"github.com/byuoitav/av-control-api/api/handlers"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel int8

		env         string
		authAddr    string
		authToken   string
		disableAuth bool
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVar(&authAddr, "auth-addr", "", "address of the auth server")
	pflag.StringVar(&authToken, "auth-token", "", "authorization token to use when calling the auth server")
	pflag.BoolVar(&disableAuth, "disable-auth", false, "disables auth checks")

	pflag.StringVarP(&env, "env", "e", "default", "The deployment environment for the API")
	pflag.Parse()

	// build the logger
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}

	// TODO get a database interface
	db := couch.DataService{
		DBAddress:  os.Getenv("DB_ADDRESS"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}

	handlers := handlers.Handlers{
		Environment: env,
		DataService: &db,
	}

	e := echo.New()

	// TODO maybe check the database health check
	// TODO add log level endpoint
	// TODO add auth

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	e.GET("/room/:room", handlers.GetRoomConfiguration)
	e.GET("/room/:room/state", handlers.GetRoomState)
	e.GET("/room/:room/graph/:type", handlers.GetRoomGraph)

	e.GET("/device/:device", handlers.GetDeviceConfiguration)
	e.GET("/device/:device/state", handlers.GetDeviceState)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("unable to bind listener", zap.Error(err))
	}

	logger.Info("Starting server", zap.String("on", lis.Addr().String()))
	err = e.Server.Serve(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		logger.Fatal("failed to serve", zap.Error(err))
	}

	//// PUT requests
	//router.PUT("/buildings/:building/rooms/:room", h.SetRoomState, auth.AuthorizeRequest("write-state", "room", h.GetRoomResource))

	//// room status
	//router.GET("/buildings/:building/rooms/:room", h.GetRoomState, auth.AuthorizeRequest("read-state", "room", h.GetRoomResource))
	//router.GET("/buildings/:building/rooms/:room/configuration", h.GetRoomByNameAndBuilding, auth.AuthorizeRequest("read-config", "room", h.GetRoomResource))
	//router.PUT("/log-level/:level", log.SetLogLevel)
	//router.GET("/log-level", log.GetLogLevel)
}
