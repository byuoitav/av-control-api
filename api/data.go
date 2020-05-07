package api

type Device struct {
	// fill out struct info
}

type DeviceType struct {
}

type DeviceService interface {
	Room(id string) ([]Device, error)
	Device(id string) (Device, error)
}
