package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
)

type goodGS struct{}
type badGS struct{}

func (g *goodGS) Get(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.StateResponse, error) {
	return avcontrol.StateResponse{
		Errors: []avcontrol.DeviceStateError{
			{
				ID: "just a filler",
			},
		},
	}, nil
}

func (g *goodGS) GetHealth(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomHealth, error) {
	return avcontrol.RoomHealth{}, errors.New("TODO")
}

func (g *goodGS) GetInfo(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomInfo, error) {
	return avcontrol.RoomInfo{}, errors.New("TODO")
}

func (g *goodGS) Set(ctx context.Context, room avcontrol.RoomConfig, req avcontrol.StateRequest) (avcontrol.StateResponse, error) {
	return avcontrol.StateResponse{
		Errors: []avcontrol.DeviceStateError{
			{
				ID: "just a filler",
			},
		},
	}, nil
}

func (g *badGS) Get(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.StateResponse, error) {
	return avcontrol.StateResponse{}, errors.New("no room to get")
}

func (g *badGS) Set(ctx context.Context, room avcontrol.RoomConfig, req avcontrol.StateRequest) (avcontrol.StateResponse, error) {
	return avcontrol.StateResponse{}, errors.New("no room to set")
}

func (g *badGS) GetHealth(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomHealth, error) {
	return avcontrol.RoomHealth{}, errors.New("TODO")
}

func (g *badGS) GetInfo(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomInfo, error) {
	return avcontrol.RoomInfo{}, errors.New("TODO")
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

	var room avcontrol.RoomConfig

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

	var res avcontrol.StateResponse

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

	var res avcontrol.StateResponse

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
