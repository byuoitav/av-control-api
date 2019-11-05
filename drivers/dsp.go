package drivers

import (
	"context"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
)

type dsp interface {
	GetVolumeByBlock(ctx context.Context, block string) (int, error)
	GetMutedByBlock(ctx context.Context, block string) (bool, error)

	SetVolumeByBlock(ctx context.Context, block string, volume int) error
	SetMutedByBlock(ctx context.Context, block string, muted bool) error
}

type DSP interface {
	Device
	dsp
}

type CreateDSPFunc func(context.Context, string) (DSP, error)

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

		d, err := create(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		volume, err := d.GetVolumeByBlock(c.Request().Context(), block)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Volume{Volume: volume})
	})

	e.GET("/:address/block/:block/volume/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		volume, err := strconv.Atoi(c.Param("volume"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		d, err := create(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if err = d.SetVolumeByBlock(c.Request().Context(), block, volume); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Volume{Volume: volume})
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

		d, err := create(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		muted, err := d.GetMutedByBlock(c.Request().Context(), block)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Muted{Muted: muted})
	})

	e.GET("/:address/block/:block/muted/:muted", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		muted, err := strconv.ParseBool(c.Param("muted"))
		switch {
		case len(addr) == 0:
			return c.String(http.StatusBadRequest, "must include the address of the dsp")
		case len(block) == 0:
			return c.String(http.StatusBadRequest, "must include a block for the dsp")
		case err != nil:
			return c.String(http.StatusBadRequest, err.Error())
		}

		d, err := create(c.Request().Context(), addr)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if err = d.SetMutedByBlock(c.Request().Context(), block, muted); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Muted{Muted: muted})
	})
}
