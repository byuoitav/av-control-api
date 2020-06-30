package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"net/http"
	"os"

	"github.com/byuoitav/av-control-api/api/couch"
	"github.com/byuoitav/av-control-api/api/handlers"
	"github.com/byuoitav/av-control-api/api/state"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel int8

		host string
		env  string

		authAddr    string
		authToken   string
		disableAuth bool

		dbAddr     string
		dbUsername string
		dbPassword string
		dbInsecure bool
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVarP(&env, "env", "e", "default", "The deployment environment for the API")
	pflag.StringVarP(&host, "host", "h", "", "host of this server. necessary to proxy requests")
	pflag.StringVar(&authAddr, "auth-addr", "", "address of the auth server")
	pflag.StringVar(&authToken, "auth-token", "", "authorization token to use when calling the auth server")
	pflag.BoolVar(&disableAuth, "disable-auth", false, "disables auth checks")
	pflag.StringVar(&dbAddr, "db-address", "", "database address")
	pflag.StringVar(&dbUsername, "db-username", "", "database username")
	pflag.StringVar(&dbPassword, "db-password", "", "database password")
	pflag.BoolVar(&dbInsecure, "db-insecure", false, "don't use SSL in database connection")
	pflag.Parse()

	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
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
			StacktraceKey:  "trace",
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
	defer logger.Sync()

	// validate flags
	if host == "" {
		logger.Fatal("--host is required. use --help for more details")
	}

	// build the data service
	dsOpts := []couch.Option{
		couch.WithEnvironment(env),
	}

	if len(dbUsername) > 0 {
		dsOpts = append(dsOpts, couch.WithBasicAuth(dbUsername, dbPassword))
	}

	if dbInsecure {
		dsOpts = append(dsOpts, couch.WithInsecure())
	}

	ds, err := couch.New(context.TODO(), dbAddr, dsOpts...)
	if err != nil {
		logger.Fatal("unable to connect to data service", zap.Error(err))
	}

	// build the getsetter
	gs, err := state.New(context.TODO(), ds, logger)
	if err != nil {
		logger.Fatal("unable to build state get/setter", zap.Error(err))
	}

	// build http stuff
	middleware := handlers.Middleware{}
	handlers := handlers.Handlers{
		Logger:      logger,
		DataService: ds,
		State:       gs,
	}

	e := echo.New()

	// TODO maybe check the database health check
	// TODO add log level endpoint
	// TODO add auth

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	api := e.Group("/v1", middleware.RequestID)
	api.GET("/driverMapping", handlers.GetDriverMapping)
	api.GET("/room/:room", handlers.GetRoomConfiguration)
	api.GET("/room/:room/state", handlers.GetRoomState)
	api.PUT("/room/:room/state", handlers.SetRoomState)

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
}
