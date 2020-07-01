package handlers

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func proxyRequest(ctx context.Context, c echo.Context, url *url.URL, log *zap.Logger) error {
	log.Info("Proxying request", zap.String("url", url.String()))
	return nil
}

// TODO should we require the port to match?
func (h *Handlers) Proxy(c *gin.Context) {
	room := c.MustGet(_cRoom).(api.Room)

	if h.Host == "" || strings.EqualFold(h.Host, room.Proxy.Host) {
		c.Next()
		return
	}

	// TODO make sure there is no cycle

	c.Abort()

	id := c.GetString(_cRequestID)
	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	// proxy the request
	url := *room.Proxy
	url.Path = c.Request.URL.Path
	log.Info("Proxying request", zap.String("url", url.String()))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, url.String(), c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "unable to build proxy request: %s", err)
		return
	}

	// add request id header

	c.String(http.StatusOK, "hi")
}
