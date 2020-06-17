package drivers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"golang.org/x/sync/singleflight"
)

var (
	ErrFuncNotSupported = errors.New("device does not support this function")
	ErrMissingAddress   = errors.New("must include address of the device")
	ErrMissingInput     = errors.New("missing input")
)

// CreateDeviceFunc is passed to CreateDeviceServer and is called to create a new Device struct whenever the Server needs to communicate  with a new Device.
type CreateDeviceFunc func(context.Context, string) (Device, error)

func CreateDeviceServer(create CreateDeviceFunc) (Server, error) {
	e := newEchoServer()
	m := &sync.Map{}
	single := &singleflight.Group{}

	newDev := func(ctx context.Context, addr string) (Device, error) {
		if dev, ok := m.Load(addr); ok {
			return dev, nil
		}

		dev, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}

		m.Store(addr, dev)
		return dev, nil
	}

	e.GET("/:address/GetPower", func(c echo.Context) error {
		addr := c.Param("address")
		if len(addr) == 0 {
			return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
		}

		val, err, _ := single.Do("GetPower"+addr, func() (interface{}, error) {
			dev, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			d, ok := dev.(DeviceWithPower)
			if !ok {
				return nil, ErrFuncNotSupported
			}

			return d.GetPower(c.Request().Context())
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// TODO return power struct
		return c.JSON(http.StatusOK, val)
	})

	e.GET("/:address/SetPower/:power", func(c echo.Context) error {
		addr := c.Param("address")
		pow, err := strconv.ParseBool(c.Param("power"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		_, err, _ = single.Do(fmt.Sprintf("SetPower%v%v", addr, pow), func() (interface{}, error) {
			dev, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			d, ok := dev.(DeviceWithPower)
			if !ok {
				return nil, ErrFuncNotSupported
			}

			return nil, d.SetPower(c.Request().Context(), pow)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// TODO return power struct
		return c.JSON(http.StatusOK, "")
	})

	// add DeviceWithAudioInput endpoints
	if _, ok := tmpl.(DeviceWithAudioInput); ok {
		e.GET("/:address/GetAudioInputs", func(c echo.Context) error {
			addr := c.Param("address")
			if len(addr) == 0 {
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			}

			val, err, _ := single.Do("GetAudioInputs"+addr, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithAudioInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return d.GetAudioInputs(c.Request().Context())
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, val)
		})

		e.GET("/:address/SetAudioInput/:output/:input", func(c echo.Context) error {
			addr := c.Param("address")
			output := c.Param("output")
			input := c.Param("input")
			switch {
			case len(addr) == 0:
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			case len(input) == 0:
				return c.String(http.StatusBadRequest, ErrMissingInput.Error())
			}

			_, err, _ := single.Do("SetAudioInput"+addr+output+input, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithAudioInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return nil, d.SetAudioInput(c.Request().Context(), output, input)
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, map[string]string{
				output: input,
			})
		})
	}

	// add DeviceWithVideoInput endpoints
	if _, ok := tmpl.(DeviceWithVideoInput); ok {
		e.GET("/:address/GetVideoInputs", func(c echo.Context) error {
			addr := c.Param("address")
			if len(addr) == 0 {
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			}

			val, err, _ := single.Do("GetVideoInputs"+addr, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithVideoInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return d.GetVideoInputs(c.Request().Context())
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, val)
		})

		e.GET("/:address/SetVideoInput/:output/:input", func(c echo.Context) error {
			addr := c.Param("address")
			output := c.Param("output")
			input := c.Param("input")
			switch {
			case len(addr) == 0:
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			case len(input) == 0:
				return c.String(http.StatusBadRequest, ErrMissingInput.Error())
			}

			_, err, _ := single.Do("SetAudioInput"+addr+output+input, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithVideoInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return nil, d.SetVideoInput(c.Request().Context(), output, input)
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, map[string]string{
				output: input,
			})
		})
	}

	// add DeviceWithAudioVideoInput endpoints
	if _, ok := tmpl.(DeviceWithAudioVideoInput); ok {
		e.GET("/:address/GetAudioVideoInputs", func(c echo.Context) error {
			addr := c.Param("address")
			if len(addr) == 0 {
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			}

			val, err, _ := single.Do("GetAudioVideoInputs"+addr, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithAudioVideoInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return d.GetAudioVideoInputs(c.Request().Context())
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, val)
		})

		e.GET("/:address/SetAudioVideoInput/:output/:input", func(c echo.Context) error {
			addr := c.Param("address")
			output := c.Param("output")
			input := c.Param("input")
			switch {
			case len(addr) == 0:
				return c.String(http.StatusBadRequest, ErrMissingAddress.Error())
			case len(input) == 0:
				return c.String(http.StatusBadRequest, ErrMissingInput.Error())
			}

			_, err, _ := single.Do("SetAudioVideoInput"+addr+output+input, func() (interface{}, error) {
				dev, err := create(c.Request().Context(), addr)
				if err != nil {
					return nil, err
				}

				d, ok := dev.(DeviceWithAudioVideoInput)
				if !ok {
					return nil, ErrNotSupported
				}

				return nil, d.SetAudioVideoInput(c.Request().Context(), output, input)
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, map[string]string{
				output: input,
			})
		})
	}

	return wrapEchoServer(e), nil
}
