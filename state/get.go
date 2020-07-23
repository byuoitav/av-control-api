package state

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/byuoitav/av-control-api/api"
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
	driver drivers.Driver
	log    *zap.Logger
}

type getDeviceStateResponse struct {
	id     avcontrol.DeviceID
	state  avcontrol.DeviceState
	errors []avcontrol.DeviceStateError
}

func (gs *getSetter) Get(ctx context.Context, room avcontrol.Room) (avcontrol.StateResponse, error) {
	if len(room.Devices) == 0 {
		return avcontrol.StateResponse{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		_, ok := drivers.Get(dev.Driver)
		if !ok {
			return avcontrol.StateResponse{}, fmt.Errorf("%s: %w", dev.Driver, ErrDriverNotRegistered)
		}
	}

	id := avcontrol.CtxRequestID(ctx)
	log := gs.logger
	if len(id) > 0 {
		log = gs.logger.With(zap.String("requestID", id))
	}

	stateResp := avcontrol.StateResponse{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceState),
	}

	resps := make(chan getDeviceStateResponse)

	for id, dev := range room.Devices {
		req := getDeviceStateRequest{
			id:     id,
			device: dev,
			driver: gs.drivers[dev.Driver],
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

	close(resps)
	return stateResp, nil
}

func (req *getDeviceStateRequest) do(ctx context.Context) (resp getDeviceStateResponse) {
	var respMu sync.Mutex
	resp = getDeviceStateResponse{
		id: req.id,
		state: avcontrol.DeviceState{
			Inputs:  make(map[string]avcontrol.Input),
			Volumes: make(map[string]int),
			Mutes:   make(map[string]bool),
		},
	}

	req.log.Info("Getting state")

	defer func() {
		// reset maps if they weren't used
		if len(resp.state.Inputs) == 0 {
			resp.state.Inputs = nil
		}

		if len(resp.state.Volumes) == 0 {
			resp.state.Volumes = nil
		}

		if len(resp.state.Mutes) == 0 {
			resp.state.Mutes = nil
		}
	}()

	driver, _ := drivers.Get(req.device.Driver)

	req.log.Debug("Getting device")

	dev, err := driver.GetDevice(ctx, req.device.Address)
	if err != nil {
		req.log.Warn("unable to get device", zap.Error(err))
		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Error: fmt.Sprintf("unable to get device: %s", err),
		})
		return resp
	}

	// TODO support getting capabilities a different way?

	driverErr := func(field string, err error) {
		req.log.Warn("unable to get "+field, zap.Error(err))

		respMu.Lock()
		defer respMu.Unlock()

		resp.errors = append(resp.errors, avcontrol.DeviceStateError{
			ID:    req.id,
			Field: field,
			Error: err.Error(),
		})
	}

	i := 0
	spreadRequests := func() {
		defer func() { i++ }()
		if i == 0 {
			return
		}

		time.Sleep(25 * time.Millisecond)
	}

	// req.log.Debug("Got capabilities", zap.Strings("capabilities", caps.GetCapabilities()))
	wg := sync.WaitGroup{}

	for _, capability := range caps.GetCapabilities() {
		switch drivers.Capability(capability) {
		case drivers.CapabilityPower:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting power")
				defer wg.Done()

				power, err := req.driver.GetPower(ctx, deviceInfo)
				if err != nil {
					driverErr("power", err)
					return
				}

				on := power.GetOn()
				req.log.Info("Got power", zap.Bool("on", on))

				respMu.Lock()
				defer respMu.Unlock()
				resp.state.PoweredOn = &on
			}()
		case drivers.CapabilityAudioInput:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting audio inputs")
				defer wg.Done()

				inputs, err := req.driver.GetAudioInputs(ctx, deviceInfo)
				if err != nil {
					driverErr("inputs.$.audio", err)
					return
				}

				req.log.Info("Got audio inputs", zap.Any("inputs", inputs.GetInputs()))

				respMu.Lock()
				defer respMu.Unlock()

				for out, in := range inputs.GetInputs() {
					input := resp.state.Inputs[out]
					save := in
					input.Audio = &save
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityVideoInput:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting video inputs")
				defer wg.Done()

				inputs, err := req.driver.GetVideoInputs(ctx, deviceInfo)
				if err != nil {
					driverErr("inputs.$.video", err)
					return
				}

				req.log.Info("Got video inputs", zap.Any("inputs", inputs.GetInputs()))

				respMu.Lock()
				defer respMu.Unlock()

				for out, in := range inputs.GetInputs() {
					input := resp.state.Inputs[out]
					save := in
					input.Video = &save
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityAudioVideoInput:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting audioVideo inputs")
				defer wg.Done()

				inputs, err := req.driver.GetAudioVideoInputs(ctx, deviceInfo)
				if err != nil {
					driverErr("inputs.$.audioVideo", err)
					return
				}

				req.log.Info("Got audioVideo inputs", zap.Any("inputs", inputs.GetInputs()))

				respMu.Lock()
				defer respMu.Unlock()

				for out, in := range inputs.GetInputs() {
					input := resp.state.Inputs[out]
					save := in
					input.AudioVideo = &save
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityBlank:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting blank")
				defer wg.Done()

				blank, err := req.driver.GetBlank(ctx, deviceInfo)
				if err != nil {
					driverErr("blank", err)
					return
				}

				blanked := blank.GetBlanked()
				req.log.Info("Got blank", zap.Bool("blanked", blanked))

				respMu.Lock()
				defer respMu.Unlock()
				resp.state.Blanked = &blanked
			}()
		case drivers.CapabilityVolume:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting volumes")
				defer wg.Done()

				audioInfo := &drivers.GetAudioInfo{
					Info:   deviceInfo,
					Blocks: req.device.Ports.OfType("volume").Names(),
				}

				vols, err := req.driver.GetVolumes(ctx, audioInfo)
				if err != nil {
					driverErr("volumes", err)
					return
				}

				req.log.Info("Got volumes", zap.Any("vols", vols.GetVolumes()))

				respMu.Lock()
				defer respMu.Unlock()

				for block, vol := range vols.GetVolumes() {
					resp.state.Volumes[block] = int(vol)
				}
			}()
		case drivers.CapabilityMute:
			wg.Add(1)
			spreadRequests()

			go func() {
				req.log.Info("Getting mutes")
				defer wg.Done()

				audioInfo := &drivers.GetAudioInfo{
					Info:   deviceInfo,
					Blocks: req.device.Ports.OfType("mute").Names(),
				}

				mutes, err := req.driver.GetMutes(ctx, audioInfo)
				if err != nil {
					driverErr("mutes", err)
					return
				}

				req.log.Info("Got mutes", zap.Any("mutes", mutes.GetMutes()))

				respMu.Lock()
				defer respMu.Unlock()

				for block, muted := range mutes.GetMutes() {
					resp.state.Mutes[block] = muted
				}
			}()
		case drivers.CapabilityInfo:
			// we don't do anything with info
		default:
			req.log.Warn("unknown capability", zap.String("capability", capability))

			respMu.Lock()
			resp.errors = append(resp.errors, api.DeviceStateError{
				ID:    req.id,
				Error: fmt.Sprintf("unknown capability %s", capability),
			})
			respMu.Unlock()

			continue
		}
	}

	wg.Wait()

	req.log.Info("Finished getting state")
	return
}