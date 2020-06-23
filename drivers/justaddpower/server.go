package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	justaddpower "github.com/byuoitav/justaddpower-driver"
	"github.com/spf13/pflag"
)

// imports

func main() {
	var (
		port int
	)

	// variable declarations

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	// other flags

	pflag.Parse()

	// create a net.Listener to run the server on
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	// import driver library
	createJ := func(ctx context.Context, addr string) (drivers.Device, error) {
		return &justaddpower.JustAddPowerReciever{
			Address: addr,
		}, nil
	}

	// create server
	server, err := drivers.NewServer(createJ)
	if err != nil {
		fmt.Printf("failed to create server: %s\n", err)
		os.Exit(1)
	}

	if err = server.Serve(lis); err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
}
