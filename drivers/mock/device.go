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

func (d *Device) VolumeBlocks() []string {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	b := make([]string, len(d.volumeBlocks))
	copy(b, d.volumeBlocks)
	return b
}

func (d *Device) MuteBlocks() []string {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	b := make([]string, len(d.muteBlocks))
	copy(b, d.muteBlocks)
	return b
}

func (d *Device) GetCapabilities(context.Context) ([]string, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()
	return d.capabilities, nil
}

func (d *Device) GetPower(context.Context) (bool, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Power") {
		return false, ErrNotSupported
	}

	time.Sleep(d.Delay)
	return *d.On, nil
}

func (d *Device) SetPower(ctx context.Context, pow bool) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Power") {
		return ErrNotSupported
	}

	time.Sleep(d.Delay)
	d.On = &pow
	return nil
}

func (d *Device) GetAudioInputs(context.Context) (map[string]string, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("AudioInput") {
		return nil, ErrNotSupported
	}

	time.Sleep(d.Delay)
	return d.AudioInputs, nil
}

func (d *Device) SetAudioInput(ctx context.Context, output, input string) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("AudioInput") {
		return ErrNotSupported
	}

	if !contains(d.audioOutputs, output) {
		return ErrInvalidOutput
	}

	time.Sleep(d.Delay)
	d.AudioInputs[output] = input
	return nil
}

func (d *Device) GetVideoInputs(context.Context) (map[string]string, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("VideoInput") {
		return nil, ErrNotSupported
	}

	time.Sleep(d.Delay)
	return d.VideoInputs, nil
}

func (d *Device) SetVideoInput(ctx context.Context, output, input string) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("VideoInput") {
		return ErrNotSupported
	}

	if !contains(d.videoOutputs, output) {
		return ErrInvalidOutput
	}

	time.Sleep(d.Delay)
	d.VideoInputs[output] = input
	return nil
}

func (d *Device) GetAudioVideoInputs(context.Context) (map[string]string, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("AudioVideoInput") {
		return nil, ErrNotSupported
	}

	time.Sleep(d.Delay)
	return d.AudioVideoInputs, nil
}

func (d *Device) SetAudioVideoInput(ctx context.Context, output, input string) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("AudioVideoInput") {
		return ErrNotSupported
	}

	if !contains(d.audioVideoOutputs, output) {
		return ErrInvalidOutput
	}

	time.Sleep(d.Delay)
	d.AudioVideoInputs[output] = input
	return nil
}

func (d *Device) GetBlank(context.Context) (bool, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Blank") {
		return false, ErrNotSupported
	}

	time.Sleep(d.Delay)
	return *d.Blanked, nil
}

func (d *Device) SetBlank(ctx context.Context, blanked bool) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Blank") {
		return ErrNotSupported
	}

	time.Sleep(d.Delay)
	d.Blanked = &blanked
	return nil
}

func (d *Device) GetVolumes(ctx context.Context, blocks []string) (map[string]int, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Volume") {
		return nil, ErrNotSupported
	}

	vols := make(map[string]int)
	for i := range blocks {
		if !contains(d.volumeBlocks, blocks[i]) {
			return nil, ErrInvalidBlock
		}

		vols[blocks[i]] = d.Volumes[blocks[i]]
	}

	time.Sleep(d.Delay)
	return vols, nil
}

func (d *Device) SetVolume(ctx context.Context, block string, level int) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Volume") {
		return ErrNotSupported
	}

	if !contains(d.volumeBlocks, block) {
		return ErrInvalidOutput
	}

	time.Sleep(d.Delay)
	d.Volumes[block] = level
	return nil
}

func (d *Device) GetMutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Mute") {
		return nil, ErrNotSupported
	}

	mutes := make(map[string]bool)
	for i := range blocks {
		if !contains(d.muteBlocks, blocks[i]) {
			return nil, ErrInvalidBlock
		}

		mutes[blocks[i]] = d.Mutes[blocks[i]]
	}

	time.Sleep(d.Delay)
	return mutes, nil
}

func (d *Device) SetMute(ctx context.Context, block string, muted bool) error {
	d.once.Do(d.init)

	d.Lock()
	defer d.Unlock()

	if !d.hasCapability("Mute") {
		return ErrNotSupported
	}

	if !contains(d.muteBlocks, block) {
		return ErrInvalidOutput
	}

	time.Sleep(d.Delay)
	d.Mutes[block] = muted
	return nil
}
