package state

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidDevice = errors.New("device is invalid in this room")
	ErrNotCapable    = errors.New("can't set this field on this device")
)

type setDeviceStateRequest struct {
	id     api.DeviceID
	device api.Device
	state  api.DeviceState
	driver drivers.DriverClient
	log    *zap.Logger
}

type setDeviceStateResponse struct {
	id     api.DeviceID
	state  api.DeviceState
	errors []api.DeviceStateError
}

func (gs *getSetter) Set(ctx context.Context, room api.Room, req api.StateRequest) (api.StateResponse, error) {
	if len(req.Devices) == 0 {
		return api.StateResponse{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		_, ok := gs.drivers[dev.Driver]
		if !ok {
			return api.StateResponse{}, fmt.Errorf("%w: %s", ErrUnknownDriver, dev.Driver)
		}
	}

	// make sure each of the devices in the request are in this room
	for id := range req.Devices {
		if _, ok := room.Devices[id]; !ok {
			return api.StateResponse{}, fmt.Errorf("%s: %w", id, ErrInvalidDevice)
		}
	}

	id := api.RequestID(ctx)
	log := gs.logger
	if len(id) > 0 {
		log = gs.logger.With(zap.String("requestID", id))
	}

	stateResp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
	}

	resps := make(chan setDeviceStateResponse)
	expectedResps := 0

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
			driver: gs.drivers[dev.Driver],
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

	close(resps)
	return stateResp, nil
}

func (req *setDeviceStateRequest) do(ctx context.Context) setDeviceStateResponse {
	var respMu sync.Mutex
	resp := setDeviceStateResponse{
		id: req.id,
		state: api.DeviceState{
			Inputs:  make(map[string]api.Input),
			Volumes: make(map[string]int),
			Mutes:   make(map[string]bool),
		},
	}

	deviceInfo := &drivers.DeviceInfo{
		Address: req.device.Address,
	}

	req.log.Info("Setting state")
	req.log.Debug("Getting capabilities")

	caps, err := req.driver.GetCapabilities(ctx, deviceInfo)
	if err != nil {
		req.log.Warn("unable to get capabilities", zap.Error(err))

		resp.errors = append(resp.errors, api.DeviceStateError{
			ID:    req.id,
			Error: fmt.Sprintf("unable to get capabilities: %s", status.Convert(err).Message()),
		})

		if len(resp.state.Inputs) == 0 {
			resp.state.Inputs = nil
		}

		if len(resp.state.Volumes) == 0 {
			resp.state.Volumes = nil
		}

		if len(resp.state.Mutes) == 0 {
			resp.state.Mutes = nil
		}

		return resp
	}

	hasCapability := func(c drivers.Capability) bool {
		for _, ability := range caps.GetCapabilities() {
			if ability == string(c) {
				return true
			}
		}

		return false
	}

	driverErr := func(field string, value interface{}, err error) {
		req.log.Warn("unable to set "+field, zap.Any("to", value), zap.Error(err))

		msg := err.Error()
		if err, ok := status.FromError(err); ok {
			msg = err.Message()
		}

		respMu.Lock()
		defer respMu.Unlock()

		resp.errors = append(resp.errors, api.DeviceStateError{
			ID:    req.id,
			Field: field,
			Value: value,
			Error: msg,
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

	req.log.Debug("Got capabilities", zap.Strings("capabilities", caps.GetCapabilities()))
	wg := sync.WaitGroup{}

	// try to set each of the fields that were passed in
	if req.state.PoweredOn != nil {
		if hasCapability(drivers.CapabilityPower) {
			powerReq := &drivers.SetPowerRequest{
				Info: deviceInfo,
				Power: &drivers.Power{
					On: *req.state.PoweredOn,
				},
			}

			wg.Add(1)
			spreadRequests()
			go func() {
				req.log.Info("Setting power", zap.Bool("on", powerReq.Power.On))
				defer wg.Done()

				if _, err := req.driver.SetPower(ctx, powerReq); err != nil {
					driverErr("poweredOn", powerReq.Power.On, err)
					return
				}

				req.log.Info("Set power")

				respMu.Lock()
				defer respMu.Unlock()
				resp.state.PoweredOn = &powerReq.Power.On
			}()
		} else {
			driverErr("poweredOn", *req.state.PoweredOn, ErrNotCapable)
		}
	}

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
		if hasCapability(drivers.CapabilityAudioInput) {
			for output, input := range req.state.Inputs {
				if input.Audio == nil {
					continue
				}

				inputReq := &drivers.SetInputRequest{
					Info:   deviceInfo,
					Output: output,
					Input:  *input.Audio,
				}

				wg.Add(1)
				spreadRequests()
				go func() {
					req.log.Info("Setting audio input", zap.String("output", inputReq.Output), zap.String("input", inputReq.Input))
					defer wg.Done()

					if _, err := req.driver.SetAudioInput(ctx, inputReq); err != nil {
						driverErr(fmt.Sprintf("input.%s.audio", inputReq.Output), inputReq.Input, err)
						return
					}

					req.log.Info("Set audio input", zap.String("output", inputReq.Output))

					respMu.Lock()
					defer respMu.Unlock()
					in := resp.state.Inputs[inputReq.Output]
					in.Audio = &inputReq.Input
					resp.state.Inputs[inputReq.Output] = in
				}()
			}
		} else {
			driverErr("input.$.audio", req.state.Inputs, ErrNotCapable)
		}
	}

	if setVideoInput {
		if hasCapability(drivers.CapabilityVideoInput) {
			for output, input := range req.state.Inputs {
				if input.Video == nil {
					continue
				}

				inputReq := &drivers.SetInputRequest{
					Info:   deviceInfo,
					Output: output,
					Input:  *input.Video,
				}

				wg.Add(1)
				spreadRequests()
				go func() {
					req.log.Info("Setting video input", zap.String("output", inputReq.Output), zap.String("input", inputReq.Input))
					defer wg.Done()

					if _, err := req.driver.SetVideoInput(ctx, inputReq); err != nil {
						driverErr(fmt.Sprintf("input.%s.video", inputReq.Output), inputReq.Input, err)
						return
					}

					req.log.Info("Set video input", zap.String("output", inputReq.Output))

					respMu.Lock()
					defer respMu.Unlock()
					in := resp.state.Inputs[inputReq.Output]
					in.Video = &inputReq.Input
					resp.state.Inputs[inputReq.Output] = in
				}()
			}
		} else {
			driverErr("input.$.video", req.state.Inputs, ErrNotCapable)
		}
	}

	if setAudioVideoInput {
		if hasCapability(drivers.CapabilityAudioVideoInput) {
			for output, input := range req.state.Inputs {
				if input.AudioVideo == nil {
					continue
				}

				inputReq := &drivers.SetInputRequest{
					Info:   deviceInfo,
					Output: output,
					Input:  *input.AudioVideo,
				}

				wg.Add(1)
				spreadRequests()
				go func() {
					req.log.Info("Setting audioVideo input", zap.String("output", inputReq.Output), zap.String("input", inputReq.Input))
					defer wg.Done()

					if _, err := req.driver.SetAudioVideoInput(ctx, inputReq); err != nil {
						driverErr(fmt.Sprintf("input.%s.audioVideo", inputReq.Output), inputReq.Input, err)
						return
					}

					req.log.Info("Set audioVideo input", zap.String("output", inputReq.Output))

					respMu.Lock()
					defer respMu.Unlock()
					in := resp.state.Inputs[inputReq.Output]
					in.AudioVideo = &inputReq.Input
					resp.state.Inputs[inputReq.Output] = in
				}()
			}
		} else {
			driverErr("input.$.audioVideo", req.state.Inputs, ErrNotCapable)
		}
	}

	if req.state.Blanked != nil {
		if hasCapability(drivers.CapabilityBlank) {
			blankReq := &drivers.SetBlankRequest{
				Info: deviceInfo,
				Blank: &drivers.Blank{
					Blanked: *req.state.Blanked,
				},
			}

			wg.Add(1)
			spreadRequests()
			go func() {
				req.log.Info("Setting blank", zap.Bool("blanked", blankReq.Blank.Blanked))
				defer wg.Done()

				if _, err := req.driver.SetBlank(ctx, blankReq); err != nil {
					driverErr("blanked", blankReq.Blank.Blanked, err)
					return
				}

				req.log.Info("Set blank")

				respMu.Lock()
				defer respMu.Unlock()
				resp.state.Blanked = &blankReq.Blank.Blanked
			}()
		} else {
			driverErr("blanked", *req.state.Blanked, ErrNotCapable)
		}
	}

	// TODO should we validate blocks from ports?
	if len(req.state.Volumes) > 0 {
		if hasCapability(drivers.CapabilityVolume) {
			for block, vol := range req.state.Volumes {
				volReq := &drivers.SetVolumeRequest{
					Info:  deviceInfo,
					Block: block,
					Level: int32(vol),
				}

				wg.Add(1)
				spreadRequests()
				go func() {
					req.log.Info("Setting volume", zap.String("block", volReq.Block), zap.Int32("level", volReq.Level))
					defer wg.Done()

					if _, err := req.driver.SetVolume(ctx, volReq); err != nil {
						driverErr(fmt.Sprintf("volumes.%s", volReq.Block), volReq.Level, err)
						return
					}

					req.log.Info("Set volume", zap.String("block", volReq.Block))

					respMu.Lock()
					defer respMu.Unlock()
					resp.state.Volumes[volReq.Block] = int(volReq.Level)
				}()
			}
		} else {
			driverErr("volumes", req.state.Volumes, ErrNotCapable)
		}
	}

	// TODO should we validate blocks from ports?
	if len(req.state.Mutes) > 0 {
		if hasCapability(drivers.CapabilityMute) {
			for block, muted := range req.state.Mutes {
				muteReq := &drivers.SetMuteRequest{
					Info:  deviceInfo,
					Block: block,
					Muted: muted,
				}

				wg.Add(1)
				spreadRequests()
				go func() {
					req.log.Info("Setting mute", zap.String("block", muteReq.Block), zap.Bool("muted", muteReq.Muted))
					defer wg.Done()

					if _, err := req.driver.SetMute(ctx, muteReq); err != nil {
						driverErr(fmt.Sprintf("mutes.%s", muteReq.Block), muteReq.Muted, err)
						return
					}

					req.log.Info("Set mute", zap.String("block", muteReq.Block))

					respMu.Lock()
					defer respMu.Unlock()
					resp.state.Mutes[muteReq.Block] = muteReq.Muted
				}()
			}
		} else {
			driverErr("mutes", req.state.Mutes, ErrNotCapable)
		}
	}

	wg.Wait()

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

	req.log.Info("Finished setting state")
	return resp
}
