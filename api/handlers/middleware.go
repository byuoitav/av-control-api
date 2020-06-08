package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/segmentio/ksuid"
)

const (
	_cRequestID = "requestID"
)

type Middleware struct{}

func (m *Middleware) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := ksuid.NewRandom()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		c.Set(_cRequestID, uid.String())
		return next(c)
	}
}
