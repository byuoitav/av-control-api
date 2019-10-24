package drivers

import "context"

type DSP interface {
	Device

	GetVolumeByBlock(ctx context.Context, addr string, block string) (int, error)
	GetMutedByBlock(ctx context.Context, addr string, block string) (bool, error)

	SetVolumeByBlock(ctx context.Context, addr string, block string, volume int) error
	SetMutedByBlock(ctx context.Context, addr string, block string, muted bool) error
}
