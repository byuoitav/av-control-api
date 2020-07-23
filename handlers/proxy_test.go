package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setLogger() *zap.Logger {
	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zapcore.Level(1)),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	log, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}

	return log
}

func TestProxyWithCycle(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://byu.edu",
	}

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	h.Proxy(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "detected proxy cycle. please try again" {
		t.Fatalf("did not error on cycle")
	}
}

func TestProxyWithFwdFor(t *testing.T) {

	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://yourmom.com",
	}

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}
	c.Request.RemoteAddr = "http://byu.edu"
	c.Request.Header.Set(_hForwardedFor, "yourmom")

	h.RequestID(c)
	h.Room(c)
	h.Proxy(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) == "detected proxy cycle. please try again" {
		t.Fatalf("errored on cycle")
	}

}

func TestProxyWithoutFwdFor(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://yourmom.com",
	}

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}
	c.Request.RemoteAddr = "http://byu.edu"

	h.RequestID(c)
	h.Room(c)
	h.Proxy(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) == "detected proxy cycle. please try again" {
		t.Fatalf("errored on cycle")
	}
}

func TestProxyWithoutHost(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "",
	}

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	h.Proxy(c)
}
