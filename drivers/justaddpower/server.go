package main

import (
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

	pflag.IntVarP(&port, "port", "P", 80, "port to run the server on")
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
	createJ := func(addr string) drivers.VideoSwitcher {
		return &justaddpower.JustAddPowerReciever{
			Address: addr,
		}
	}

	// create server
	server := drivers.CreateVideoSwitcherServer(createJ)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
}
