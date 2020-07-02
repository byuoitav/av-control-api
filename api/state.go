package api

import "context"

type StateGetSetter interface {
	Get(context.Context, Room) (StateResponse, error)
	Set(context.Context, Room, StateRequest) (StateResponse, error)

	DriverStates(context.Context) (map[string]string, error)
}
