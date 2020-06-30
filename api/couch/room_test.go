package couch

import (
	"context"
	"net/url"
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"github.com/go-kivik/kivikmock/v3"
	"github.com/google/go-cmp/cmp"
)

func TestRoom(t *testing.T) {
	client, mock, err := kivikmock.New()
	if err != nil {
		t.Fatalf("unable to create kivik mock: %s", err)
	}

	ds := &DataService{
		client:       client,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
		environment:  "default",
	}

	db := mock.NewDB()
	mock.ExpectDB().WithName(ds.database).WillReturn(db)
	db.ExpectGet().WithDocID("ITB-1101").WillReturn(kivikmock.DocumentT(t, `{
		"_id": "ITB-1101",
		"proxy": "http://ITB-1101-CP1.byu.edu:17000",
		"devices": {
			"ITB-1101-D1": {
				"driver": "Sony Bravia",
				"address": "ITB-1101-D1.byu.edu"
			},
			"ITB-1101-D2": {
				"driver": "Sony ADCP",
				"address": "ITB-1101-D2.byu.edu"
			},
			"ITB-1101-DSP1": {
				"driver": "QSC",
				"address": "ITB-1101-DSP1.byu.edu",
				"ports": [
					{
						"name": "Mic1Gain",
						"type": "volume"
					},
					{
						"name": "Mic1Mute",
						"type": "mute"
					}
				]
			}
		}
	}`))

	room, err := ds.Room(context.Background(), "ITB-1101")
	if err != nil {
		t.Fatalf("unable to get mapping: %s", err)
	}

	expectedURL, _ := url.Parse("http://ITB-1101-CP1.byu.edu:17000")

	expected := api.Room{
		ID:    "ITB-1101",
		Proxy: expectedURL,
		Devices: map[api.DeviceID]api.Device{
			"ITB-1101-D1": api.Device{
				Driver:  "Sony Bravia",
				Address: "ITB-1101-D1.byu.edu",
			},
			"ITB-1101-D2": api.Device{
				Driver:  "Sony ADCP",
				Address: "ITB-1101-D2.byu.edu",
			},
			"ITB-1101-DSP1": api.Device{
				Driver:  "QSC",
				Address: "ITB-1101-DSP1.byu.edu",
				Ports: api.Ports{
					{
						Name: "Mic1Gain",
						Type: "volume",
					},
					{
						Name: "Mic1Mute",
						Type: "mute",
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, room); diff != "" {
		t.Errorf("generated incorrect mapping (-want, +got):\n%s", diff)
	}
}
