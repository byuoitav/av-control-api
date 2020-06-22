package drivers

import (
	"context"
	sync "sync"
	"time"
)

type mockTV struct {
	delay time.Duration
	sync.Mutex

	on       bool
	avInputs map[string]string
	blank    bool
	volume   int
	mute     bool
}

func (m *mockTV) GetPower(context.Context) (bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	return m.on, nil
}

func (m *mockTV) SetPower(ctx context.Context, pow bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	m.on = pow
	return nil
}

func (m *mockTV) GetAudioVideoInputs(context.Context) (map[string]string, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	return m.avInputs, nil
}

func (m *mockTV) SetAudioVideoInput(ctx context.Context, output, input string) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	m.avInputs[output] = input
	return nil
}

func (m *mockTV) GetBlank(context.Context) (bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	return m.blank, nil
}

func (m *mockTV) SetBlank(ctx context.Context, blank bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	m.blank = blank
	return nil
}

func (m *mockTV) GetVolumes(context.Context, []string) (map[string]int, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	return map[string]int{
		"": m.volume,
	}, nil
}

func (m *mockTV) SetVolume(ctx context.Context, block string, level int) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	m.volume = level
	return nil
}

func (m *mockTV) GetMutes(context.Context, []string) (map[string]bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	return map[string]bool{
		"": m.mute,
	}, nil
}

func (m *mockTV) SetMute(ctx context.Context, block string, mute bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.delay)
	m.mute = mute
	return nil
}
