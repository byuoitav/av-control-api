package handlers

import (
	"context"
	"net/url"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// TODO should we require the port to match?
func (h *Handlers) shouldProxy(room api.Room) bool {
	if h.Host == "" || strings.EqualFold(h.Host, room.Proxy.Host) {
		return false
	}

	return true
}

func proxyRequest(ctx context.Context, c echo.Context, url *url.URL, log *zap.Logger) error {
	log.Info("Proxying request", zap.String("url", url.String()))
	return nil
}
