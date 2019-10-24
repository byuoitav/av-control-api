package drivers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type dsp interface {
	GetVolumeByBlock(ctx context.Context, addr, block string) (int, error)
	GetMutedByBlock(ctx context.Context, addr, block string) (bool, error)

	SetVolumeByBlock(ctx context.Context, addr, block string, volume int) error
	SetMutedByBlock(ctx context.Context, addr, block string, muted bool) error
}

type DSP interface {
	Device
	dsp
}

func CreateDSPServer(dsp DSP) Server {
	e := newEchoServer()

	addDeviceRoutes(e, dsp)
	addDSPRoutes(e, dsp)

	return wrapEchoServer(e)
}

func addDSPRoutes(e *echo.Echo, d dsp) {
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

		volume, err := d.GetVolumeByBlock(c.Request().Context(), addr, block)
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

		if err = d.SetVolumeByBlock(c.Request().Context(), addr, block, volume); err != nil {
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

		muted, err := d.GetMutedByBlock(c.Request().Context(), addr, block)
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

		if err = d.SetMutedByBlock(c.Request().Context(), addr, block, muted); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, Muted{Muted: muted})
	})
}
