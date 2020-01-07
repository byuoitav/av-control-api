package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/london-driver"
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
		return london.NewDSP(addr, london.WithDelay(300*time.Second)), nil
	}

	server, err := drivers.CreateDSPServer(create)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
