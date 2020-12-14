package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/gin-gonic/gin"
)

type goodGS struct{}
type badGS struct{}

func boolP(b bool) *bool {
	return &b
}

func stringP(s string) *string {
	return &s
}

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
	return avcontrol.RoomHealth{
		Devices: map[avcontrol.DeviceID]avcontrol.DeviceHealth{
			"ITB-1101-D1": {
				Healthy: boolP(true),
			},
		},
	}, nil
}

func (g *goodGS) GetInfo(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomInfo, error) {
	return avcontrol.RoomInfo{
		Devices: map[avcontrol.DeviceID]avcontrol.DeviceInfo{
			"ITB-1101-D1": {
				Info: "here's the info",
			},
		},
	}, nil
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
	// return avcontrol.RoomHealth{}, errors.New("no room to get health")
	return avcontrol.RoomHealth{}, errors.New("can't get health")
}

func (g *badGS) GetInfo(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomInfo, error) {
	return avcontrol.RoomInfo{}, errors.New("can't get info")
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/state", nil)
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/", nil)
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

func TestRoomHealthPass(t *testing.T) {
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/health", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	h.GetRoomHealth(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var res avcontrol.RoomHealth

	json.Unmarshal(body, &res)

	if *res.Devices["ITB-1101-D1"].Healthy != true {
		t.Fatalf("correct room no returned: %s", string(body))
	}
}

func TestRoomHealthFail(t *testing.T) {
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/health", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	//getroomhealth no worky???
	h.GetRoomHealth(c)
	fmt.Printf("helloo1111\n\n\n\n")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "can't get health" {
		t.Fatalf("unexpected response: %s", string(body))
	}
}

func TestRoomInfoPass(t *testing.T) {
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/info", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	h.GetRoomInfo(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var res avcontrol.RoomInfo

	json.Unmarshal(body, &res)

	if res.Devices["ITB-1101-D1"].Info != "here's the info" {
		t.Fatalf("correct room no returned: %s", string(body))
	}
}

func TestRoomInfoFail(t *testing.T) {
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
	c.Request, _ = http.NewRequest(http.MethodGet, "/room/:room/health", nil)
	c.Params = gin.Params{
		{
			Key:   "room",
			Value: "ITB-1101",
		},
	}

	h.RequestID(c)
	h.Room(c)
	h.GetRoomInfo(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "can't get info" {
		t.Fatalf("unexpected response: %s", string(body))
	}
}
