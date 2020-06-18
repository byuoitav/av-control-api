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
	}

	driverResps map[string]string
	apiResp     api.StateResponse
}

var (
	via  = "ITB-1101-VIA1"
	hdmi = "ITB-1101-HDMI1"
	sign = "ITB-1101-SIGN1"
	pc   = "ITB-1101-PC1"
)

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
					Input: map[string]api.Input{
						"ITB-1101-D1": {
							Audio: &via,
							Video: &via,
						},
					},
					Volumes: map[string]int{
						"ITB-1101-D1": 30,
					},
					Mutes: map[string]bool{
						"ITB-1101-D1": false,
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{
						"ITB-1101-VIA1": 50,
					},
					Mutes: map[string]bool{
						"ITB-1101-VIA1": true,
					},
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
					Input: map[string]api.Input{
						"ITB-1101-D1": {
							Audio: &hdmi,
							Video: &hdmi,
						},
					},
					Volumes: map[string]int{
						"ITB-1101-D1": 100,
					},
					Mutes: map[string]bool{
						"ITB-1101-D1": true,
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{
						via: 50,
					},
					Mutes: map[string]bool{
						via: true,
					},
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
					Input: map[string]api.Input{
						"ITB-1101-D1": {
							Audio: &hdmi,
							Video: &via,
						},
					},
				},
				"ITB-1101-AMP1": {
					Volumes: map[string]int{
						"ITB-1101-AMP1": 30,
					},
					Mutes: map[string]bool{
						"ITB-1101-AMP1": false,
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{
						via: 50,
					},
					Mutes: map[string]bool{
						via: true,
					},
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
					Input: map[string]api.Input{
						"ITB-1101-D1": {
							Audio: &hdmi,
							Video: &hdmi,
						},
					},
					Volumes: map[string]int{
						"ITB-1101-D1": 100,
					},
					Mutes: map[string]bool{
						"ITB-1101-D1": true,
					},
				},
				"ITB-1101-D2": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: map[string]api.Input{
						"ITB-1101-D2": {
							Audio: &sign,
							Video: &sign,
						},
					},
					Volumes: map[string]int{
						"ITB-1101-D2": 100,
					},
					Mutes: map[string]bool{
						"ITB-1101-D2": true,
					},
				},
				"ITB-1101-D3": {
					PoweredOn: boolP(true),
					Blanked:   boolP(false),
					Input: map[string]api.Input{
						"ITB-1101-D3": {
							Audio: &pc,
							Video: &pc,
						},
					},
					Volumes: map[string]int{
						"ITB-1101-D3": 100,
					},
					Mutes: map[string]bool{
						"ITB-1101-D3": true,
					},
				},
				"ITB-1101-VIA1": {
					Volumes: map[string]int{
						via: 50,
					},
					Mutes: map[string]bool{
						via: true,
					},
				},
				"ITB-1101-MIC1": {
					Volumes: map[string]int{
						"ITB-1101-MIC1": 50,
					},
					Mutes: map[string]bool{
						"ITB-1101-MIC1": true,
					},
				},
				"ITB-1101-MIC2": {
					Volumes: map[string]int{
						"ITB-1101-MIC2": 50,
					},
					Mutes: map[string]bool{
						"ITB-1101-MIC2": true,
					},
				},
				"ITB-1101-MIC3": {
					Volumes: map[string]int{
						"ITB-1101-MIC3": 50,
					},
					Mutes: map[string]bool{
						"ITB-1101-MIC3": true,
					},
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
