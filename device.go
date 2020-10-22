package avcontrol

import "context"

type (
	Device interface{}

	DeviceWithPower interface {
		Power(context.Context) (bool, error)
		SetPower(context.Context, bool) error
	}

	DeviceWithAudioInput interface {
		AudioInputs(ctx context.Context) (map[string]string, error)
		SetAudioInput(ctx context.Context, output, input string) error
	}

	DeviceWithVideoInput interface {
		VideoInputs(ctx context.Context) (map[string]string, error)
		SetVideoInput(ctx context.Context, output, input string) error
	}

	DeviceWithAudioVideoInput interface {
		AudioVideoInputs(ctx context.Context) (map[string]string, error)
		SetAudioVideoInput(ctx context.Context, output, input string) error
	}

	DeviceWithBlank interface {
		Blank(context.Context) (bool, error)
		SetBlank(context.Context, bool) error
	}

	DeviceWithVolume interface {
		Volumes(ctx context.Context, blocks []string) (map[string]int, error)
		SetVolume(context.Context, string, int) error
	}

	DeviceWithMute interface {
		Mutes(ctx context.Context, blocks []string) (map[string]bool, error)
		SetMute(context.Context, string, bool) error
	}

	DeviceWithHealth interface {
		// Healthy returns a nil error if the device is healthy.
		Healthy(context.Context) error
	}

	DeviceWithInfo interface {
		// Info returns a struct of info about the device.
		Info(context.Context) (interface{}, error)
	}

	/* TODO gotta figure out what we want from active signal - video, audio...?
	DeviceWithActiveSignal interface {
		// ActiveVideoSignal returns true if the current
		ActiveVideoSignal(context.Context) (map[string]bool, error)
	}
	*/
)
