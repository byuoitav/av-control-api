package mock

import (
	"context"
)

type WithPower struct {
	PoweredOn bool
	Error     error
	SetError  error
}

func (d WithPower) Power(ctx context.Context) (bool, error) {
	return d.PoweredOn, d.Error
}

func (d WithPower) SetPower(ctx context.Context, poweredOn bool) error {
	if d.SetError == nil {
		d.PoweredOn = poweredOn
	}

	return d.Error
}

type WithAudioInput struct {
	Inputs   map[string]string
	Error    error
	SetError error
}

func (d WithAudioInput) AudioInputs(ctx context.Context) (map[string]string, error) {
	return d.Inputs, d.Error
}

func (d WithAudioInput) SetAudioInput(ctx context.Context, output, input string) error {
	if d.SetError == nil {
		d.Inputs[output] = input
	}

	return d.SetError
}

type WithVideoInput struct {
	Inputs   map[string]string
	Error    error
	SetError error
}

func (d WithVideoInput) VideoInputs(ctx context.Context) (map[string]string, error) {
	return d.Inputs, d.Error
}

func (d WithVideoInput) SetVideoInput(ctx context.Context, output, input string) error {
	if d.SetError == nil {
		d.Inputs[output] = input
	}

	return d.SetError
}

type WithAudioVideoInput struct {
	Inputs   map[string]string
	Error    error
	SetError error
}

func (d WithAudioVideoInput) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	return d.Inputs, d.Error
}

func (d WithAudioVideoInput) SetAudioVideoInput(ctx context.Context, output, input string) error {
	if d.SetError == nil {
		d.Inputs[output] = input
	}

	return d.SetError
}

type WithBlank struct {
	Blanked  bool
	Error    error
	SetError error
}

func (d WithBlank) Blank(ctx context.Context) (bool, error) {
	return d.Blanked, d.Error
}

func (d WithBlank) SetBlank(ctx context.Context, blanked bool) error {
	if d.SetError == nil {
		d.Blanked = blanked
	}

	return d.Error
}

type WithVolume struct {
	Vols     map[string]int
	Error    error
	SetError error
}

func (d WithVolume) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	vols := make(map[string]int)
	for _, block := range blocks {
		if vol, ok := d.Vols[block]; ok {
			vols[block] = vol
		}
	}

	return vols, d.Error
}

func (d WithVolume) SetVolume(ctx context.Context, block string, vol int) error {
	if d.SetError == nil {
		d.Vols[block] = vol
	}

	return d.SetError
}

type WithMute struct {
	Ms       map[string]bool
	Error    error
	SetError error
}

func (d WithMute) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	ms := make(map[string]bool)
	for _, block := range blocks {
		if muted, ok := d.Ms[block]; ok {
			ms[block] = muted
		}
	}

	return ms, d.Error
}

func (d WithMute) SetMute(ctx context.Context, block string, muted bool) error {
	if d.SetError == nil {
		d.Ms[block] = muted
	}

	return d.SetError
}

type WithHealth struct {
	Error error
}

func (d WithHealth) Healthy(ctx context.Context) error {
	return d.Error
}

type WithInfo struct {
	I     interface{}
	Error error
}

func (d WithInfo) Info(ctx context.Context) (interface{}, error) {
	return d.I, d.Error
}
