package couch

import (
	"context"
	"errors"
	"net/url"
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/go-kivik/kivikmock/v3"
	"github.com/matryer/is"
)

func TestRoomConfig(t *testing.T) {
	is := is.New(t)

	client, mock, err := kivikmock.New()
	is.NoErr(err)

	ds, err := NewWithClient(context.Background(), client)
	is.NoErr(err)

	db := mock.NewDB()
	mock.ExpectDB().WithName(ds.database).WillReturn(db)
	db.ExpectGet().WithDocID("ITB-1101").WillReturn(kivikmock.DocumentT(t, `{
		"_id": "ITB-1101",
		"proxy": "http://ITB-1101-CP1.byu.edu:17000",
		"devices": {
			"ITB-1101-D1": {
				"driver": "sony/bravia",
				"address": "ITB-1101-D1.byu.edu"
			},
			"ITB-1101-D2": {
				"driver": "sony/adcp",
				"address": "ITB-1101-D2.byu.edu"
			},
			"ITB-1101-DSP1": {
				"driver": "qsc",
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

	room, err := ds.RoomConfig(context.Background(), "ITB-1101")
	is.NoErr(err)

	expectedURL, err := url.Parse("http://ITB-1101-CP1.byu.edu:17000")
	is.NoErr(err)

	is.Equal(room, avcontrol.RoomConfig{
		ID:    "ITB-1101",
		Proxy: expectedURL,
		Devices: map[avcontrol.DeviceID]avcontrol.DeviceConfig{
			"ITB-1101-D1": {
				Driver:  "sony/bravia",
				Address: "ITB-1101-D1.byu.edu",
			},
			"ITB-1101-D2": {
				Driver:  "sony/adcp",
				Address: "ITB-1101-D2.byu.edu",
			},
			"ITB-1101-DSP1": {
				Driver:  "qsc",
				Address: "ITB-1101-DSP1.byu.edu",
				Ports: avcontrol.PortConfigs{
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
	})
}

func TestBadRoom(t *testing.T) {
	is := is.New(t)

	client, mock, err := kivikmock.New()
	is.NoErr(err)

	ds, err := NewWithClient(context.Background(), client)
	is.NoErr(err)

	db := mock.NewDB()
	mock.ExpectDB().WithName(ds.database).WillReturn(db)
	db.ExpectGet().WithDocID("ITB-1101").WillReturnError(errors.New("doc doesn't exist"))

	_, err = ds.RoomConfig(context.Background(), "ITB-1101")
	is.Equal(err.Error(), "unable to get/scan room: doc doesn't exist")
}

func TestRoomConvertFail(t *testing.T) {
	is := is.New(t)

	r := room{
		ID:    "ITB-1101",
		Proxy: ":foo",
	}

	_, err := r.convert()
	is.Equal(err.Error(), `unable to parse proxy url: parse ":foo": missing protocol scheme`)
}
