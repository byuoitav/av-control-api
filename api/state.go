package api

import "context"

type StateGetSetter interface {
	Get(ctx context.Context, room []Device) (StateResponse, error)
	Set(ctx context.Context, room []Device, req StateRequest) (StateResponse, error)
}
