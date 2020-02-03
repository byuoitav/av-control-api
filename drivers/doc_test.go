package drivers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/qsc-driver"
	"github.com/spf13/pflag"
)

// This is a simple example of the main function of a driver server. It first builds the CreateDSPFunc for a QSC DSP using the qsc-driver library, then builds and starts the HTTP server with the standard DSP endpoints.
func Example() {
	create := func(ctx context.Context, addr string) (DSP, error) {
		return &qsc.DSP{
			Address: addr,
		}, nil
	}

	// bind a net.Listener to :8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	server, err := CreateDSPServer(create)
	if err != nil {
		// handle error
	}

	if err = server.Serve(lis); err != nil {
		// handle error
	}
}

// It makes sense to allow flags to be passed to a driver server to affect the behavior of the program. We like to use https://github.com/spf13/pflag for POSIX/GNU-style flags. This example lets the port be set by passing -p <port> or --port <port>.
func Example_flagsDriverServer() {
	var port int

	pflag.IntVarP(&port, "port", "p", 8080, "port to run the server on")
	pflag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// handle error
	}

	create := func(ctx context.Context, addr string) (Display, error) {
		return nec.NewProjector(addr, nec.WithDelay(300*time.Second))
	}

	server, err := CreateDisplayServer(create)
	if err != nil {
		// handle error
	}

	if err = server.Serve(lis); err != nil {
		// handle error
	}
}
