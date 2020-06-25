package drivers

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type wrappedGrpcServer struct {
	*grpc.Server
}

func (w *wrappedGrpcServer) Stop(ctx context.Context) error {
	w.Server.Stop()
	return nil
}

func newGrpcServer(newDev NewDeviceFunc) Server {
	g := &grpcDriverServer{
		newDevice: newDev,
		single:    &singleflight.Group{},
	}

	server := &wrappedGrpcServer{
		Server: grpc.NewServer(),
	}

	RegisterDriverServer(server.Server, g)
	return server
}

type grpcDriverServer struct {
	newDevice NewDeviceFunc
	single    *singleflight.Group
}

func (g *grpcDriverServer) GetCapabilities(ctx context.Context, info *DeviceInfo) (*Capabilities, error) {
	dev, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if d, ok := dev.(DeviceWithCapabilities); ok {
		caps, err := d.GetCapabilities(ctx)
		if err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}

		return &Capabilities{
			Capabilities: caps,
		}, nil
	}

	var caps []string

	if _, ok := dev.(DeviceWithPower); ok {
		caps = append(caps, string(CapabilityPower))
	}

	if _, ok := dev.(DeviceWithAudioInput); ok {
		caps = append(caps, string(CapabilityAudioInput))
	}

	if _, ok := dev.(DeviceWithVideoInput); ok {
		caps = append(caps, string(CapabilityVideoInput))
	}

	if _, ok := dev.(DeviceWithAudioVideoInput); ok {
		caps = append(caps, string(CapabilityAudioVideoInput))
	}

	if _, ok := dev.(DeviceWithBlank); ok {
		caps = append(caps, string(CapabilityBlank))
	}

	if _, ok := dev.(DeviceWithVolume); ok {
		caps = append(caps, string(CapabilityVolume))
	}

	if _, ok := dev.(DeviceWithMute); ok {
		caps = append(caps, string(CapabilityMute))
	}

	if _, ok := dev.(DeviceWithInfo); ok {
		caps = append(caps, string(CapabilityInfo))
	}

	return &Capabilities{
		Capabilities: caps,
	}, nil
}

func (g *grpcDriverServer) GetPower(ctx context.Context, info *DeviceInfo) (*Power, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithPower)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetPower not implemented")
	}

	val, err, _ := g.single.Do("GetPower"+info.String(), func() (interface{}, error) {
		pow, err := dev.GetPower(ctx)
		if err != nil {
			return nil, err
		}

		return pow, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Power{
		On: val.(bool),
	}, nil
}

func (g *grpcDriverServer) SetPower(ctx context.Context, req *SetPowerRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithPower)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetPower not implemented")
	}

	_, err, _ = g.single.Do("SetPower"+req.String(), func() (interface{}, error) {
		return nil, dev.SetPower(ctx, req.GetPower().GetOn())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetAudioInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithAudioInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetAudioInputs not implemented")
	}

	val, err, _ := g.single.Do("GetAudioInputs"+info.String(), func() (interface{}, error) {
		inputs, err := dev.GetAudioInputs(ctx)
		if err != nil {
			return nil, err
		}

		return inputs, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Inputs{
		Inputs: val.(map[string]string),
	}, nil
}

func (g *grpcDriverServer) SetAudioInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithAudioInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetAudioInput not implemented")
	}

	_, err, _ = g.single.Do("SetAudioInput"+req.String(), func() (interface{}, error) {
		return nil, dev.SetAudioInput(ctx, req.GetOutput(), req.GetInput())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetVideoInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithVideoInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetVideoInputs not implemented")
	}

	val, err, _ := g.single.Do("GetVideoInputs"+info.String(), func() (interface{}, error) {
		inputs, err := dev.GetVideoInputs(ctx)
		if err != nil {
			return nil, err
		}

		return inputs, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Inputs{
		Inputs: val.(map[string]string),
	}, nil
}

func (g *grpcDriverServer) SetVideoInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithVideoInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetVideoInput not implemented")
	}

	_, err, _ = g.single.Do("SetVideoInput"+req.String(), func() (interface{}, error) {
		return nil, dev.SetVideoInput(ctx, req.GetOutput(), req.GetInput())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetAudioVideoInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithAudioVideoInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetAudioVideoInputs not implemented")
	}

	val, err, _ := g.single.Do("GetAudioVideoInputs"+info.String(), func() (interface{}, error) {
		inputs, err := dev.GetAudioVideoInputs(ctx)
		if err != nil {
			return nil, err
		}

		return inputs, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Inputs{
		Inputs: val.(map[string]string),
	}, nil
}

func (g *grpcDriverServer) SetAudioVideoInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithAudioVideoInput)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetAudioVideoInput not implemented")
	}

	_, err, _ = g.single.Do("SetAudioVideoInput"+req.String(), func() (interface{}, error) {
		return nil, dev.SetAudioVideoInput(ctx, req.GetOutput(), req.GetInput())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetBlank(ctx context.Context, info *DeviceInfo) (*Blank, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithBlank)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetBlank not implemented")
	}

	val, err, _ := g.single.Do("GetBlank"+info.String(), func() (interface{}, error) {
		blanked, err := dev.GetBlank(ctx)
		if err != nil {
			return nil, err
		}

		return blanked, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Blank{
		Blanked: val.(bool),
	}, nil
}

func (g *grpcDriverServer) SetBlank(ctx context.Context, req *SetBlankRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithBlank)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetBlank not implemented")
	}

	_, err, _ = g.single.Do("SetBlank"+req.String(), func() (interface{}, error) {
		return nil, dev.SetBlank(ctx, req.GetBlank().GetBlanked())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetVolumes(ctx context.Context, info *GetAudioInfo) (*Volumes, error) {
	device, err := g.newDevice(ctx, info.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithVolume)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetVolumes not implemented")
	}

	val, err, _ := g.single.Do("GetVolumes"+info.String(), func() (interface{}, error) {
		volumes, err := dev.GetVolumes(ctx, info.GetBlocks())
		if err != nil {
			return nil, err
		}

		return volumes, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	// have to cast each value
	vols := make(map[string]int32)
	for k, v := range val.(map[string]int) {
		vols[k] = int32(v)
	}

	return &Volumes{
		Volumes: vols,
	}, nil
}

func (g *grpcDriverServer) SetVolume(ctx context.Context, req *SetVolumeRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithVolume)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetVolume not implemented")
	}

	_, err, _ = g.single.Do("SetVolume"+req.String(), func() (interface{}, error) {
		return nil, dev.SetVolume(ctx, req.GetBlock(), int(req.GetLevel()))
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g *grpcDriverServer) GetMutes(ctx context.Context, info *GetAudioInfo) (*Mutes, error) {
	device, err := g.newDevice(ctx, info.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithMute)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method GetMutes not implemented")
	}

	val, err, _ := g.single.Do("GetMutes"+info.String(), func() (interface{}, error) {
		mutes, err := dev.GetMutes(ctx, info.GetBlocks())
		if err != nil {
			return nil, err
		}

		return mutes, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &Mutes{
		Mutes: val.(map[string]bool),
	}, nil
}

func (g *grpcDriverServer) SetMute(ctx context.Context, req *SetMuteRequest) (*empty.Empty, error) {
	device, err := g.newDevice(ctx, req.GetInfo().GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	dev, ok := device.(DeviceWithMute)
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "method SetMute not implemented")
	}

	_, err, _ = g.single.Do("SetMute"+req.String(), func() (interface{}, error) {
		return nil, dev.SetMute(ctx, req.GetBlock(), req.GetMuted())
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}
