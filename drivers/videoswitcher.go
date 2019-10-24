package drivers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type VideoSwitcher interface {
	Device

	// TODO notes about being 1 indexed

	GetInputByOutput(ctx context.Context, addr, output string) (string, error)
	SetInputByOutput(ctx context.Context, addr, output, input string) error

	// TODO active input ?
}

// TODO should we just make an explicit input/output struct that these return in their http calls?
func CreateVideoSwitcherServer(vs VideoSwitcher) Server {
	e := newEchoServer()

	addDeviceRoutes(e, vs)
	addVideoSwitcherRoutes(e, vs)

	return wrapEchoServer(e)
}

func addVideoSwitcherRoutes(e *echo.Echo, vs VideoSwitcher) {
	e.GET("/:address/output/:output/input", func(c echo.Context) error {
		addr := c.Param("address")
		output := c.Param("output")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the video switcher")
		case len(output) == 0:
			return c.String(http.StatusBadRequest, "must include an output port for the video switcher")
		}

		input, err := vs.GetInputByOutput(c.Request().Context(), addr, output)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Input{Input: fmt.Sprintf("%v:%v", input, output)})
	})

	e.GET("/:address/output/:output/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		output := c.Param("output")
		input := c.Param("input")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the video switcher")
		case len(output) == 0:
			return c.String(http.StatusBadRequest, "must include an output port")
		case len(input) == 0:
			return c.String(http.StatusBadRequest, "must include an input portr")
		}

		if err := vs.SetInputByOutput(c.Request().Context(), addr, output, input); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Input{Input: fmt.Sprintf("%v:%v", input, output)})
	})
}
