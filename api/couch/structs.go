package couch

import (
	"regexp"

	"github.com/byuoitav/av-control-api/api"
)

type device struct {
	ID      string            `json:"_id"`
	TypeID  string            `json:"typeID"`
	Address string            `json:"address"`
	Proxy   map[string]string `json:"proxy"`
	Ports   []port            `json:"ports"`
	Type    deviceType
}

type deviceType struct {
	ID       string             `json:"_id"`
	Commands map[string]command `json:"commands"`
}

type command struct {
	URLs  map[string]string `json:"addresses"`
	Order *int              `json:"order,omitempty"`
}

type port struct {
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	Incoming  bool     `json:"incoming"`
	Type      string   `json:"type"`
}

func (d device) convert() (api.Device, error) {
	toReturn := api.Device{
		ID:      api.DeviceID(d.ID),
		Type:    d.Type.convert(),
		Address: d.Address,
		Proxy:   make(map[*regexp.Regexp]string),
	}

	for i := range d.Ports {
		toReturn.Ports = append(toReturn.Ports, d.Ports[i].convert())
	}

	for k, v := range d.Proxy {
		exp, err := regexp.Compile(k)
		if err != nil {
			return api.Device{}, err
		}

		toReturn.Proxy[exp] = v
	}

	return toReturn, nil
}

func (dt deviceType) convert() api.DeviceType {
	toReturn := api.DeviceType{
		ID:       dt.ID,
		Commands: make(map[string]api.Command),
	}

	for key, val := range dt.Commands {
		toReturn.Commands[key] = val.convert()
	}

	return toReturn
}

func (c command) convert() api.Command {
	toReturn := api.Command{
		Order: c.Order,
		URLs:  make(map[string]string),
	}

	for key, val := range c.URLs {
		toReturn.URLs[key] = val
	}

	return toReturn
}

func (p port) convert() api.Port {
	var endpoints api.Endpoints
	for _, e := range p.Endpoints {
		endpoints = append(endpoints, api.DeviceID(e))
	}
	return api.Port{
		Name:      p.Name,
		Endpoints: endpoints,
		Incoming:  p.Incoming,
		Type:      p.Type,
	}
}
