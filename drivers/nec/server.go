package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/nec-driver"
	"github.com/spf13/pflag"
)

func main() {
	// parse flags
	var port int

	// TODO add flags for timeout, etc
	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.Parse()

	// bind to given port
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(ctx context.Context, addr string) (drivers.Device, error) {
		return nec.NewProjector(addr, nec.WithDelay(300*time.Second)), err // TODO add options
	}

	// create server
	server, err := drivers.NewServer(create)

	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
