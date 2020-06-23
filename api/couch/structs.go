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

type driverMapping struct {
	Mapping map[string]struct {
		BaseURLs map[string]string `json:"BaseURLs"`
	} `json:"mapping"`
}

func (d driverMapping) convert() api.DriverMapping {
	toReturn := make(map[string]struct {
		BaseURLs map[string]string `json:"BaseURLs"`
	})
	for k, v := range d.Mapping {
		toReturn[k] = v
	}

	return toReturn
}

func (d device) convert() (api.Device, error) {
	toReturn := api.Device{
		ID:      api.DeviceID(d.ID),
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

func (p port) convert() api.Port {
	return api.Port{
		Name: p.Name,
		Type: p.Type,
	}
}
