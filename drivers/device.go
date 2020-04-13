package drivers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/sync/singleflight"
)

type Device interface {
	GetInfo(ctx context.Context) (interface{}, error)
}

type CreateDeviceFunc func(context.Context, string) (Device, error)

func addDeviceRoutes(e *echo.Echo, create CreateDeviceFunc) {
	single := &singleflight.Group{}

	e.GET("/:address/info", func(c echo.Context) error { //
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the device")
		}

		val, err, _ := single.Do(addr, func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return d.GetInfo(c.Request().Context())
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, val)
	})
}
