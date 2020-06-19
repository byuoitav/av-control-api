package drivers

import (
	"context"
	"net"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echopprof "github.com/sevenNt/echo-pprof"
)

type wrappedEchoServer struct {
	*echo.Echo
}

func newEchoServer() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	// /healthz simply reports on the server being up and thus returns
	// a generic string rather than real health data. Real health data
	// will be left up to each individual driver as necessary.
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	echopprof.Wrap(e)

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

func (e *wrappedEchoServer) Stop(ctx context.Context) error {
	return e.Shutdown(ctx)
}

type power struct {
	Power string `json:"power"`
}

type blanked struct {
	Blanked bool `json:"blanked"`
}

type input struct {
	Input string `json:"input"`
}

type muted struct {
	Muted bool `json:"muted"`
}

type volume struct {
	Volume int `json:"volume"`
}

type activeSignal struct {
	ActiveSignal bool `json:"activeSignal"`
}
