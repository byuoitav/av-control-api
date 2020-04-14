package drivers

import (
	"context"
	"fmt"
	"internal/singleflight"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

type VideoSwitcher interface {
	Device
	// TODO notes about being 1 indexed

	// TODO should we just make an explicit input/output struct that these return in their http calls?
	GetInputByOutput(ctx context.Context, output string) (string, error)
	SetInputByOutput(ctx context.Context, output, input string) error

	// TODO active input ?
}

type CreateVideoSwitcherFunc func(context.Context, string) (VideoSwitcher, error)

func CreateVideoSwitcherServer(create CreateVideoSwitcherFunc) Server {
	e := newEchoServer()
	m := &sync.Map{}

	vs := func(ctx context.Context, addr string) (VideoSwitcher, error) {
		if vs, ok := m.Load(addr); ok {
			return vs.(VideoSwitcher), nil
		}

		vs, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}

		m.Store(addr, vs)
		return vs, nil
	}

	dev := func(ctx context.Context, addr string) (Device, error) {
		return vs(ctx, addr)
	}

	addDeviceRoutes(e, dev)
	addVideoSwitcherRoutes(e, vs)

	return wrapEchoServer(e)
}

func addVideoSwitcherRoutes(e *echo.Echo, create CreateVideoSwitcherFunc) {
	single := &singleflight.Group{}

	e.GET("/:address/output/:output/input", func(c echo.Context) error {
		addr := c.Param("address")
		out := c.Param("output")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the video switcher")
		case len(out) == 0:
			return c.String(http.StatusBadRequest, "must include an output port for the video switcher")
		}

		val, err, _ := single.Do(addr+out+"input", func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return d.GetInputByOutput(c.Request().Context(), out)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		in, ok := val.(string)
		if !ok {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unexpected response: expected %T, got: %T", in, val))
		}

		return c.JSON(http.StatusOK, input{Input: fmt.Sprintf("%v:%v", in, out)})
	})

	e.GET("/:address/output/:output/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		out := c.Param("output")
		in := c.Param("input")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the video switcher")
		case len(out) == 0:
			return c.String(http.StatusBadRequest, "must include an output port")
		case len(in) == 0:
			return c.String(http.StatusBadRequest, "must include an input port")
		}

		_, err, _ := single.Do(fmt.Sprintf("%v%v%v", addr, out, in), func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return nil, d.SetInputByOutput(c.Request().Context(), out, in)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, input{Input: fmt.Sprintf("%v:%v", in, out)})
	})
}
