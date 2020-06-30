package couch

import (
	"fmt"
	"net/url"

	"github.com/byuoitav/av-control-api/api"
	"golang.org/x/net/context"
)

type room struct {
	ID      string            `json:"_id"`
	Proxy   string            `json:"proxy"`
	Devices map[string]device `json:"devices"`
}

type device struct {
	Address string `json:"address"`
	Driver  string `json:"driver"`
	Ports   []port `json:"ports"`
}

type port struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Room gets a room
func (d *DataService) Room(ctx context.Context, id string) (api.Room, error) {
	var room room

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, id).ScanDoc(&room); err != nil {
		return api.Room{}, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return room.convert()
}

func (r room) convert() (api.Room, error) {
	url, err := url.Parse(r.Proxy)
	if err != nil {
		return api.Room{}, fmt.Errorf("unable to parse proxy url: %w", err)
	}

	room := api.Room{
		ID:      r.ID,
		Proxy:   url,
		Devices: make(map[api.DeviceID]api.Device),
	}

	for id, dev := range r.Devices {
		room.Devices[api.DeviceID(id)] = dev.convert()
	}

	return room, nil
}

func (d device) convert() api.Device {
	dev := api.Device{
		Address: d.Address,
		Driver:  d.Driver,
	}

	for i := range d.Ports {
		dev.Ports = append(dev.Ports, d.Ports[i].convert())
	}

	return dev
}

func (p port) convert() api.Port {
	return api.Port{
		Name: p.Name,
		Type: p.Type,
	}
}
