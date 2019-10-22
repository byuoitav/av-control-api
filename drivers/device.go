package drivers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
)

type Device interface {
	GetInfo(ctx context.Context, addr string) (interface{}, error)
}

func addDeviceRoutes(e *echo.Echo, dev Device) {
	e.GET("/:address/info", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include an address of the display")
		}

		info, err := dev.GetInfo(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, info)
	})
}
