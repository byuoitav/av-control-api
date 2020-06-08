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

type getStateTest struct {
	name          string
	room          string
	env           string
	deviceService interface {
		api.DeviceService
		SetBaseURL(string)
	}

	httpResps map[string]string
	apiResp   api.StateResponse
}

var getTests = []getStateTest{
	getStateTest{
		name:          "Simple",
		deviceService: &mock.SimpleRoom{},
		env:           "default",
		httpResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":   `{"power": "on"}`,
			"/ITB-1101-D1.av/GetAVInput": `{"input": "hdmi!1"}`,
			"/ITB-1101-D1.av/GetBlanked": `{"blanked": false}`,
			"/ITB-1101-D1.av/GetMuted":   `{"muted": false}`,
			"/ITB-1101-D1.av/GetVolume":  `{"volume": 30}`,
		},
		apiResp: api.StateResponse{
			OutputGroups: map[api.DeviceID]api.OutputGroupState{
				"ITB-1101-D1": api.OutputGroupState{
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-VIA1"),
						Video:            deviceID("ITB-1101-VIA1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-VIA1"),
							api.DeviceID("ITB-1101-HDMI1"),
						},
					},
					Volume: intP(30),
					Muted:  boolP(false),
				},
			},
		},
	},
	getStateTest{
		name:          "Simple2",
		deviceService: &mock.SimpleRoom{},
		env:           "default",
		httpResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":   `{"power": "standby"}`,
			"/ITB-1101-D1.av/GetAVInput": `{"input": "hdmi!2"}`,
			"/ITB-1101-D1.av/GetBlanked": `{"blanked": true}`,
			"/ITB-1101-D1.av/GetMuted":   `{"muted": true}`,
			"/ITB-1101-D1.av/GetVolume":  `{"volume": 100}`,
		},
		apiResp: api.StateResponse{
			OutputGroups: map[api.DeviceID]api.OutputGroupState{
				"ITB-1101-D1": api.OutputGroupState{
					PoweredOn: boolP(false),
					Blanked:   boolP(true),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-HDMI1"),
						Video:            deviceID("ITB-1101-HDMI1"),
						CanSetSeparately: boolP(true),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(100),
					Muted:  boolP(true),
				},
			},
		},
	},
}

func TestGetState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			// start http server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, tt.httpResps[r.URL.Path])
			}))
			t.Cleanup(func() {
				ts.Close()
			})

			tt.deviceService.SetBaseURL(ts.URL)

			room, err := tt.deviceService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			resp, err := GetDevices(ctx, room, tt.env)
			if err != nil {
				t.Errorf("unable to get room state: %s", err)
			}

			if diff := cmp.Diff(tt.apiResp, resp); diff != "" {
				t.Errorf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}
