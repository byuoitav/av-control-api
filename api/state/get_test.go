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
	deviceService interface {
		api.DeviceService
		SetBaseURL(string)
	}
	env string

	httpResps map[string]string
	apiResp   api.StateResponse
}

var getTests = []getStateTest{
	getStateTest{
		name:          "simple",
		deviceService: &mock.SimpleRoom{},
		env:           "default",
		httpResps:     map[string]string{},
		apiResp:       api.StateResponse{},
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
				fmt.Printf("closing server: %s\n", ts.URL)
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

			if !cmp.Equal(resp, tt.apiResp) {
				t.Errorf("generated incorrect response:\n\tgot %+v\n\texpected: %+v", resp, tt.apiResp)
			}
		})
	}
}
