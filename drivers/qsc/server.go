package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/qsc-driver"
	"github.com/spf13/pflag"
)

func main() {
	var port int
	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")

	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(ctx context.Context, addr string) (drivers.DSP, error) {
		return &qsc.DSP{
			Address: addr,
		}, nil
	}

	server, err := drivers.CreateDSPServer(create)
	if err != nil {
		fmt.Printf("Error while trying to create DSP Server: %s\n", err)
		os.Exit(1)
	}

	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}

}
