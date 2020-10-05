package state

import (
	"context"
	"errors"
	"fmt"
	"sync"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/drivers"
	"go.uber.org/zap"
)

var (
	ErrNoStateGettable     = errors.New("can't get the state of any devices in this room")
	ErrDriverNotRegistered = errors.New("driver not registered")
)

type getDeviceStateRequest struct {
	id     avcontrol.DeviceID
	device avcontrol.Device
	driver *avcontrol.Driver
	log    *zap.Logger
}

type getDeviceStateResponse struct {
	id     avcontrol.DeviceID
	state  avcontrol.DeviceState
	errors []avcontrol.DeviceStateError
	sync.Mutex
}

func (gs *GetSetter) Get(ctx context.Context, room avcontrol.Room) (avcontrol.StateResponse, error) {
	if len(room.Devices) == 0 {
		return avcontrol.StateResponse{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		if gs.Drivers.Get(dev.Driver) == nil {
			return avcontrol.StateResponse{}, fmt.Errorf("%s: %w", dev.Driver, ErrDriverNotRegistered)
		}
	}

	id := avcontrol.CtxRequestID(ctx)
	log := gs.Logger
	if len(id) > 0 {
		log = gs.Logger.With(zap.String("requestID", id))
	}

	resps := make(chan getDeviceStateResponse)
	defer close(resps)

	stateResp := avcontrol.StateResponse{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceState),
	}

	for id, dev := range room.Devices {
		req := getDeviceStateRequest{
			id:     id,
			device: dev,
			driver: gs.Drivers.Get(dev.Driver),
			log:    log.With(zap.String("deviceID", string(id))),
		}

		go func() {
			resps <- req.do(ctx)
		}()
	}

	received := 0
	for resp := range resps {
		received++

		stateResp.Devices[resp.id] = resp.state
		stateResp.Errors = append(stateResp.Errors, resp.errors...)

		if received == len(room.Devices) {
			break
		}
	}

	sortErrors(stateResp.Errors)
	return stateResp, nil
}

func (req *getDeviceStateRequest) do(ctx context.Context) getDeviceStateResponse {
	resp := getDeviceStateResponse{
		id: req.id,
	}

	req.log.Info("Getting state")
	req.log.Debug("Getting device")

	dev, err := req.driver.GetDevice(ctx, req.device.Address)
	if err != nil {
		req.log.Warn("unable to get device", zap.Error(err))
		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Error: fmt.Sprintf("unable to get device: %s", err),
		})
		return resp
	}

	req.log.Debug("Got device")

	handleErr := func(field string, err error) {
		req.log.Warn("unable to get "+field, zap.Error(err))

		resp.Lock()
		defer resp.Unlock()

		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Field: field,
			Error: err.Error(),
		})
	}

	// get every field possible
	wg := sync.WaitGroup{}

	if dev, ok := dev.(drivers.DeviceWithPower); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting power")
			defer wg.Done()

			power, err := dev.GetPower(ctx)
			if err != nil {
				handleErr("power", err)
				return
			}

			req.log.Info("Got power", zap.Bool("poweredOn", power))

			resp.Lock()
			defer resp.Unlock()
			resp.state.PoweredOn = &power
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithAudioInput); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting audio inputs")
			defer wg.Done()

			inputs, err := dev.GetAudioInputs(ctx)
			if err != nil {
				handleErr("inputs.$.audio", err)
				return
			}

			req.log.Info("Got audio inputs", zap.Any("inputs", inputs))

			resp.Lock()
			defer resp.Unlock()

			if resp.state.Inputs == nil {
				resp.state.Inputs = make(map[string]avcontrol.Input)
			}

			for out, in := range inputs {
				input := resp.state.Inputs[out]
				save := in
				input.Audio = &save
				resp.state.Inputs[out] = input
			}
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithVideoInput); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting video inputs")
			defer wg.Done()

			inputs, err := dev.GetVideoInputs(ctx)
			if err != nil {
				handleErr("inputs.$.video", err)
				return
			}

			req.log.Info("Got video inputs", zap.Any("inputs", inputs))

			resp.Lock()
			defer resp.Unlock()

			if resp.state.Inputs == nil {
				resp.state.Inputs = make(map[string]avcontrol.Input)
			}

			for out, in := range inputs {
				input := resp.state.Inputs[out]
				save := in
				input.Video = &save
				resp.state.Inputs[out] = input
			}
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithAudioVideoInput); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting audioVideo inputs")
			defer wg.Done()

			inputs, err := dev.GetAudioVideoInputs(ctx)
			if err != nil {
				handleErr("inputs.$.audioVideo", err)
				return
			}

			req.log.Info("Got audioVideo inputs", zap.Any("inputs", inputs))

			resp.Lock()
			defer resp.Unlock()

			if resp.state.Inputs == nil {
				resp.state.Inputs = make(map[string]avcontrol.Input)
			}

			for out, in := range inputs {
				input := resp.state.Inputs[out]
				save := in
				input.AudioVideo = &save
				resp.state.Inputs[out] = input
			}
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithBlank); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting blank")
			defer wg.Done()

			blank, err := dev.GetBlank(ctx)
			if err != nil {
				handleErr("blank", err)
				return
			}

			req.log.Info("Got blank", zap.Bool("blank", blank))

			resp.Lock()
			defer resp.Unlock()
			resp.state.Blanked = &blank
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithVolume); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting volumes")
			defer wg.Done()

			vols, err := dev.GetVolumes(ctx, req.device.Ports.OfType("volume").Names())
			if err != nil {
				handleErr("volumes", err)
				return
			}

			req.log.Info("Got volumes", zap.Any("volumes", vols))

			resp.Lock()
			defer resp.Unlock()
			resp.state.Volumes = vols
		}()
	}

	if dev, ok := dev.(drivers.DeviceWithMute); ok {
		wg.Add(1)

		go func() {
			req.log.Info("Getting mutes")
			defer wg.Done()

			mutes, err := dev.GetMutes(ctx, req.device.Ports.OfType("mute").Names())
			if err != nil {
				handleErr("mutes", err)
				return
			}

			req.log.Info("Got mutes", zap.Any("mutes", mutes))

			resp.Lock()
			defer resp.Unlock()
			resp.state.Mutes = mutes
		}()
	}

	wg.Wait()

	req.log.Info("Finished getting state")
	return resp
}
