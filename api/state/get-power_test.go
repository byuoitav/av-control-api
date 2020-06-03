package state

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

func urlParse(rawurl string) *url.URL {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err.Error())
	}

	return url
}

type stateTest struct {
	name          string
	room          string
	deviceService api.DeviceService
	env           string
	resp          generateActionsResponse
}

var getPowerTest = []stateTest{
	stateTest{
		name:          "simple",
		deviceService: mock.SimpleRoom{},
		env:           "default",
		resp: generateActionsResponse{
			Actions: []action{
				action{
					ID: "ITB-1101-D1",
					Req: &http.Request{
						Method: http.MethodGet,
						URL:    urlParse("http://ITB-1101-D1.av/GetPower"),
					},
				},
			},
			ExpectedUpdates: 1,
		},
	},
}

func TestGetPower(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, tt := range getPowerTest {
		t.Run(tt.name, func(t *testing.T) {
			room, err := tt.deviceService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			var getPower getPower
			resp := getPower.GenerateActions(ctx, room, tt.env)

			if !cmp.Equal(resp, tt.resp) {
				t.Errorf("generated incorrect actions:\n\tgot %+v\n\texpected: %+v", tt.name)
			}
		})
	}
}
