package avcontrol

import "context"

// StateGetSetter represents something that can get and set device state.
type StateGetSetter interface {
	// Get gets the state for the given Room.
	Get(context.Context, Room) (StateResponse, error)

	// Set sets the state for the given Room, based on the data in the StateRequest.
	Set(context.Context, Room, StateRequest) (StateResponse, error)
}
