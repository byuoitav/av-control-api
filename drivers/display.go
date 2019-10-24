package drivers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Display interface {
	Device

	GetPower(ctx context.Context, addr string) (string, error)
	GetBlanked(ctx context.Context, addr string) (bool, error)
	GetInput(ctx context.Context, addr string) (string, error)
	GetActiveSignal(ctx context.Context, addr string) (bool, error)

	SetPower(ctx context.Context, addr, power string) error
	SetBlanked(ctx context.Context, addr string, blanked bool) error
	SetInput(ctx context.Context, addr string, input string) error
}

func CreateDisplayServer(disp Display) Server {
	e := newEchoServer()

	addDeviceRoutes(e, disp)
	addDisplayRoutes(e, disp)

	return wrapEchoServer(e)
}

func addDisplayRoutes(e *echo.Echo, disp Display) {
	// power
	e.GET("/:address/power", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the display")
		}

		power, err := disp.GetPower(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Power{Power: power})
	})

	e.GET("/:address/power/:power", func(c echo.Context) error {
		addr := c.Param("address")
		power := c.Param("power")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the display")
		case len(power) == 0:
			return c.String(http.StatusBadRequest, "must include a power state to set")
		}

		if err := disp.SetPower(c.Request().Context(), addr, power); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Power{Power: power})
	})

	// blanked
	e.GET("/:address/blanked", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the display")
		}

		blanked, err := disp.GetBlanked(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Blanked{Blanked: blanked})
	})

	e.GET("/:address/blanked/:blanked", func(c echo.Context) error {
		addr := c.Param("address")
		blanked, err := strconv.ParseBool(c.Param("blanked"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the display")
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := disp.SetBlanked(c.Request().Context(), addr, blanked); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Blanked{Blanked: blanked})
	})

	// input
	e.GET("/:address/input", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the display")
		}

		input, err := disp.GetInput(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Input{Input: input})
	})

	e.GET("/:address/input/:input", func(c echo.Context) error {
		addr := c.Param("address")
		input := c.Param("input")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the display")
		case len(input) == 0:
			return c.String(http.StatusBadRequest, "must include a input to set")
		}

		if err := disp.SetInput(c.Request().Context(), addr, input); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Input{Input: input})
	})

	// active signal
	e.GET("/:address/activesignal", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the display")
		}

		asignal, err := disp.GetActiveSignal(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, ActiveSignal{ActiveSignal: asignal})
	})
}
