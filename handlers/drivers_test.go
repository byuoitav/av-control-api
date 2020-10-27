package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStatsPass(t *testing.T) {
	t.SkipNow()

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

	h.Stats(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	var states map[string]map[string]string

	err = json.Unmarshal(body, &states)
	if err != nil {
		t.Fatalf("error unmarshaling json: %s %s", err, body)
	}

	if states["driverStates"]["key"] != "val" {
		t.Fatalf("didn't work")
	}
}

func TestStatsFail(t *testing.T) {
	t.SkipNow()

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

	h.Stats(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading resp body: %s", err)
	}

	if string(body) != "no states to get" {
		t.Fatalf("unexpected response: %s", string(body))
	}
}
