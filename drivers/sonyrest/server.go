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
	var psk string
	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.StringVarP(&psk, "psk", "k", "", "pre-shared key for the device")

	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(ctx context.Context, addr string) (drivers.Device, error) {
		return &sonyrest.TV{
			Address: addr,
			PSK:     psk,
		}, nil
	}

	server, err := drivers.NewServer(create)
	if err != nil {
		fmt.Printf("failed to create server: %s\n", err)
		os.Exit(1)
	}

	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
