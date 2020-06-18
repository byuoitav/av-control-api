package state

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/log"
	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

var devID1 = "ITB-1101-VIA1"
var devID2 = "ITB-1101-HDMI1"

var setInputTest = []stateTest{
	{
		name:        "Simple",
		dataService: &mock.SimpleRoom{},
		env:         "default",
		resp: generatedActions{
			Actions: []action{
				{
					ID:  "ITB-1101-D1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-D1.av/SetAVInput/hdmi!1"),
				},
			},
			ExpectedUpdates: 1,
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Input: map[string]api.Input{
						"hdmi1": api.Input{
							Video: &devID1,
							Audio: &devID1,
						},
					},
				},
			},
		},
	},
	{
		name:        "SimpleSeparate",
		dataService: &mock.SimpleSeparateInput{},
		env:         "default",
		resp: generatedActions{
			Actions: []action{
				{
					ID:  "ITB-1101-D1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-D1.av/SetAVInput/hdmi!2"),
				},
				{
					ID:  "ITB-1101-SW1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-SW1.av/SetVideoInput/2"),
				},
				{
					ID:  "ITB-1101-SW1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-SW1.av/SetAudioInput/1"),
				},
			},
			ExpectedUpdates: 1,
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Input: map[string]api.Input{
						"hdmi1": api.Input{
							Video: &devID2,
							Audio: &devID1,
						},
					},
				},
			},
		},
	},
}

func TestSetInput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, tt := range setInputTest {
		t.Run(tt.name, func(t *testing.T) {
			room, err := tt.dataService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			set := setInput{
				Logger:      log.Logger{},
				Environment: tt.env,
			}

			resp := set.GenerateActions(ctx, room, tt.req)

			if diff := cmp.Diff(tt.resp, resp); diff != "" {
				t.Errorf("generated incorrect actions (-want, +got):\n%s", diff)
			}
		})
	}
}
