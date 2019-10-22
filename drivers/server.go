package drivers

import (
	"net"

	"github.com/labstack/echo"
)

type Server interface {
	Serve(lis net.Listener) error
}

type wrappedEchoServer struct {
	*echo.Echo
}

func (e *wrappedEchoServer) Serve(lis net.Listener) error {
	return e.Server.Serve(lis)
}

func wrapEchoServer(e *echo.Echo) Server {
	return &wrappedEchoServer{
		Echo: e,
	}
}
