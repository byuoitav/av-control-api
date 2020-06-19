package drivers

import (
	"context"
)

type mockTV struct {
	on       bool
	avInputs map[string]string
	blank    bool
	volume   int
	mute     bool
}

func (m *mockTV) GetPower(context.Context) (bool, error) {
	return m.on, nil
}

func (m *mockTV) SetPower(ctx context.Context, pow bool) error {
	m.on = pow
	return nil
}

func (m *mockTV) GetAudioVideoInputs(context.Context) (map[string]string, error) {
	return m.avInputs, nil
}

func (m *mockTV) SetAudioVideoInput(ctx context.Context, output, input string) error {
	m.avInputs[output] = input
	return nil
}

func (m *mockTV) GetBlank(context.Context) (bool, error) {
	return m.blank, nil
}

func (m *mockTV) SetBlank(ctx context.Context, blank bool) error {
	m.blank = blank
	return nil
}

func (m *mockTV) GetVolumes(context.Context, []string) (map[string]int, error) {
	return map[string]int{
		"": m.volume,
	}, nil
}

func (m *mockTV) SetVolume(ctx context.Context, block string, level int) error {
	m.volume = level
	return nil
}

func (m *mockTV) GetMutes(context.Context, []string) (map[string]bool, error) {
	return map[string]bool{
		"": m.mute,
	}, nil
}

func (m *mockTV) SetMute(ctx context.Context, block string, mute bool) error {
	m.mute = mute
	return nil
}
