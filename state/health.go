package state

import (
	"context"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type getDeviceHealthRequest struct {
	id     avcontrol.DeviceID
	device avcontrol.DeviceConfig
	driver avcontrol.Driver
	log    *zap.Logger
}

type getDeviceHealthResponse struct {
	id     avcontrol.DeviceID
	health avcontrol.DeviceHealth
}

func (gs *GetSetter) GetHealth(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomHealth, error) {
	if len(room.Devices) == 0 {
		return avcontrol.RoomHealth{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		if gs.DriverRegistry.Get(dev.Driver) == nil {
			return avcontrol.RoomHealth{}, fmt.Errorf("%s: %w", dev.Driver, ErrDriverNotRegistered)
		}
	}

	id := avcontrol.CtxRequestID(ctx)
	log := gs.Logger
	if len(id) > 0 {
		log = gs.Logger.With(zap.String("requestID", id))
	}

	resps := make(chan getDeviceHealthResponse)
	defer close(resps)

	roomHealth := avcontrol.RoomHealth{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceHealth),
	}

	for id, dev := range room.Devices {
		req := getDeviceHealthRequest{
			id:     id,
			device: dev,
			driver: gs.DriverRegistry.Get(dev.Driver),
			log:    log.With(zap.String("deviceID", string(id))),
		}

		go func() {
			resps <- req.do(ctx)
		}()
	}

	for i := 0; i < len(room.Devices); i++ {
		resp := <-resps
		roomHealth.Devices[resp.id] = resp.health
	}

	return roomHealth, nil
}

func (req *getDeviceHealthRequest) do(ctx context.Context) getDeviceHealthResponse {
	resp := getDeviceHealthResponse{
		id: req.id,
	}

	req.log.Info("Getting health")
	req.log.Debug("Getting device")

	dev, err := req.driver.CreateDevice(ctx, req.device.Address)
	if err != nil {
		req.log.Warn("unable to get device", zap.Error(err))
		str := fmt.Sprintf("unable to get device: %s", err)
		resp.health.Error = &str
		return resp
	}

	req.log.Debug("Got device")

	if dev, ok := dev.(avcontrol.DeviceWithHealth); ok {
		healthy := true
		if err := dev.Healthy(ctx); err != nil {
			req.log.Warn("unable to get health", zap.Error(err))
			str := err.Error()
			resp.health.Error = &str
			healthy = false
		}

		resp.health.Healthy = &healthy
	}

	req.log.Info("Finished getting health")

	return resp
}
