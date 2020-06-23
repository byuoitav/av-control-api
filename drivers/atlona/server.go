package main

import (
	"context"
	"fmt"
	"net"
	"os"

	atlona "github.com/byuoitav/atlona-driver"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/spf13/pflag"
)

// imports

func main() {
	var (
		port     int
		username string
		password string
	)

	// Variable declarations

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&username, "username", "u", "", "username for device")
	pflag.StringVarP(&password, "password", "p", "", "password for device")

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
	createVS := func(ctx context.Context, addr string) (drivers.Device, error) {
		//this dont work here for some reason
		vs, err := atlona.CreateVideoSwitcher(ctx, addr, username, password, drivers.Log.Named(addr))
		if err != nil {
			return nil, fmt.Errorf("failed to discover device: %w", err)
		}

		return vs, nil
	}

	// create server
	server, err := drivers.NewServer(createVS)
	if err != nil {
		fmt.Printf("failed to create server: %s\n", err)
		os.Exit(1)
	}
	if err = server.Serve(lis); err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
}
