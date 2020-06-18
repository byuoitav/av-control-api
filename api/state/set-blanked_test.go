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

var setBlankedTest = []stateTest{
	{
		name:        "Simple",
		dataService: &mock.SimpleRoom{},
		env:         "default",
		resp: generatedActions{
			Actions: []action{
				{
					ID:  "ITB-1101-D1",
					Req: newRequest(http.MethodGet, "http://host/ITB-1101-D1.av/SetBlanked/true"),
				},
			},
			ExpectedUpdates: 1,
		},
		req: api.StateRequest{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					Blanked: boolP(true),
				},
			},
		},
	},
}

func TestSetBlanked(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, tt := range setBlankedTest {
		t.Run(tt.name, func(t *testing.T) {
			room, err := tt.dataService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			set := setBlanked{
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
