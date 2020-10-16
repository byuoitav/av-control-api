package avcontrol

import "context"

// TODO all of these interfaces should delete the `Get` prefix
type (
	Device interface{}

	DeviceWithPower interface {
		GetPower(context.Context) (bool, error)
		SetPower(context.Context, bool) error
	}

	DeviceWithAudioInput interface {
		GetAudioInputs(ctx context.Context) (map[string]string, error)
		SetAudioInput(ctx context.Context, output, input string) error
	}

	DeviceWithVideoInput interface {
		GetVideoInputs(ctx context.Context) (map[string]string, error)
		SetVideoInput(ctx context.Context, output, input string) error
	}

	DeviceWithAudioVideoInput interface {
		GetAudioVideoInputs(ctx context.Context) (map[string]string, error)
		SetAudioVideoInput(ctx context.Context, output, input string) error
	}

	DeviceWithBlank interface {
		GetBlank(context.Context) (bool, error)
		SetBlank(context.Context, bool) error
	}

	DeviceWithVolume interface {
		GetVolumes(ctx context.Context, blocks []string) (map[string]int, error)
		SetVolume(context.Context, string, int) error
	}

	DeviceWithMute interface {
		GetMutes(ctx context.Context, blocks []string) (map[string]bool, error)
		SetMute(context.Context, string, bool) error
	}

	DeviceWithHealth interface {
		// Healthy returns a nil error if the device is healthy.
		Healthy(context.Context) error
	}
)
