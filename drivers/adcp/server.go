package main

import (
	"context"
	"fmt"
	"net"
	"os"

	adcp "github.com/byuoitav/adcp-driver"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/spf13/pflag"
)

func main() {
	// parse flags
	var port int

	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.Parse()

	// bind to given port
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	// import display lib
	display := func(ctx context.Context, addr string) (drivers.DisplayDSP, error) {
		return &adcp.Projector{
			Address: addr,
		}, nil
	}

	// create server
	server := drivers.CreateDisplayDSPServer(display)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
