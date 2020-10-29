package mock

type TV struct {
	WithPower
	WithAudioVideoInput
	WithBlank
	WithVolume
	WithMute
	WithHealth
	WithInfo
}

type TVSeparateInput struct {
	WithPower
	WithAudioInput
	WithVideoInput
	WithBlank
	WithVolume
	WithMute
	WithHealth
	WithInfo
}

type Projector struct {
	WithPower
	WithAudioVideoInput
	WithBlank
	WithHealth
	WithInfo
}

type BasicVideoSwitcher struct {
	WithAudioVideoInput
	WithHealth
	WithInfo
}

type VideoSwitcher struct {
	WithAudioInput
	WithVideoInput
	WithHealth
	WithInfo
}

type VideoSwitcherDSP struct {
	WithAudioInput
	WithVideoInput
	WithVolume
	WithMute
	WithHealth
	WithInfo
}

type DSP struct {
	WithVolume
	WithMute
	WithHealth
	WithInfo
}

type AVOverIPReceiver struct {
	WithAudioVideoInput
	WithHealth
	WithInfo
}
