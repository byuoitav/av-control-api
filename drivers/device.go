package drivers

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"golang.org/x/sync/singleflight"
)

type CreateDeviceFunc func(context.Context, string) (interface{}, error)

func CreateServer(tmpl interface{}, create CreateDeviceFunc) Server {
	e := newEchoServer()
	m := &sync.Map{}
	single := &singleflight.Group{}

	dev := func(ctx context.Context, addr string) (interface{}, error) {
		if d, ok := m.Load(addr); ok {
			return d, nil
		}

		d, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}

		m.Store(addr, d)
		return d, nil
	}

	// TODO add all endpoints

	if dai, ok := tmpl.(DeviceWithAudioInput); ok {
		e.GET("/:address/info", func(c echo.Context) error { //
			addr := c.Param("address")
			if len(addr) == 0 {
				return c.String(http.StatusBadRequest, "must include the address of the device")
			}

			val, err, _ := single.Do("0"+addr, func() (interface{}, error) {
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
}

func addDeviceRoutes(e *echo.Echo, create CreateDeviceFunc) {
	single := &singleflight.Group{}

	e.GET("/:address/info", func(c echo.Context) error { //
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, "must include the address of the device")
		}

		val, err, _ := single.Do("0"+addr, func() (interface{}, error) {
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
