package main

import (
	"net"
)

func main() {
	// display := &nec.Display{}

	// server := drivers.CreateDisplayServer(display)

	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
	}

	if err = server.Serve(lis); err != nil {
	}
}
