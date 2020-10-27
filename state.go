package avcontrol

import "context"

// StateGetSetter represents something that can get and set device state.
type StateGetSetter interface {
	// Get gets the state for the given Room.
	Get(context.Context, RoomConfig) (StateResponse, error)

	// Set sets the state for the given Room, based on the data in the StateRequest.
	Set(context.Context, RoomConfig, StateRequest) (StateResponse, error)

	// GetHealth returns the health status of the room
	GetHealth(context.Context, RoomConfig) (RoomHealth, error)

	// GetInfo returns the info about each device in the room
	GetInfo(context.Context, RoomConfig) (RoomInfo, error)
}
