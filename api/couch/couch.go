package couch

import "github.com/byuoitav/av-control-api/api"

type DataService struct {
}

func (d *DataService) Room(id string) ([]api.Device, error) {
	return []api.Device{}, nil
}
