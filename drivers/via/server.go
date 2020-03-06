package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	via "github.com/byuoitav/kramer-driver/via"
	"github.com/spf13/pflag"
)

// imports

func main() {
	var (
		port     int
		username string
		password string
	)

	pflag.IntVarP(&port, "port", "P", 80, "port to run the server on")
	pflag.StringVarP(&username, "username", "u", "su", "username for device")
	pflag.StringVarP(&password, "password", "p", "supass", "password for device")
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
	createDsp := func(ctx context.Context, addr string) (drivers.DSP, error) {
		return &via.VIA{
			Address:  addr,
			Username: username,
			Password: password,
		}, nil
	}

	// create server 
	server, err := drivers.CreateDSPServer(createDsp)
	if err != nil {
		fmt.Printf("Error while trying to create DSP Server: %s\n", err)
		os.Exit(1)
	}

	if err = server.Serve(lis); err != nil {
		fmt.Printf("failed to listen: %s\n", err)
		os.Exit(1)
	}
}
