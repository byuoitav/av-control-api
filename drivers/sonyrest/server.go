package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/sonyrest-driver"
	"github.com/spf13/pflag"
)

func main() {
	var port int
	pflag.IntVarP(&port, "port", "p", 80, "port to run the server on")

	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(ctx context.Context, addr string) (drivers.Display, error) {
		return &sonyrest.TV{
			Address: addr,
			PSK:     "T3CL1T3",
		}, nil
	}

	server, err := drivers.CreateDisplayServer(create)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}

}
