package state

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

var getInputTest = []stateTest{
	stateTest{
		name:          "simpleSeparateInput",
		deviceService: mock.SimpleSeparateInput{},
		env:           "default",
		resp: generateActionsResponse{
			Actions: []action{
				action{
					ID:  "ITB-1101-D1",
					Req: newRequest(http.MethodGet, "http://ITB-1101-CP1.byu.edu/ITB-1101-D1.av/GetAVInput"),
				},
				action{
					ID:  "ITB-1101-SW1",
					Req: newRequest(http.MethodGet, "http://ITB-1101-CP1.byu.edu/ITB-1101-SW1.av/GetVideoInput"),
				},
				action{
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
			room, err := tt.deviceService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			var get getInput
			resp := get.GenerateActions(ctx, room, tt.env)

			if diff := cmp.Diff(tt.resp, resp); diff != "" {
				t.Errorf("generated incorrect actions (-want, +got):\n%s", diff)
			}
		})
	}
}

func Equals(r1, r2 generateActionsResponse) bool {
	if len(r1.Actions) != len(r2.Actions) || len(r1.Errors) != len(r2.Errors) || r1.ExpectedUpdates != r2.ExpectedUpdates {
		return false
	}

	for i := range r1.Actions {
		if r1.Actions[i].ID != r2.Actions[i].ID {
			return false
		}
		if r1.Actions[i].Req.Method != r2.Actions[i].Req.Method {
			return false
		}
		// urls doesn't work and idk why
		if r1.Actions[i].Req.URL != r2.Actions[i].Req.URL {
			fmt.Printf("bad urls: %v %v\n", r1.Actions[i].Req.URL, r2.Actions[i].Req.URL)
			return false
		}
	}

	return true
}
