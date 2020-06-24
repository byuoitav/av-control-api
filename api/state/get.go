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
	ErrNoStateGettable = errors.New("can't get the state of any devices in this room")
	ErrUnknownDriver   = errors.New("unknown driver")
)

type getDeviceStateRequest struct {
	id     api.DeviceID
	device api.Device
	driver drivers.DriverClient
	log    api.Logger
}

type getDeviceStateResponse struct {
	id     api.DeviceID
	state  api.DeviceState
	errors []api.DeviceStateError
}

func (gs *getSetter) Get(ctx context.Context, room api.Room) (api.StateResponse, error) {
	stateResp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
	}

	id := api.RequestID(ctx)
	log := gs.log.With(zap.String("requestID", id))

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		_, ok := gs.drivers[dev.Driver]
		if !ok {
			return stateResp, fmt.Errorf("%w: %s", ErrUnknownDriver, dev.Driver)
		}
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
			resps <- req.Do(ctx)

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

	close(resps)

	return stateResp, nil
}

func (req *getDeviceStateRequest) Do(ctx context.Context) getDeviceStateResponse {
	var respMu sync.Mutex
	resp := getDeviceStateResponse{
		id: req.id,
		state: api.DeviceState{
			Inputs: make(map[string]api.Input),
		},
	}

	deviceInfo := &drivers.DeviceInfo{
		Address: req.device.Address,
	}

	req.log.Info("Getting state")
	req.log.Debug("Getting capabilities")

	caps, err := req.driver.GetCapabilities(ctx, deviceInfo)
	if err != nil {
		req.log.Warn("unable to get capabilities", zap.Error(err))

		resp.errors = append(resp.errors, api.DeviceStateError{
			ID:    req.id,
			Error: fmt.Sprintf("unable to get capabilities: %s", status.Convert(err).Message()),
		})

		return resp
	}

	driverErr := func(field string, err error) {
		req.log.Warn("unable to get "+field, zap.Error(err))

		respMu.Lock()
		defer respMu.Unlock()

		resp.errors = append(resp.errors, api.DeviceStateError{
			ID:    req.id,
			Field: field,
			Error: status.Convert(err).Message(),
		})
	}

	spreadRequests := func(i int) {
		if i == 0 {
			return
		}

		time.Sleep(25 * time.Millisecond)
	}

	req.log.Debug("Got capabilities", zap.Strings("capabilities", caps.GetCapabilities()))
	wg := sync.WaitGroup{}

	for i, capability := range caps.GetCapabilities() {
		switch drivers.Capability(capability) {
		case drivers.CapabilityPower:
			wg.Add(1)
			spreadRequests(i)

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
			spreadRequests(i)

			go func() {
				req.log.Info("Getting audio input")
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
					input.Audio = &in
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityVideoInput:
			wg.Add(1)
			spreadRequests(i)

			go func() {
				req.log.Info("Getting video input")
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
					input.AudioVideo = &in
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityAudioVideoInput:
			wg.Add(1)
			spreadRequests(i)

			go func() {
				req.log.Info("Getting audio-video input")
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
					input.AudioVideo = &in
					resp.state.Inputs[out] = input
				}
			}()
		case drivers.CapabilityBlank:
			wg.Add(1)
			spreadRequests(i)

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
			spreadRequests(i)

			go func() {
				req.log.Info("Getting volumes")
				defer wg.Done()

				audioInfo := &drivers.GetAudioInfo{
					Info:   deviceInfo,
					Blocks: req.device.TypePorts("audio").Names(),
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
			spreadRequests(i)

			go func() {
				req.log.Info("Getting mutes")
				defer wg.Done()

				audioInfo := &drivers.GetAudioInfo{
					Info:   deviceInfo,
					Blocks: req.device.TypePorts("audio").Names(),
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
	return resp
}
