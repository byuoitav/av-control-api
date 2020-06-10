package state

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api/log"
	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

var getInputTest = []stateTest{
	{
		name:        "simpleSeparateInput",
		dataService: &mock.SimpleSeparateInput{},
		env:         "default",
		resp: generatedActions{
			Actions: []action{
				{
					ID:  "ITB-1101-D1",
					Req: newRequest(http.MethodGet, "http://ITB-1101-CP1.byu.edu/ITB-1101-D1.av/GetAVInput"),
				},
				{
					ID:  "ITB-1101-SW1",
					Req: newRequest(http.MethodGet, "http://ITB-1101-CP1.byu.edu/ITB-1101-SW1.av/GetVideoInput"),
				},
				{
					ID:  "ITB-1101-SW1",
					Req: newRequest(http.MethodGet, "http://ITB-1101-CP1.byu.edu/ITB-1101-SW1.av/GetAudioInput"),
				},
			},
			ExpectedUpdates: 1,
		},
	},
}

func TestGetInput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, tt := range getInputTest {
		t.Run(tt.name, func(t *testing.T) {
			room, err := tt.dataService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			get := getInput{
				Logger:      log.Logger{},
				Environment: tt.env,
			}

			resp := get.GenerateActions(ctx, room)

			if diff := cmp.Diff(tt.resp, resp); diff != "" {
				t.Errorf("generated incorrect actions (-want, +got):\n%s", diff)
			}
		})
	}
}
