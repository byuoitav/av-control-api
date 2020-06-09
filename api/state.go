package api

import "context"

type StateGetSetter interface {
	Get(ctx context.Context, room Room) (StateResponse, error)
	Set(ctx context.Context, room Room, req StateRequest) (StateResponse, error)
}
