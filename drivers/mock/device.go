package mock

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNotSupported  = errors.New("not supported")
	ErrInvalidOutput = errors.New("this output does not exist")
	ErrInvalidBlock  = errors.New("this block does not exist")
)

type Device struct {
	Delay time.Duration

	On               *bool
	AudioInputs      map[string]string
	VideoInputs      map[string]string
	AudioVideoInputs map[string]string
	Blanked          *bool
	Volumes          map[string]int
	Mutes            map[string]bool

	once              sync.Once
	capabilities      []string
	audioOutputs      []string
	videoOutputs      []string
	audioVideoOutputs []string
	volumeBlocks      []string
	muteBlocks        []string

	sync.Mutex
}

func (d *Device) init() {
	d.Lock()
	defer d.Unlock()

	if d.On != nil {
		d.capabilities = append(d.capabilities, "Power")
	}

	if len(d.AudioInputs) > 0 {
		d.capabilities = append(d.capabilities, "AudioInput")
		for k := range d.AudioInputs {
			d.audioOutputs = append(d.audioOutputs, k)
		}
	}

	if len(d.VideoInputs) > 0 {
		d.capabilities = append(d.capabilities, "VideoInput")
		for k := range d.VideoInputs {
			d.videoOutputs = append(d.videoOutputs, k)
		}
	}

	if len(d.AudioVideoInputs) > 0 {
		d.capabilities = append(d.capabilities, "AudioVideoInput")
		for k := range d.AudioVideoInputs {
			d.audioVideoOutputs = append(d.audioVideoOutputs, k)
		}
	}

	if d.Blanked != nil {
		d.capabilities = append(d.capabilities, "Blank")
	}

	if len(d.Volumes) > 0 {
		d.capabilities = append(d.capabilities, "Volume")
		for k := range d.Volumes {
			d.volumeBlocks = append(d.volumeBlocks, k)
		}
	}

	if len(d.Mutes) > 0 {
		d.capabilities = append(d.capabilities, "Mute")
		for k := range d.Mutes {
			d.muteBlocks = append(d.muteBlocks, k)
		}
	}
}

func (d *Device) hasCapability(c string) bool {
	d.once.Do(d.init)

	for i := range d.capabilities {
		if d.capabilities[i] == c {
			return true
		}
	}

	return false
}

func contains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}

	return false
}

func (d *Device) GetCapabilities(context.Context) ([]string, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()
	return d.capabilities, nil
}

func (d *Device) GetPower(context.Context) (bool, error) {
	d.once.Do(d.init)

	if !d.hasCapability("Power") {
		return false, ErrNotSupported
	}

	d.Lock()
	defer d.Unlock()
	time.Sleep(d.Delay)
	return *d.On, nil
}

func (d *Device) SetPower(ctx context.Context, pow bool) error {
	d.once.Do(d.init)

	if !d.hasCapability("Power") {
		return ErrNotSupported
	}

	d.Lock()
	defer d.Unlock()
	time.Sleep(d.Delay)
	d.On = &pow
	return nil
}

func (d *Device) GetAudioInputs(context.Context) (map[string]string, error) {
	d.once.Do(d.init)

	if !d.hasCapability("AudioInput") {
		return nil, ErrNotSupported
	}

	d.Lock()
	defer d.Unlock()
	time.Sleep(d.Delay)
	return d.AudioInputs, nil
}

func (d *Device) SetAudioInput(ctx context.Context, output, input string) error {
	d.once.Do(d.init)

	if !d.hasCapability("AudioInput") {
		return ErrNotSupported
	}

	if !contains(d.audioOutputs, output) {
		return ErrInvalidOutput
	}

	d.Lock()
	defer d.Unlock()
	time.Sleep(d.Delay)
	d.AudioInputs[output] = input
	return nil
}

/*
func (d *Device) GetBlank(context.Context) (bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return m.Blank, nil
}

func (d *Device) SetBlank(ctx context.Context, blank bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Blank = blank
	return nil
}

func (d *Device) GetVolumes(context.Context, []string) (map[string]int, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return map[string]int{
		"": m.Volume,
	}, nil
}

func (d *Device) SetVolume(ctx context.Context, block string, level int) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Volume = level
	return nil
}

func (d *Device) GetMutes(context.Context, []string) (map[string]bool, error) {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	return map[string]bool{
		"": m.Mute,
	}, nil
}

func (d *Device) SetMute(ctx context.Context, block string, mute bool) error {
	m.Lock()
	defer m.Unlock()
	time.Sleep(m.Delay)
	m.Mute = mute
	return nil
}
*/
