package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/epson-driver"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel int8
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.Parse()

	drivers.Config.Level.SetLevel(zapcore.Level(logLevel))

	create := func(ctx context.Context, addr string) (drivers.DisplayDSP, error) {
		logger := drivers.Log.Named(addr)

		return epson.NewProjector(addr, epson.WithLogger(logger)), nil
	}

	server, err := drivers.CreateDisplayDSPServer(create)
	if err != nil {
		fmt.Printf("failed to create server: %s\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	drivers.Log.Infof("Starting server on: %s", lis.Addr().String())
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
