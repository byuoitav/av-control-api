package state

import (
	"context"
	"fmt"

	avcontrol "github.com/byuoitav/av-control-api"
	"go.uber.org/zap"
)

type getDeviceInfoRequest struct {
	id     avcontrol.DeviceID
	device avcontrol.DeviceConfig
	driver avcontrol.Driver
	log    *zap.Logger
}

type getDeviceInfoResponse struct {
	id   avcontrol.DeviceID
	info avcontrol.DeviceInfo
}

func (gs *GetSetter) GetInfo(ctx context.Context, room avcontrol.RoomConfig) (avcontrol.RoomInfo, error) {
	if len(room.Devices) == 0 {
		return avcontrol.RoomInfo{}, nil
	}

	// make sure the driver for every device in the room exists
	for _, dev := range room.Devices {
		if gs.DriverRegistry.Get(dev.Driver) == nil {
			return avcontrol.RoomInfo{}, fmt.Errorf("%s: %w", dev.Driver, ErrDriverNotRegistered)
		}
	}

	id := avcontrol.CtxRequestID(ctx)
	log := gs.Logger
	if len(id) > 0 {
		log = gs.Logger.With(zap.String("requestID", id))
	}

	resps := make(chan getDeviceInfoResponse)
	defer close(resps)

	roomHealth := avcontrol.RoomInfo{
		Devices: make(map[avcontrol.DeviceID]avcontrol.DeviceInfo),
	}

	for id, dev := range room.Devices {
		req := getDeviceInfoRequest{
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
		roomHealth.Devices[resp.id] = resp.info
	}

	return roomHealth, nil
}

func (req *getDeviceInfoRequest) do(ctx context.Context) getDeviceInfoResponse {
	resp := getDeviceInfoResponse{
		id: req.id,
	}

	req.log.Info("Getting info")
	req.log.Debug("Getting device")

	dev, err := req.driver.CreateDevice(ctx, req.device.Address)
	if err != nil {
		req.log.Warn("unable to get device", zap.Error(err))
		str := fmt.Sprintf("unable to get device: %s", err)
		resp.info.Error = &str
		return resp
	}

	req.log.Debug("Got device")

	if dev, ok := dev.(avcontrol.DeviceWithInfo); ok {
		info, err := dev.Info(ctx)
		if err != nil {
			req.log.Warn("unable to get info", zap.Error(err))
			str := err.Error()
			resp.info.Error = &str
			return resp
		}

		resp.info.Info = &info
	}

	req.log.Info("Finished getting info")

	return resp
}
