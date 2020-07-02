package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"net/http"
	"os"

	"github.com/byuoitav/av-control-api/api/couch"
	"github.com/byuoitav/av-control-api/api/handlers"
	"github.com/byuoitav/av-control-api/api/state"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel string

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
	pflag.StringVarP(&logLevel, "log-level", "L", "", "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
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

	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		fmt.Printf("invalid log level: %s\n", err.Error())
		os.Exit(1)
	}

	config := zap.Config{
		Level: zap.NewAtomicLevelAt(level),
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

	log, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}
	defer log.Sync()

	// validate flags
	if host == "" {
		log.Fatal("--host is required. use --help for more details")
	}

	// context for setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build the data service
	if dbInsecure {
		dbAddr = "http://" + dbAddr
	} else {
		dbAddr = "https://" + dbAddr
	}

	dsOpts := []couch.Option{
		couch.WithEnvironment(env),
	}

	if len(dbUsername) > 0 {
		dsOpts = append(dsOpts, couch.WithBasicAuth(dbUsername, dbPassword))
	}

	ds, err := couch.New(ctx, dbAddr, dsOpts...)
	if err != nil {
		log.Fatal("unable to connect to data service", zap.Error(err))
	}

	// build the getsetter
	gs, err := state.New(ctx, ds, log)
	if err != nil {
		log.Fatal("unable to build state get/setter", zap.Error(err))
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
	debug.GET("/statsz", handlers.Stats)
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

	api := r.Group("/v1", handlers.RequestID, handlers.Log)

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
