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

var getVolumeTest = []stateTest{
	{
		name: "simpleSeparateInput",
		dataService: &mock.SimpleSeparateInput{
			BaseURL: "http://host",
		},
		env: "default",
		resp: generatedActions{
			Actions: []action{
				{
					ID:  "ITB-1101-AMP1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-AMP1.av/GetVolume"),
				},
				{
					ID:  "ITB-1101-VIA1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-VIA1.av/GetVolume"),
				},
			},
			ExpectedUpdates: 2,
		},
	},
}

func TestGetVolume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, tt := range getVolumeTest {
		t.Run(tt.name, func(t *testing.T) {
			room, err := tt.dataService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			get := getVolume{
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
