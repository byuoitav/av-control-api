package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"github.com/gin-gonic/gin"
)

type goodGS struct{}
type badGS struct{}

func (g *goodGS) Get(ctx context.Context, room api.Room) (api.StateResponse, error) {
	return api.StateResponse{
		Errors: []api.DeviceStateError{
			{
				ID: "just a filler",
			},
		},
	}, nil
}

func (g *goodGS) Set(ctx context.Context, room api.Room, req api.StateRequest) (api.StateResponse, error) {
	return api.StateResponse{
		Errors: []api.DeviceStateError{
			{
				ID: "just a filler",
			},
		},
	}, nil
}

func (g *goodGS) DriverStates(context.Context) (map[string]string, error) {
	return map[string]string{"key": "val"}, nil
}

func (g *badGS) Get(ctx context.Context, room api.Room) (api.StateResponse, error) {
	return api.StateResponse{}, errors.New("no room to get")
}

func (g *badGS) Set(ctx context.Context, room api.Room, req api.StateRequest) (api.StateResponse, error) {
	return api.StateResponse{}, errors.New("no room to set")
}

func (g *badGS) DriverStates(context.Context) (map[string]string, error) {
	return nil, errors.New("no states to get")
}

func TestGetRoomConfiguration(t *testing.T) {
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
	h.GetRoomConfiguration(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var room api.Room

	json.Unmarshal(body, &room)
	if room.ID != "ITB-1101" {
		t.Fatalf("correct room no returned: %v", body)
	}
}

func TestRoomStatePass(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://byu.edu",
		State:       &goodGS{},
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
	h.GetRoomState(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var res api.StateResponse

	json.Unmarshal(body, &res)
	if res.Errors[0].ID != "just a filler" {
		t.Fatalf("correct room no returned: %v", body)
	}
}

func TestRoomStateFail(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://byu.edu",
		State:       &badGS{},
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
	h.GetRoomState(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "no room to get" {
		t.Fatalf("unexpected response: %s", string(body))
	}
}

func TestSetRoomStatePass(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://byu.edu",
		State:       &goodGS{},
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
	h.SetRoomState(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var res api.StateResponse

	json.Unmarshal(body, &res)
	if res.Errors[0].ID != "just a filler" {
		t.Fatalf("correct room no returned: %v", body)
	}
}

func TestSetRoomStateFail(t *testing.T) {
	log := setLogger()
	defer log.Sync()

	d := goodDS{}
	h := Handlers{
		Logger:      log,
		DataService: &d,
		Host:        "http://byu.edu",
		State:       &badGS{},
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
	h.SetRoomState(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "no room to set" {
		t.Fatalf("wrong body received: %s", string(body))
	}
}
