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

// RoomConfig gets a room
func (d *DataService) RoomConfig(ctx context.Context, id string) (avcontrol.RoomConfig, error) {
	var room room

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, id).ScanDoc(&room); err != nil {
		return avcontrol.RoomConfig{}, fmt.Errorf("unable to get/scan room: %w", err)
	}

	return room.convert()
}

func (r room) convert() (avcontrol.RoomConfig, error) {
	url, err := url.Parse(r.Proxy)
	if err != nil {
		return avcontrol.RoomConfig{}, fmt.Errorf("unable to parse proxy url: %w", err)
	}

	room := avcontrol.RoomConfig{
		ID:      r.ID,
		Proxy:   url,
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceConfig),
	}

	for id, dev := range r.Devices {
		room.Devices[avcontrol.DeviceID(id)] = dev.convert()
	}

	return room, nil
}

func (d device) convert() avcontrol.DeviceConfig {
	dev := avcontrol.DeviceConfig{
		Address: d.Address,
		Driver:  d.Driver,
	}

	for i := range d.Ports {
		dev.Ports = append(dev.Ports, d.Ports[i].convert())
	}

	return dev
}

func (p port) convert() avcontrol.PortConfig {
	return avcontrol.PortConfig{
		Name: p.Name,
		Type: p.Type,
	}
}
