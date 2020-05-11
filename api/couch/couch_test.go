package couch

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

// func TestRoom(t *testing.T) {
// 	d := DataService{
// 		DBAddress:  strings.Trim(os.Getenv("DB_ADDRESS"), "https://"),
// 		DBUsername: os.Getenv("DB_USERNAME"),
// 		DBPassword: os.Getenv("DB_PASSWORD"),
// 	}

// 	devices, err := d.Room(context.TODO(), "ITB-1108A")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	for i := range devices {
// 		fmt.Printf("ID: %s\n", devices[i].ID)
// 		fmt.Printf("Address: %s\n", devices[i].Address)
// 		fmt.Printf("Type: %s\n", devices[i].Type.ID)
// 	}

// 	devices, err = d.Room(context.TODO(), "ITB-1108B")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	for i := range devices {
// 		fmt.Printf("ID: %s\n", devices[i].ID)
// 		fmt.Printf("Address: %s\n", devices[i].Address)
// 		fmt.Printf("Type: %s\n", devices[i].Type.ID)
// 	}
// }

func TestDevice(t *testing.T) {
	d := DataService{
		DBAddress:  strings.Trim(os.Getenv("DB_ADDRESS"), "https://"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}

	device, err := d.Device(context.TODO(), "ITB-1108A-CP1")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("ID: %s\n", device.ID)
	fmt.Printf("Address: %s\n", device.Address)
	fmt.Printf("Type: %s\n", device.Type.ID)

	device, err = d.Device(context.TODO(), "ITB-1108A-CP2")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("ID: %s\n", device.ID)
	fmt.Printf("Address: %s\n", device.Address)
	fmt.Printf("Type: %s\n", device.Type.ID)

	device, err = d.Device(context.TODO(), "ITB-1108A-CP3")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("ID: %s\n", device.ID)
	fmt.Printf("Address: %s\n", device.Address)
	fmt.Printf("Type: %s\n", device.Type.ID)
}
