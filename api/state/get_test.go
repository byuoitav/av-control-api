package state

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/api/log"
	"github.com/byuoitav/av-control-api/api/mock"
	"github.com/google/go-cmp/cmp"
)

type getStateTest struct {
	name        string
	room        string
	env         string
	dataService interface {
		api.DataService
		SetBaseURL(string)
	}

	driverResps map[string]string
	apiResp     api.StateResponse
}

var getTests = []getStateTest{
	{
		name:        "Simple/1",
		dataService: &mock.SimpleRoom{},
		env:         "default",
		driverResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":    `{"power": "on"}`,
			"/ITB-1101-D1.av/GetAVInput":  `{"input": "hdmi!1"}`,
			"/ITB-1101-D1.av/GetBlanked":  `{"blanked": false}`,
			"/ITB-1101-D1.av/GetMuted":    `{"muted": false}`,
			"/ITB-1101-D1.av/GetVolume":   `{"volume": 30}`,
			"/ITB-1101-VIA1.av/GetVolume": `{"volume": 50}`,
			"/ITB-1101-VIA1.av/GetMuted":  `{"muted": true}`,
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-VIA1"),
						Video:            deviceID("ITB-1101-VIA1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(30),
					Muted:  boolP(false),
				},
				"ITB-1101-VIA1": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
			},
		},
	},
	{
		name:        "Simple/2",
		dataService: &mock.SimpleRoom{},
		env:         "default",
		driverResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":    `{"poweredOn": false}`,
			"/ITB-1101-D1.av/GetAVInput":  `{"input": "hdmi!2"}`,
			"/ITB-1101-D1.av/GetBlanked":  `{"blanked": true}`,
			"/ITB-1101-D1.av/GetMuted":    `{"muted": true}`,
			"/ITB-1101-D1.av/GetVolume":   `{"volume": 100}`,
			"/ITB-1101-VIA1.av/GetVolume": `{"volume": 50}`,
			"/ITB-1101-VIA1.av/GetMuted":  `{"muted": true}`,
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(false),
					Blanked:   boolP(true),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-HDMI1"),
						Video:            deviceID("ITB-1101-HDMI1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(100),
					Muted:  boolP(true),
				},
				"ITB-1101-VIA1": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
			},
		},
	},
	{
		name:        "SimpleSeparateInput/1",
		dataService: &mock.SimpleSeparateInput{},
		env:         "default",
		driverResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":       `{"poweredOn": true}`,
			"/ITB-1101-D1.av/GetAVInput":     `{"input": "hdmi!2"}`,
			"/ITB-1101-D1.av/GetBlanked":     `{"blanked": false}`,
			"/ITB-1101-SW1.av/GetVideoInput": `{"1": "1"}`,
			"/ITB-1101-SW1.av/GetAudioInput": `{"2": "1"}`,
			"/ITB-1101-AMP1.av/GetVolume":    `{"volume": 30}`,
			"/ITB-1101-AMP1.av/GetMuted":     `{"muted": false}`,
			"/ITB-1101-VIA1.av/GetVolume":    `{"volume": 50}`,
			"/ITB-1101-VIA1.av/GetMuted":     `{"muted": true}`,
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-HDMI1"),
						Video:            deviceID("ITB-1101-VIA1"),
						CanSetSeparately: boolP(true),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
				},
				"ITB-1101-AMP1": {
					Volume: intP(30),
					Muted:  boolP(false),
				},
				"ITB-1101-VIA1": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
			},
		},
	},
	{
		name:        "JustAddPower/1",
		dataService: &mock.JustAddPowerRoom{},
		env:         "default",
		driverResps: map[string]string{
			"/ITB-1101-D1.av/GetPower":            `{"poweredOn": true}`,
			"/ITB-1101-D1.av/GetAVInput":          `{"input": "hdmi!2"}`,
			"/ITB-1101-D1.av/GetBlanked":          `{"blanked": false}`,
			"/ITB-1101-D1.av/GetMuted":            `{"muted": true}`,
			"/ITB-1101-D1.av/GetVolume":           `{"volume": 100}`,
			"/ITB-1101-D2.av/GetPower":            `{"poweredOn": true}`,
			"/ITB-1101-D2.av/GetAVInput":          `{"input": "hdmi!2"}`,
			"/ITB-1101-D2.av/GetBlanked":          `{"blanked": false}`,
			"/ITB-1101-D2.av/GetMuted":            `{"muted": true}`,
			"/ITB-1101-D2.av/GetVolume":           `{"volume": 100}`,
			"/ITB-1101-D3.av/GetPower":            `{"poweredOn": true}`,
			"/ITB-1101-D3.av/GetAVInput":          `{"input": "hdmi!2"}`,
			"/ITB-1101-D3.av/GetBlanked":          `{"blanked": false}`,
			"/ITB-1101-D3.av/GetMuted":            `{"muted": true}`,
			"/ITB-1101-D3.av/GetVolume":           `{"volume": 100}`,
			"/ITB-1101-RX1.av/GetStream":          `{"input": "10.66.76.185"}`,
			"/ITB-1101-RX2.av/GetStream":          `{"input": "10.66.76.188"}`,
			"/ITB-1101-RX3.av/GetStream":          `{"input": "10.66.76.187"}`,
			"/ITB-1101-VIA1.av/GetVolume":         `{"volume": 50}`,
			"/ITB-1101-VIA1.av/GetMuted":          `{"muted": true}`,
			"/ITB-1101-DSP1.av/Mic1/volume/level": `{"volume": 50}`,
			"/ITB-1101-DSP1.av/Mic1/mute/status":  `{"muted": true}`,
			"/ITB-1101-DSP1.av/Mic2/volume/level": `{"volume": 50}`,
			"/ITB-1101-DSP1.av/Mic2/mute/status":  `{"muted": true}`,
			"/ITB-1101-DSP1.av/Mic3/volume/level": `{"volume": 50}`,
			"/ITB-1101-DSP1.av/Mic3/mute/status":  `{"muted": true}`,
		},
		apiResp: api.StateResponse{
			Devices: map[api.DeviceID]api.DeviceState{
				"ITB-1101-D1": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-HDMI1"),
						Video:            deviceID("ITB-1101-HDMI1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-PC1"),
							api.DeviceID("ITB-1101-SIGN1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(100),
					Muted:  boolP(true),
				},
				"ITB-1101-D2": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-SIGN1"),
						Video:            deviceID("ITB-1101-SIGN1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-PC1"),
							api.DeviceID("ITB-1101-SIGN1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(100),
					Muted:  boolP(true),
				},
				"ITB-1101-D3": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: &api.Input{
						Audio:            deviceID("ITB-1101-PC1"),
						Video:            deviceID("ITB-1101-PC1"),
						CanSetSeparately: boolP(false),
						AvailableInputs: []api.DeviceID{
							api.DeviceID("ITB-1101-HDMI1"),
							api.DeviceID("ITB-1101-PC1"),
							api.DeviceID("ITB-1101-SIGN1"),
							api.DeviceID("ITB-1101-VIA1"),
						},
					},
					Volume: intP(100),
					Muted:  boolP(true),
				},
				"ITB-1101-VIA1": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
				"ITB-1101-MIC1": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
				"ITB-1101-MIC2": {
					Volume: intP(50),
					Muted:  boolP(true),
				},
				"ITB-1101-MIC3": {
					Volume: intP(50),
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
				fmt.Fprintln(w, tt.driverResps[r.URL.Path])
			}))
			t.Cleanup(func() {
				ts.Close()
			})

			tt.dataService.SetBaseURL(ts.URL)

			room, err := tt.dataService.Room(ctx, tt.room)
			if err != nil {
				t.Errorf("unable to get room: %s", err)
			}

			gs := &GetSetter{
				Environment: tt.env,
				Logger:      log.Logger{},
			}

			resp, err := gs.Get(ctx, room)
			if err != nil {
				t.Errorf("unable to get room state: %s", err)
			}

			if diff := cmp.Diff(tt.apiResp, resp); diff != "" {
				t.Errorf("generated incorrect response (-want, +got):\n%s", diff)
			}
		})
	}
}
