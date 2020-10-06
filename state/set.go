package state

import (
	"context"
	"errors"
	"fmt"
	"sync"

	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

var (
	ErrInvalidDevice = errors.New("device is invalid in this room")
	ErrNotCapable    = errors.New("can't set this field on this device")
	ErrInvalidBlock  = errors.New("invalid block")
)

type setDeviceStateRequest struct {
	id     avcontrol.DeviceID
	device avcontrol.DeviceConfig
	state  avcontrol.DeviceState
	driver avcontrol.Driver
	log    *zap.Logger
}

type setDeviceStateResponse struct {
	id     avcontrol.DeviceID
	state  avcontrol.DeviceState
	errors []avcontrol.DeviceStateError
	sync.Mutex
}

func (gs *GetSetter) Set(ctx context.Context, room avcontrol.RoomConfig, req avcontrol.StateRequest) (avcontrol.StateResponse, error) {
	if len(room.Devices) == 0 {
		return avcontrol.StateResponse{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		if gs.DriverRegistry.Get(dev.Driver) == nil {
			return avcontrol.StateResponse{}, fmt.Errorf("%s: %w", dev.Driver, ErrDriverNotRegistered)
		}
	}

	// make sure each of the devices in the request are in this room
	for id := range req.Devices {
		if _, ok := room.Devices[id]; !ok {
			return avcontrol.StateResponse{}, fmt.Errorf("%s: %w", id, ErrInvalidDevice)
		}
	}

	id := avcontrol.CtxRequestID(ctx)
	log := gs.Logger
	if len(id) > 0 {
		log = gs.Logger.With(zap.String("requestID", id))
	}

	expectedResps := 0
	resps := make(chan setDeviceStateResponse)
	defer close(resps)

	stateResp := avcontrol.StateResponse{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceState),
	}

	for id, dev := range room.Devices {
		state, ok := req.Devices[id]
		if !ok {
			continue
		}

		expectedResps++
		req := setDeviceStateRequest{
			id:     id,
			device: dev,
			state:  state,
			driver: gs.DriverRegistry.Get(dev.Driver),
			log:    log.With(zap.String("deviceID", string(id))),
		}

		go func() {
			resps <- req.do(ctx)
		}()
	}

	for resp := range resps {
		expectedResps--

		stateResp.Devices[resp.id] = resp.state
		stateResp.Errors = append(stateResp.Errors, resp.errors...)

		if expectedResps == 0 {
			break
		}
	}

	sortErrors(stateResp.Errors)
	return stateResp, nil
}

func (req *setDeviceStateRequest) do(ctx context.Context) setDeviceStateResponse {
	resp := setDeviceStateResponse{
		id: req.id,
	}

	req.log.Info("Setting state")
	req.log.Debug("Getting device")

	dev, err := req.driver.CreateDevice(ctx, req.device.Address)
	if err != nil {
		req.log.Warn("unable to get device", zap.Error(err))
		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Error: fmt.Sprintf("unable to get device: %s", err),
		})
		return resp
	}

	req.log.Debug("Got device")

	handleErr := func(field string, value interface{}, err error) {
		req.log.Warn("unable to set "+field, zap.Any("to", value), zap.Error(err))

		resp.Lock()
		defer resp.Unlock()

		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Field: field,
			Value: value,
			Error: err.Error(),
		})
	}

	// set every field in req.state
	wg := sync.WaitGroup{}

	// setting power happens before everything else is run
	if req.state.PoweredOn != nil {
		if dev, ok := dev.(avcontrol.DeviceWithPower); ok {
			req.log.Info("Setting power", zap.Bool("poweredOn", *req.state.PoweredOn))

			if err := dev.SetPower(ctx, *req.state.PoweredOn); err != nil {
				handleErr("poweredOn", *req.state.PoweredOn, err)
			} else {
				req.log.Info("Set power")
				resp.state.PoweredOn = req.state.PoweredOn
			}
		} else {
			handleErr("poweredOn", *req.state.PoweredOn, ErrNotCapable)
		}
	}

	// figure out which inputs we set
	var setAudioInput, setVideoInput, setAudioVideoInput bool
	for _, input := range req.state.Inputs {
		if input.Audio != nil {
			setAudioInput = true
		}

		if input.Video != nil {
			setVideoInput = true
		}

		if input.AudioVideo != nil {
			setAudioVideoInput = true
		}
	}

	if setAudioInput {
		if dev, ok := dev.(avcontrol.DeviceWithAudioInput); ok {
			for output, input := range req.state.Inputs {
				if input.Audio == nil {
					continue
				}

				wg.Add(1)
				go func(output, input string) {
					req.log.Info("Setting audio input", zap.String("output", output), zap.String("input", input))
					defer wg.Done()

					if err := dev.SetAudioInput(ctx, output, input); err != nil {
						handleErr(fmt.Sprintf("input.%s.audio", output), input, err)
						return
					}

					req.log.Info("Set audio input", zap.String("output", output))

					resp.Lock()
					defer resp.Unlock()

					if resp.state.Inputs == nil {
						resp.state.Inputs = make(map[string]avcontrol.Input)
					}

					in := resp.state.Inputs[output]
					in.Audio = &input
					resp.state.Inputs[output] = in
				}(output, *input.Audio)
			}
		} else {
			handleErr("input.$.audio", req.state.Inputs, ErrNotCapable)
		}
	}

	if setVideoInput {
		if dev, ok := dev.(avcontrol.DeviceWithVideoInput); ok {
			for output, input := range req.state.Inputs {
				if input.Video == nil {
					continue
				}

				wg.Add(1)
				go func(output, input string) {
					req.log.Info("Setting video input", zap.String("output", output), zap.String("input", input))
					defer wg.Done()

					if err := dev.SetVideoInput(ctx, output, input); err != nil {
						handleErr(fmt.Sprintf("input.%s.video", output), input, err)
						return
					}

					req.log.Info("Set video input", zap.String("output", output))

					resp.Lock()
					defer resp.Unlock()

					if resp.state.Inputs == nil {
						resp.state.Inputs = make(map[string]avcontrol.Input)
					}

					in := resp.state.Inputs[output]
					in.Video = &input
					resp.state.Inputs[output] = in
				}(output, *input.Video)
			}
		} else {
			handleErr("input.$.video", req.state.Inputs, ErrNotCapable)
		}
	}

	if setAudioVideoInput {
		if dev, ok := dev.(avcontrol.DeviceWithAudioVideoInput); ok {
			for output, input := range req.state.Inputs {
				if input.AudioVideo == nil {
					continue
				}

				wg.Add(1)
				go func(output, input string) {
					req.log.Info("Setting audioVideo input", zap.String("output", output), zap.String("input", input))
					defer wg.Done()

					if err := dev.SetAudioVideoInput(ctx, output, input); err != nil {
						handleErr(fmt.Sprintf("input.%s.audioVideo", output), input, err)
						return
					}

					req.log.Info("Set audioVideo input", zap.String("output", output))

					resp.Lock()
					defer resp.Unlock()

					if resp.state.Inputs == nil {
						resp.state.Inputs = make(map[string]avcontrol.Input)
					}

					in := resp.state.Inputs[output]
					in.AudioVideo = &input
					resp.state.Inputs[output] = in
				}(output, *input.AudioVideo)
			}
		} else {
			handleErr("input.$.audioVideo", req.state.Inputs, ErrNotCapable)
		}
	}

	if req.state.Blanked != nil {
		if dev, ok := dev.(avcontrol.DeviceWithBlank); ok {
			wg.Add(1)

			go func() {
				req.log.Info("Setting blank", zap.Bool("blanked", *req.state.Blanked))
				defer wg.Done()

				if err := dev.SetBlank(ctx, *req.state.Blanked); err != nil {
					handleErr("blanked", *req.state.Blanked, err)
					return
				}

				req.log.Info("Set blank")

				resp.Lock()
				defer resp.Unlock()
				resp.state.Blanked = req.state.Blanked
			}()
		} else {
			handleErr("blanked", *req.state.Blanked, ErrNotCapable)
		}
	}

	if len(req.state.Volumes) > 0 {
		if dev, ok := dev.(avcontrol.DeviceWithVolume); ok {
			validBlocks := req.device.Ports.OfType("volume").Names()
			for block, vol := range req.state.Volumes {
				if !containsString(validBlocks, block) {
					handleErr(fmt.Sprintf("volumes.%s", block), vol, ErrInvalidBlock)
					continue
				}

				wg.Add(1)
				go func(block string, vol int) {
					req.log.Info("Setting volume", zap.String("block", block), zap.Int("level", vol))
					defer wg.Done()

					if err := dev.SetVolume(ctx, block, vol); err != nil {
						handleErr(fmt.Sprintf("volumes.%s", block), vol, err)
						return
					}

					req.log.Info("Set volume", zap.String("block", block))

					resp.Lock()
					defer resp.Unlock()

					if resp.state.Volumes == nil {
						resp.state.Volumes = make(map[string]int)
					}

					resp.state.Volumes[block] = vol
				}(block, vol)
			}
		} else {
			handleErr("volumes", req.state.Volumes, ErrNotCapable)
		}
	}

	if len(req.state.Mutes) > 0 {
		if dev, ok := dev.(avcontrol.DeviceWithMute); ok {
			validBlocks := req.device.Ports.OfType("mute").Names()
			for block, muted := range req.state.Mutes {
				if !containsString(validBlocks, block) {
					handleErr(fmt.Sprintf("mutes.%s", block), muted, ErrInvalidBlock)
					continue
				}

				wg.Add(1)
				go func(block string, muted bool) {
					req.log.Info("Setting mute", zap.String("block", block), zap.Bool("muted", muted))
					defer wg.Done()

					if err := dev.SetMute(ctx, block, muted); err != nil {
						handleErr(fmt.Sprintf("mutes.%s", block), muted, err)
						return
					}

					req.log.Info("Set mute", zap.String("block", block))

					resp.Lock()
					defer resp.Unlock()

					if resp.state.Mutes == nil {
						resp.state.Mutes = make(map[string]bool)
					}

					resp.state.Mutes[block] = muted
				}(block, muted)
			}
		} else {
			handleErr("mutes", req.state.Mutes, ErrNotCapable)
		}
	}

	wg.Wait()

	req.log.Info("Finished setting state")
	return resp
}
