package handlers

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) Proxy(c *gin.Context) {
	room := c.MustGet(_cRoom).(avcontrol.RoomConfig)

	if room.Proxy.Host == "" || h.Host == "" || strings.EqualFold(h.Host, room.Proxy.Host) {
		c.Next()
		return
	}

	c.Abort()

	// make sure there is no cycle
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if strings.Contains(c.GetHeader(_hForwardedFor), ip) {
		c.String(http.StatusBadRequest, "detected proxy cycle. please try again")
		return
	}

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

	req.Header = c.Request.Header.Clone()

	// set X-Forwarded-For
	fwdFor := req.Header.Get(_hForwardedFor)
	if fwdFor == "" {
		req.Header.Set(_hForwardedFor, ip)
	} else {
		req.Header.Set(_hForwardedFor, fwdFor+", "+ip)
	}

	// set X-Request-ID
	if req.Header.Get(_hRequestID) == "" {
		req.Header.Set(_hRequestID, id)
	}

	// send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "unable to make proxy request: %s", err)
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get(_hContentType), resp.Body, nil)
}
