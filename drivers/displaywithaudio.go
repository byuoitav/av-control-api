package drivers

import "context"

type DisplayWithAudio interface {
	Display

	GetVolume(ctx context.Context, addr string) (int, error)
	GetMuted(ctx context.Context, addr string) (bool, error)
}
