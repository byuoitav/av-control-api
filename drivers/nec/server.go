package main

import (
	"fmt"
	"net"
	"os"

	"github.com/byuoitav/av-control-api/drivers"
	"github.com/byuoitav/nec-driver"
	"github.com/spf13/pflag"
)

func main() {
	// parse flags
	var port int

	// TODO add flags for timeout, etc
	pflag.IntVarP(&port, "port", "p", 80, "port to run the server on")
	pflag.Parse()

	// bind to given port
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	create := func(addr string) drivers.Display {
		return nec.NewProjector(addr) // TODO add options
	}

	// create server
	server := drivers.CreateDisplayServer(create)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
