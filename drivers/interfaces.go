package drivers

import "context"

type DeviceWithAudioInput interface {
	GetAudioInputs(ctx context.Context) (map[string]string, error)
	SetAudioInput(ctx context.Context, output, input string) error
}

type DeviceWithVideoInput interface {
	GetVideoInputs(ctx context.Context) (map[string]string, error)
	SetVideoInput(ctx context.Context, output, input string) error
}

type DeviceWithAudioVideoInput interface {
	GetAudioVideoInputs(ctx context.Context) (map[string]string, error)
	SetAudioVideoInput(ctx context.Context, output, input string) error
}

type DeviceWithPower interface {
	GetPower(context.Context) (bool, error)
	SetPower(context.Context, bool) error
}

type DeviceWithBlank interface {
	GetBlank(context.Context) (bool, error)
	SetBlank(context.Context, bool) error
}

type DeviceWithVolume interface {
	GetVolumes(context.Context) (map[string]int, error)
	SetVolume(context.Context, string, int) error
}

type DeviceWithMute interface {
	GetMutes(context.Context) (map[string]bool, error)
	SetMute(context.Context, string, bool) error
}

type DeviceWithInfo interface {
	GetInfo(context.Context) (interface{}, error)
}
