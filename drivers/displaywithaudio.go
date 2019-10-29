package drivers

import "context"

type simpleAudioDevice interface {
	GetVolume(ctx context.Context) (int, error)
	GetMuted(ctx context.Context) (bool, error)
}

type SimpleAudioDisplay interface {
	Device
	display
	simpleAudioDevice
}
