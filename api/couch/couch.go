package couch

import (
	"errors"
	"fmt"

	_ "github.com/go-kivik/couchdb/v4"
	kivik "github.com/go-kivik/kivik/v4"
	"golang.org/x/net/context"
)

type DataService struct {
	DBAddress  string
	DBUsername string
	DBPassword string
}

// Room gets a room
func (d *DataService) Room(ctx context.Context, id string) ([]device, error) {
	addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, d.DBAddress)
	// fmt.Printf("address: %s\n", addr)
	client, err := kivik.New("couch", addr)
	if err != nil {
		return []device{}, fmt.Errorf("unable to connect to couch: %s", err)
	}

	db := client.DB(context.TODO(), "devices")

	roomQuery := kivik.Options{
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{
				"$regex": id + "-.*",
			},
		},
	}

	devices, err := db.Find(context.TODO(), roomQuery)
	if err != nil {
		return []device{}, fmt.Errorf("unable to find devices in room %s: %s", id, err)
	}

	var toReturn []device
	added := false

	for devices.Next() {
		if devices.EOQ() {
			break
		}

		var dev device
		if err = devices.ScanDoc(&dev); err != nil {
			fmt.Printf("error scanning in device doc\n")
			continue
		}

		toReturn = append(toReturn, dev)
		added = true
	}

	if added {
		return toReturn, nil
	}
	return []device{}, errors.New("unable to get room")
}

// Device gets a device
func (d *DataService) Device(ctx context.Context, id string) (device, error) {
	addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, d.DBAddress)
	client, err := kivik.New("couch", addr)
	if err != nil {
		return device{}, fmt.Errorf("unable to connect to couch: %s", err)
	}

	db := client.DB(ctx, "devices")

	var dev device
	if err = db.Get(ctx, id).ScanDoc(&dev); err != nil {
		return dev, fmt.Errorf("error retrieving device doc: %s", err)
	}

	dt, err := d.DeviceType(ctx, dev.TypeID)
	if err != nil {
		return dev, fmt.Errorf("error retrieving device type doc: %s", err)
	}

	dev.Type = dt

	return dev, nil
}

func (d *DataService) DeviceType(ctx context.Context, id string) (deviceType, error) {
	addr := fmt.Sprintf("https://%s:%s@%s", d.DBUsername, d.DBPassword, d.DBAddress)
	client, err := kivik.New("couch", addr)
	if err != nil {
		return deviceType{}, fmt.Errorf("unable to connect to couch: %s", err)
	}

	db := client.DB(ctx, "device-types")

	var dt deviceType
	if err = db.Get(ctx, id).ScanDoc(&dt); err != nil {
		return dt, fmt.Errorf("error retrieving device type doc: %s", err)
	}

	return dt, nil
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
