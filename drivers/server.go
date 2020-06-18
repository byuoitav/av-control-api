package drivers

import (
	"net"
)

type Server interface {
	Serve(lis net.Listener) error
}
