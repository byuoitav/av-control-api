package drivers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
)

type Device interface {
	GetInfo(ctx context.Context) (interface{}, error)
}

type CreateDeviceFunc func(context.Context, string) (Device, error)

func addDeviceRoutes(e *echo.Echo, create CreateDeviceFunc, ctx context.Context) {
	e.GET("/:address/info", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the device")
		}

		dev, err := create(ctx, addr)
		if err != nil {
			return err
		}
		info, err := dev.GetInfo(c.Request().Context())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, info)
	})
}
