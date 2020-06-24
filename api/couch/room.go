package couch

import (
	"fmt"

	"github.com/byuoitav/av-control-api/api"
	"golang.org/x/net/context"
)

type room struct {
	ID           string            `json:"_id"`
	ProxyBaseURL string            `json:"proxyBaseURL"`
	Devices      map[string]device `json:"devices"`
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
	room := api.Room{
		ID:           r.ID,
		ProxyBaseURL: r.ProxyBaseURL,
		Devices:      make(map[api.DeviceID]api.Device),
	}

	for id, dev := range r.Devices {
		apiDev, err := dev.convert()
		if err != nil {
			return room, fmt.Errorf("unable to convert device %q: %w", id, err)
		}

		room.Devices[api.DeviceID(id)] = apiDev
	}

	return room, nil
}

func (d device) convert() (api.Device, error) {
	dev := api.Device{
		Address: d.Address,
		Driver:  d.Driver,
	}

	for i := range d.Ports {
		dev.Ports = append(dev.Ports, d.Ports[i].convert())
	}

	return dev, nil
}

func (p port) convert() api.Port {
	return api.Port{
		Name: p.Name,
		Type: p.Type,
	}
}
