package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/london-driver"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	var port int

	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")

	pflag.Parse()

	create := func(ctx context.Context, addr string) (drivers.DSP, error) {
		logger := drivers.Log.Named(addr)

		return london.New(addr, london.WithLogger(logger)), nil
	}

	server, err := drivers.CreateDSPServer(create)
	if err != nil {
		fmt.Printf("failed to create server: %s\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start listener: %s\n", err)
		os.Exit(1)
	}

	drivers.Config.Level.SetLevel(zap.DebugLevel)

	drivers.Log.Infof("Starting server on: %s", lis.Addr().String())
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
