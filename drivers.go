package avcontrol

import (
	"context"
)

type DriverRegistry interface {
	Register(string, Driver) error
	MustRegister(string, Driver)
	Get(string) Driver
	List() []string
}

type Driver interface {
	CreateDevice(context.Context, string) (Device, error)
	ParseConfig(map[string]interface{}) error
}
