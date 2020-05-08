package couch

import "github.com/byuoitav/av-control-api/api"

type device struct {
	ID      string            `json:"_id"`
	TypeID  string            `json:"typeID"`
	Address string            `json:"address"`
	Proxy   map[string]string `json:"proxy"`
	Ports   []port            `json:"ports"`
}

type deviceType struct {
	ID       string             `json:"_id"`
	Commands map[string]command `json:"commands"`
}

type command struct {
	URLs  map[string]string `json:"urls"`
	Order *int              `json:"order,omitempty"`
}

type port struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Incoming bool   `json:"incoming"`
	Outgoing bool   `json:"outgoing"`
	Type     string `json:"type"`
}

func (d device) convert() api.Device {
	toReturn := api.Device{
		ID:      d.ID,
		TypeID:  d.TypeID,
		Address: d.Address,
		Proxy:   d.Proxy,
	}

	for i := range d.Ports {
		toReturn.Ports = append(toReturn.Ports, d.Ports[i].convert())
	}

	return toReturn
}

func (dt deviceType) convert() api.DeviceType {
	toReturn := api.DeviceType{
		ID: dt.ID,
	}

	for key, val := range dt.Commands {
		toReturn.Commands[key] = val.convert()
	}

	return toReturn
}

func (c command) convert() api.Command {
	toReturn := api.Command{
		Order: c.Order,
	}

	for key, val := range c.URLs {
		toReturn.URLs[key] = val
	}

	return toReturn
}

func (p port) convert() api.Port {
	return api.Port{
		Name:     p.Name,
		Endpoint: api.DeviceID(p.Endpoint),
		Incoming: p.Incoming,
		Outgoing: p.Outgoing,
		Type:     p.Type,
	}
}
