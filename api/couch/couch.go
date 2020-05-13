package couch

import (
	"errors"
	"fmt"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	_ "github.com/go-kivik/couchdb/v4"
	kivik "github.com/go-kivik/kivik/v4"
	"golang.org/x/net/context"
)

type DataService struct {
	DBAddress     string
	DBUsername    string
	DBPassword    string
	devicesDB     *kivik.DB
	deviceTypesDB *kivik.DB
}

// Room gets a room as an array of devices
func (d *DataService) Room(ctx context.Context, id string) ([]api.Device, error) {
	if d.devicesDB == nil {
		a := strings.Trim(d.DBAddress, "https://")
		addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, a)
		// fmt.Printf("address: %s\n", addr)
		client, err := kivik.New("couch", addr)
		if err != nil {
			return []api.Device{}, fmt.Errorf("unable to connect to couch: %s", err)
		}

		d.devicesDB = client.DB(context.TODO(), "devices")
	}

	roomQuery := kivik.Options{
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{
				"$regex": id + "-.*",
			},
		},
		"limit": 100,
	}

	devices, err := d.devicesDB.Find(ctx, roomQuery)
	if err != nil {
		return []api.Device{}, fmt.Errorf("unable to find devices in room %s: %s", id, err)
	}

	dt, err := d.AllDeviceTypes(ctx)
	if err != nil {
		fmt.Println("error retrieving device types: %s", err)
	}

	var toReturn []api.Device
	added := false

	for devices.Next() {
		var dev device
		if err = devices.ScanDoc(&dev); err != nil {
			fmt.Printf("error scanning in device doc\n")
			continue
		}

		for i := range dt {
			if dev.TypeID == dt[i].ID {
				dev.Type = dt[i]
				break
			}
		}

		add, err := dev.convert()
		if err != nil {
			fmt.Println("error converting doc into api.Device: %s", err)
		}
		toReturn = append(toReturn, add)
		added = true
	}

	if added {
		// fmt.Printf("toReturn: %+v\n", toReturn)
		return toReturn, nil
	}

	return []api.Device{}, errors.New("unable to get room")
}

// Device gets a device
func (d *DataService) Device(ctx context.Context, id string) (api.Device, error) {
	if d.devicesDB == nil {
		a := strings.Trim(d.DBAddress, "https://")
		addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, a)
		// fmt.Printf("address: %s\n", addr)

		client, err := kivik.New("couch", addr)
		if err != nil {
			return api.Device{}, fmt.Errorf("unable to connect to couch: %s", err)
		}

		d.devicesDB = client.DB(ctx, "devices")
	}

	var dev device
	if err := d.devicesDB.Get(ctx, id).ScanDoc(&dev); err != nil {
		return api.Device{}, fmt.Errorf("error retrieving device doc: %s", err)
	}

	dt, err := d.DeviceType(ctx, dev.TypeID)
	if err != nil {
		return api.Device{}, fmt.Errorf("error retrieving device type doc for %s: %s", id, err)
	}

	dev.Type = dt

	toReturn, err := dev.convert()
	if err != nil {
		return api.Device{}, fmt.Errorf("unable to convert doc into api.Device: %w", err)
	}

	return toReturn, nil
}

// DeviceType is for a single device query because we only need the one device type doc
func (d *DataService) DeviceType(ctx context.Context, id string) (deviceType, error) {
	if d.deviceTypesDB == nil {
		a := strings.Trim(d.DBAddress, "https://")
		addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, a)
		// fmt.Printf("address: %s\n", addr)

		client, err := kivik.New("couch", addr)
		if err != nil {
			return deviceType{}, fmt.Errorf("unable to connect to couch: %s", err)
		}

		d.deviceTypesDB = client.DB(ctx, "device-types")
	}

	var dt deviceType
	if err := d.deviceTypesDB.Get(ctx, id).ScanDoc(&dt); err != nil {
		return deviceType{}, fmt.Errorf("error retrieving device type doc: %s", err)
	}

	return dt, nil
}

// AllDeviceTypes is for an array of devices, we just get all the device type docs
func (d *DataService) AllDeviceTypes(ctx context.Context) ([]deviceType, error) {
	if d.deviceTypesDB == nil {
		a := strings.Trim(d.DBAddress, "https://")
		addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, a)
		// fmt.Printf("address: %s\n", addr)

		client, err := kivik.New("couch", addr)
		if err != nil {
			return []deviceType{}, fmt.Errorf("unable to connect to couch: %s", err)
		}

		d.deviceTypesDB = client.DB(ctx, "device-types")
	}

	query := kivik.Options{
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{
				"$regex": ".*",
			},
		},
		"limit": 100,
	}

	types, err := d.deviceTypesDB.Find(ctx, query)
	if err != nil {
		return []deviceType{}, fmt.Errorf("error finding all device type docs: %s", err)
	}

	var toReturn []deviceType

	for types.Next() {
		var dt deviceType
		if err = types.ScanDoc(&dt); err != nil {
			fmt.Printf("error scanning in device type for %s: %s\n", types.ID(), err)
			continue
		}

		toReturn = append(toReturn, dt)
	}

	return toReturn, nil
}

// IsHealthy is a healthcheck for the database
func (d *DataService) IsHealthy(ctx context.Context, dbName string) (bool, error) {
	addr := fmt.Sprintf("http://%s:%s@%s:5984", d.DBUsername, d.DBPassword, d.DBAddress)
	client, err := kivik.New("couch", addr)
	if err != nil {
		return false, fmt.Errorf("unable to connect to couch: %s", err)
	}

	alive, err := client.DBExists(ctx, dbName)
	if err != nil {
		return false, err
	}

	return alive, err
}
