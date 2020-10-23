package mock

type TV struct {
	WithPower
	WithAudioVideoInput
	WithBlank
	WithVolume
	WithMute
}

type TVSeparateInput struct {
	WithPower
	WithAudioInput
	WithVideoInput
	WithBlank
	WithVolume
	WithMute
}

type Projector struct {
	WithPower
	WithAudioVideoInput
	WithBlank
}

type BasicVideoSwitcher struct {
	WithAudioVideoInput
}

type VideoSwitcher struct {
	WithAudioInput
	WithVideoInput
}

type VideoSwitcherDSP struct {
	WithAudioInput
	WithVideoInput
	WithVolume
	WithMute
}

type DSP struct {
	WithVolume
	WithMute
}

type AVOverIPReceiver struct {
	WithAudioVideoInput
}
