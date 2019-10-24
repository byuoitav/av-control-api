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

func newEchoServer() *echo.Echo {
	e := echo.New()

	return e
}

func wrapEchoServer(e *echo.Echo) Server {
	return &wrappedEchoServer{
		Echo: e,
	}
}

func (e *wrappedEchoServer) Serve(lis net.Listener) error {
	return e.Server.Serve(lis)
}
