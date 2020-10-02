package couch

import (
	"context"
	"errors"
	"net/url"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/go-kivik/kivikmock/v3"
	"github.com/google/go-cmp/cmp"
)

// TODO use testify - require

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

	expected := avcontrol.Room{
		ID:    "ITB-1101",
		Proxy: expectedURL,
		Devices: map[avcontrol.DeviceID]avcontrol.Device{
			"ITB-1101-D1": avcontrol.Device{
				Driver:  "Sony Bravia",
				Address: "ITB-1101-D1.byu.edu",
			},
			"ITB-1101-D2": avcontrol.Device{
				Driver:  "Sony ADCP",
				Address: "ITB-1101-D2.byu.edu",
			},
			"ITB-1101-DSP1": avcontrol.Device{
				Driver:  "QSC",
				Address: "ITB-1101-DSP1.byu.edu",
				Ports: avcontrol.Ports{
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

func TestBadRoom(t *testing.T) {
	errWanted := errors.New("unable to get/scan room: doc doesn't exist")
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
	db.ExpectGet().WithDocID("ITB-1101").WillReturnError(errors.New("doc doesn't exist"))

	_, err = ds.Room(context.Background(), "ITB-1101")
	if err == nil {
		t.Fatalf("did not get error back")
	}

	if diff := cmp.Diff(errWanted.Error(), err.Error()); diff != "" {
		t.Errorf("generated incorrect error (-want, +got):\n%s", diff)
	}
}

func TestRoomConvertFail(t *testing.T) {
	errWanted := errors.New(`unable to parse proxy url: parse ":foo": missing protocol scheme`)
	r := room{
		ID:    "ITB-1101",
		Proxy: ":foo",
	}

	_, err := r.convert()
	if err == nil {
		t.Fatalf("no error converting room")
	}

	if diff := cmp.Diff(errWanted.Error(), err.Error()); diff != "" {
		t.Errorf("generated incorrect error (-want, +got):\n%s", diff)
	}
}
