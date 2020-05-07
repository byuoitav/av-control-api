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
}

func (dt deviceType) convert() api.DeviceType {
}

func (c command) convert() api.Command {
}

func (p port) convert() api.Port {
}
