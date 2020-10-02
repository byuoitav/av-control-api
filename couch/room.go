package couch

import (
	"fmt"
	"net/url"

	avcontrol "github.com/byuoitav/av-control-api"
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
func (d *DataService) Room(ctx context.Context, id string) (avcontrol.Room, error) {
	var room room

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, id).ScanDoc(&room); err != nil {
		return avcontrol.Room{}, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return room.convert()
}

func (r room) convert() (avcontrol.Room, error) {
	url, err := url.Parse(r.Proxy)
	if err != nil {
		return avcontrol.Room{}, fmt.Errorf("unable to parse proxy url: %w", err)
	}

	room := avcontrol.Room{
		ID:      r.ID,
		Proxy:   url,
		Devices: make(map[avcontrol.DeviceID]avcontrol.Device),
	}

	for id, dev := range r.Devices {
		room.Devices[avcontrol.DeviceID(id)] = dev.convert()
	}

	return room, nil
}

func (d device) convert() avcontrol.Device {
	dev := avcontrol.Device{
		Address: d.Address,
		Driver:  d.Driver,
	}

	for i := range d.Ports {
		dev.Ports = append(dev.Ports, d.Ports[i].convert())
	}

	return dev
}

func (p port) convert() avcontrol.Port {
	return avcontrol.Port{
		Name: p.Name,
		Type: p.Type,
	}
}
