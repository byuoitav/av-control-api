package couch

import (
	"fmt"

	"github.com/byuoitav/av-control-api/api"
	kivik "github.com/go-kivik/kivik/v4"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// Room gets a room
func (d *DataService) Room(ctx context.Context, id string) (api.Room, error) {
	var devs []device
	var types []deviceType

	group, gctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error

		devs, err = d.Devices(gctx, id)
		if err != nil {
			return fmt.Errorf("unable to get devices: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		var err error

		types, err = d.DeviceTypes(gctx)
		if err != nil {
			return fmt.Errorf("unable to get device types: %w", err)
		}

		return err
	})

	if err := group.Wait(); err != nil {
		return api.Room{}, err
	}

	room := api.Room{
		ID: id,
	}

	// fill the devices with the types
	for i := range devs {
		for j := range types {
			if devs[i].TypeID == types[j].ID {
				devs[i].Type = types[j]
				break
			}
		}

		if devs[i].Type.ID == "" {
			return room, fmt.Errorf("no device type %q found", devs[i].TypeID)
		}

		apiDev, err := devs[i].convert()
		if err != nil {
			return room, fmt.Errorf("unable to convert %q: %w", devs[i].ID, err)
		}

		room.Devices = append(room.Devices, apiDev)
	}

	return room, nil
}

func (d *DataService) Devices(ctx context.Context, roomID string) ([]device, error) {
	var devs []device

	db := d.client.DB(ctx, "devices")
	roomQuery := kivik.Options{
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{
				"$regex": roomID + "-.*",
			},
		},
	}

	rows, err := db.Find(ctx, roomQuery)
	if err != nil {
		return devs, fmt.Errorf("unable to get all docs: %w", err)
	}

	for rows.Next() {
		var dev device
		if err = rows.ScanDoc(&dev); err != nil {
			return devs, fmt.Errorf("unable to scan %q: %w", rows.ID(), err)
		}

		devs = append(devs, dev)
	}

	return devs, nil
}

func (d *DataService) DriverMapping(ctx context.Context, roomID string) (api.DriverMapping, error) {
	var mapping driverMapping

	db := d.client.DB(ctx, "devices")

	err := db.Get(context.TODO(), roomID).ScanDoc(&mapping)
	if err != nil {
		return api.DriverMapping{}, fmt.Errorf("unable to scan in driver mapping for %s: %s", roomID, err)
	}

	toReturn := mapping.convert()

	return toReturn, nil
}

// DeviceTypes gets all of the device type documents available
// TODO i'm sure we can optimize this to only get the type documents associated with the devices in this room
// (in a single request)
func (d *DataService) DeviceTypes(ctx context.Context) ([]deviceType, error) {
	var types []deviceType

	db := d.client.DB(ctx, "device-types")

	opts := kivik.Options{
		"include_docs": true,
	}

	rows, err := db.AllDocs(ctx, opts)
	if err != nil {
		return types, fmt.Errorf("unable to get all docs: %w", err)
	}

	for rows.Next() {
		var typ deviceType
		if err = rows.ScanDoc(&typ); err != nil {
			return types, fmt.Errorf("unable to scan %q: %w", rows.ID(), err)
		}

		types = append(types, typ)
	}

	return types, nil
}
