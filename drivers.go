package avcontrol

import "context"

type DriverRegistry interface {
	Register(string, Driver)
	MustRegister(string, Driver)
	Get(string) Driver
	List() []string
}

type Driver interface {
	CreateDevice(context.Context, string) (Device, error)
}
