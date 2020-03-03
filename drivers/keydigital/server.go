package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	keydigital "github.com/byuoitav/keydigital-driver"
	"github.com/spf13/pflag"
)

// imports

func main() {
	var (
		port     int
		username string
		password string
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
	createVS := func(ctx context.Context, addr string) (drivers.VideoSwitcher, error) {
		vs, err := keydigital.CreateVideoSwitcher(context.TODO(), addr)
		if err != nil {
			return nil, fmt.Errorf("failed to discover device: %w", err)
		}

		return vs, nil
	}

	// create server
	server := drivers.CreateVideoSwitcherFunc(createVS)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
}
