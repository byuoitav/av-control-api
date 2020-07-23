package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"github.com/gin-gonic/gin"
)

type goodDS struct{}
type badDS struct{}

func (d *goodDS) Room(ctx context.Context, id string) (api.Room, error) {
	return api.Room{
		ID: "ITB-1101",
		Proxy: &url.URL{
			Scheme: "http",
			Host:   "byu.edu",
			Path:   "/room/ITB-1101",
		},
	}, nil
}

func (d *goodDS) DriverMapping(ctx context.Context) (api.DriverMapping, error) {
	return api.DriverMapping{}, nil
}

func (d *badDS) Room(ctx context.Context, id string) (api.Room, error) {
	return api.Room{}, errors.New("no room")
}

func (d *badDS) DriverMapping(ctx context.Context) (api.DriverMapping, error) {
	return api.DriverMapping{}, nil
}

func TestRequestIDWithID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Request.Header.Set(_hRequestID, "ID")

	handler := Handlers{}

	handler.RequestID(c)

	id := c.Keys[_cRequestID]
	if id != "ID" {
		t.Fatalf("request header changed when it shouldn't have: %s", id)
	}
}

func TestRequestIDWithOutID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler := Handlers{}

	handler.RequestID(c)

	id := c.Keys[_cRequestID]
	if id == "" {
		t.Fatalf("request header didn't change when it should have: %s", id)
	}
}

func TestLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler := Handlers{}
	log := setLogger()
	defer log.Sync()
	handler.Logger = log

	handler.RequestID(c)
	handler.Log(c)
}

func TestRoomWithRoomNoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	log := setLogger()
	defer log.Sync()
	d := goodDS{}

	handler := Handlers{
		Logger:      log,
		DataService: &d,
	}

	handler.RequestID(c)
	handler.Room(c)

	room := c.MustGet(_cRoom).(api.Room)
	if room.ID != "ITB-1101" {
		t.Fatalf("incorrect room gotten %s", c.Keys[_cRoom].(api.Room).ID)
	}
}

func TestRoomWithRoomError(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	log := setLogger()
	defer log.Sync()
	d := badDS{}

	handler := Handlers{
		Logger:      log,
		DataService: &d,
	}

	handler.RequestID(c)
	handler.Room(c)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if c.Keys[_cRoom] != nil {
		t.Fatalf("correctly got room when it shouldn't have")
	}

	if string(body) != "unable to get room no room" {
		t.Fatalf("wrong error generated: %s", body)
	}
}

func TestRoomWithoutRoom(t *testing.T) {
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler := Handlers{}
	log := setLogger()
	defer log.Sync()
	handler.Logger = log

	handler.RequestID(c)
	handler.Room(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if c.Keys[_cRoom] != nil {
		t.Fatalf("correctly got room when it shouldn't have")
	}

	if string(body) != "must include room" {
		t.Fatalf("wrong error generated: %s", body)
	}
}
