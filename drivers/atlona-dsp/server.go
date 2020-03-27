package main

import (
	"context"
	"fmt"
	"net"
	"os"

	atlona "github.com/byuoitav/atlona-driver"
	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/common/log"
	"github.com/spf13/pflag"
)

func main() {
	var port int
	var logLevel string
	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	nerr := log.SetLevel(logLevel)
	if nerr != nil {
		fmt.Printf("could not set log level: %v\n", nerr)
		os.Exit(1)
	}
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(ctx context.Context, addr string) (drivers.DSP, error) {
		return &atlona.Amp60{
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
