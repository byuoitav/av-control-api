package adcp

import (
	"context"

	"github.com/byuoitav/adcp-driver"
)

func NewDevice(ctx context.Context, addr string) (interface{}, error) {
	return &adcp.Projector{
		Address: addr,
	}, nil
}
