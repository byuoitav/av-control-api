package state

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

type setStateTest struct {
	name          string
	room          string
	env           string
	deviceService interface {
		api.DeviceService
		SetBaseURL(string)
	}

	apiReq      api.StateRequest
	driverResps map[string]string
	apiResp     api.StateResponse
}

var setTests = []setStateTest{
	setStateTest{
		name:          "Simple/Power",
		env:           "default",
		deviceService: &mock.SimpleRoom{},
		apiReq: api.StateRequest{
			OutputGroups: map[api.DeviceID]api.OutputGroupState{
				"ITB-1101-D1": api.OutputGroupState{
					PoweredOn: boolP(true),
				},
			},
		},
		driverResps: map[string]string{
			"/ITB-1101-D1.av/SetPower/on": `{"power": "on"}`,
		},
		apiResp: api.StateResponse{
			OutputGroups: map[api.DeviceID]api.OutputGroupState{
				"ITB-1101-D1": api.OutputGroupState{
					PoweredOn: boolP(true),
				},
			},
		},
	},
}

func TestSetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			// start http server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, tt.driverResps[r.URL.Path])
			}))
			t.Cleanup(func() {
				ts.Close()
			})

			tt.deviceService.SetBaseURL(ts.URL)

			room, err := tt.deviceService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			resp, err := SetDevices(ctx, tt.apiReq, room, tt.env)
			if err != nil {
				t.Errorf("unable to get room state: %s", err)
			}

			if diff := cmp.Diff(tt.apiResp, resp); diff != "" {
				t.Errorf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}
