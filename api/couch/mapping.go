package couch

import (
	"fmt"

	"github.com/byuoitav/av-control-api/api"
	"golang.org/x/net/context"
)

type driverMapping struct {
	Drivers map[string]struct {
		BaseURLs map[string]string `json:"baseURLs"`
	} `json:"drivers"`
}

func (d *DataService) DriverMapping(ctx context.Context) (api.DriverMapping, error) {
	var mapping driverMapping

	db := d.client.DB(ctx, d.database)
	if err := db.Get(ctx, d.mappingDocID).ScanDoc(&mapping); err != nil {
		return api.DriverMapping{}, fmt.Errorf("unable to get/scan driver mapping: %w", err)
	}

	return mapping.convert(d.environment), nil
}

func (d driverMapping) convert(env string) api.DriverMapping {
	mapping := make(api.DriverMapping)

	for k, v := range d.Drivers {
		if baseURL, ok := v.BaseURLs[env]; ok {
			mapping[k] = baseURL
		}
	}

	return mapping
}
