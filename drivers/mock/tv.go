package mock

import (
	"context"
	"sync"
	"time"
)

type TV struct {
	Delay time.Duration

	On       bool
	AVInputs map[string]string
	Blank    bool
	Volume   int
	Mute     bool

	sync.Mutex
}

func (m *TV) GetPower(context.Context) (bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return m.On, nil
}

func (m *TV) SetPower(ctx context.Context, pow bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.On = pow
	return nil
}

func (m *TV) GetAudioVideoInputs(context.Context) (map[string]string, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return m.AVInputs, nil
}

func (m *TV) SetAudioVideoInput(ctx context.Context, output, input string) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.AVInputs[output] = input
	return nil
}

func (m *TV) GetBlank(context.Context) (bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return m.Blank, nil
}

func (m *TV) SetBlank(ctx context.Context, blank bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Blank = blank
	return nil
}

func (m *TV) GetVolumes(context.Context, []string) (map[string]int, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return map[string]int{
		"": m.Volume,
	}, nil
}

func (m *TV) SetVolume(ctx context.Context, block string, level int) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Volume = level
	return nil
}

func (m *TV) GetMutes(context.Context, []string) (map[string]bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return map[string]bool{
		"": m.Mute,
	}, nil
}

func (m *TV) SetMute(ctx context.Context, block string, mute bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Mute = mute
	return nil
}
