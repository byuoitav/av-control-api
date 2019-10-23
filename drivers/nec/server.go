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

	pflag.IntVarP(&port, "port", "p", 80, "port to run the server on")
	pflag.Parse()

	// bind to given port
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	// import display lib
	display := &nec.Projector{}

	// create server
	server := drivers.CreateDisplayServer(display)
	if err = server.Serve(lis); err != nil {
		fmt.Printf("error while listening: %s\n", err)
		os.Exit(1)
	}
}
