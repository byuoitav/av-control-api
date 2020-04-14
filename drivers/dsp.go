package drivers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"golang.org/x/sync/singleflight"
)

type DSP interface {
	Device
	GetVolumeByBlock(ctx context.Context, block string) (int, error)
	GetMutedByBlock(ctx context.Context, block string) (bool, error)

	SetVolumeByBlock(ctx context.Context, block string, volume int) error
	SetMutedByBlock(ctx context.Context, block string, muted bool) error
}

/*
DSP is an interface with the methods required for a DSP library to implement. It is a combination of the Device interface as well as DSP specific functions. The API will send volume levels between 0 and 100, inclusive. Drivers implementing this interface should adjust the [0-100] volume level to the appropriate level for the device.

A driver library implementing this interface should look something like this:
	type QSC struct {
		Address string
		Username string
		Password string
	}

	func (q *QSC) GetInfo(ctx context.Context) (interface{}, error) {
		// open a connection with the dsp, return some info about the device...
	}

	func (q *QSC) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
		// open a connection with the dsp, return the volume for on block...
	}

	func (q *QSC) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
		// open a connection with the dsp, return the muted status for block...
	}

	func (q *QSC) SetVolumeByBlock(ctx context.Context, block string, volume int) (error) {
		// open a connection with the dsp, set the volume on block...
	}

	func (q *QSC) SetMutedByBlock(ctx context.Context, block string, muted bool) (error) {
		// open a connection with the dsp, set the muted status on block...
	}
*/

// CreateDSPFunc is passed to CreateDSPServer and is called to create a new DSP struct whenever the Server needs to communicate with a new DSP address.
type CreateDSPFunc func(ctx context.Context, addr string) (DSP, error)

// CreateDSPServer returns a Server with the appropriate endpoints for a DSP added to it.
func CreateDSPServer(create CreateDSPFunc) (Server, error) {
	e := newEchoServer()
	m := &sync.Map{}

	dsp := func(ctx context.Context, addr string) (DSP, error) {
		if dsp, ok := m.Load(addr); ok {
			return dsp.(DSP), nil
		}

		dsp, err := create(ctx, addr)
		if err != nil {
			return nil, err
		}
		m.Store(addr, dsp)
		return dsp, nil
	}

	dev := func(ctx context.Context, addr string) (Device, error) {
		return dsp(ctx, addr)
	}

	addDeviceRoutes(e, dev)
	addDSPRoutes(e, dsp)

	return wrapEchoServer(e), nil
}

func addDSPRoutes(e *echo.Echo, create CreateDSPFunc) {
	single := &singleflight.Group{}

	// volume
	e.GET("/:address/block/:block/volume", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		}

		val, err, _ := single.Do("0"+addr+block, func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return d.GetVolumeByBlock(c.Request().Context(), block)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		vol, ok := val.(int)
		if !ok {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unexpected response: expected %T, got: %T", vol, val))
		}

		return c.JSON(http.StatusOK, volume{Volume: vol})
	})

	e.GET("/:address/block/:block/volume/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		vol, err := strconv.Atoi(c.Param("volume"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		_, err, _ = single.Do(fmt.Sprintf("1%v%v%v", addr, block, vol), func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return nil, d.SetVolumeByBlock(c.Request().Context(), block, vol)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, volume{Volume: vol})
	})

	// muting
	e.GET("/:address/block/:block/muted", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		}

		val, err, _ := single.Do("2"+addr+block, func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return d.GetMutedByBlock(c.Request().Context(), block)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		mute, ok := val.(bool)
		if !ok {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("unexpected response: expected %T, got: %T", mute, val))
		}

		return c.JSON(http.StatusOK, muted{Muted: mute})
	})

	e.GET("/:address/block/:block/muted/:muted", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		mute, err := strconv.ParseBool(c.Param("muted"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		_, err, _ = single.Do(fmt.Sprintf("3%v%v%v", addr, block, mute), func() (interface{}, error) {
			d, err := create(c.Request().Context(), addr)
			if err != nil {
				return nil, err
			}

			return nil, d.SetMutedByBlock(c.Request().Context(), block, mute)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, muted{Muted: mute})
	})
}
