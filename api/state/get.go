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
)

var (
	ErrNoStateGettable = errors.New("can't get the state of any devices in this room")
	ErrUnknownDriver   = errors.New("unknown driver")
)

// Get .
func (gs *getSetter) Get(ctx context.Context, room api.Room) (api.StateResponse, error) {
	resp := api.StateResponse{
		Devices: make(map[api.DeviceID]api.DeviceState),
	}

	id := api.RequestID(ctx)
	log := gs.log.With(zap.String("requestID", id))

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		_, ok := gs.drivers[dev.Driver]
		if !ok {
			return resp, fmt.Errorf("%w: %s", ErrUnknownDriver, dev.Driver)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(room.Devices))

	getState := func(id api.DeviceID, dev api.Device) {
		defer wg.Done()

		driver := gs.drivers[dev.Driver]
		log := log.With(zap.String("deviceID", string(id)))

		state, errors := getStateForDevice(ctx, id, dev, driver, log)

		resp.Devices[id] = state
		resp.Errors = append(resp.Errors, errors...)
	}

	// TODO mutex on resp
	for id, dev := range room.Devices {
		go getState(id, dev)
	}

	wg.Wait()
	return resp, nil
}

func getStateForDevice(ctx context.Context, id api.DeviceID, dev api.Device, driver drivers.DriverClient, log api.Logger) (api.DeviceState, []api.DeviceStateError) {
	var state api.DeviceState
	var errors []api.DeviceStateError

	deviceInfo := &drivers.DeviceInfo{
		Address: dev.Address,
	}

	log.Info("Getting state")
	log.Debug("Getting capabilities")

	// TODO status (grpc) on errors
	caps, err := driver.GetCapabilities(ctx, deviceInfo)
	if err != nil {
		log.Warn("unable to get capabilities", zap.Error(err))

		errors = append(errors, api.DeviceStateError{
			ID:    id,
			Error: fmt.Sprintf("unable to get capabilities: %s", err),
		})
		return state, errors
	}

	// TODO mutex on state/errors
	log.Debug("Got capabilities", zap.Strings("capabilities", caps.GetCapabilities()))
	wg := sync.WaitGroup{}

	for _, capability := range caps.GetCapabilities() {
		switch drivers.Capability(capability) {
		case drivers.CapabilityPower:
			wg.Add(1)

			go func() {
				log.Info("Getting power")
				defer wg.Done()

				power, err := driver.GetPower(ctx, deviceInfo)
				if err != nil {
					errors = append(errors, api.DeviceStateError{
						ID:    id,
						Field: "power",
						Error: err.Error(),
					})
					return
				}

				log.Info("Got power", zap.Bool("on", power.On))
				state.PoweredOn = &power.On
			}()
		default:
			log.Warn("unknown capability", zap.String("capability", capability))

			errors = append(errors, api.DeviceStateError{
				ID:    id,
				Error: fmt.Sprintf("don't know how to handle capability %q", capability),
			})
			continue
		}

		// TODO don't wait on the last one
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	log.Info("Finished getting state")
	return state, errors
}
