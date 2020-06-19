package drivers

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func newGrpcServer(newDev NewDeviceFunc) Server {
	g := &grpcServer{
		newDevice: newDev,
		single:    &singleflight.Group{},
	}

	server := grpc.NewServer()
	RegisterDriverServer(server, g)

	return server
}

type grpcServer struct {
	newDevice NewDeviceFunc
	single    *singleflight.Group
}

func (g *grpcServer) GetCapabilities(ctx context.Context, info *DeviceInfo) (*Capabilities, error) {
	device, err := g.newDevice(ctx, info.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	var caps []string

	if _, ok := device.(DeviceWithPower); ok {
		caps = append(caps, string(CapabilityPower))
	}

	if _, ok := device.(DeviceWithAudioInput); ok {
		caps = append(caps, string(CapabilityAudioInput))
	}

	if _, ok := device.(DeviceWithVideoInput); ok {
		caps = append(caps, string(CapabilityVideoInput))
	}

	if _, ok := device.(DeviceWithAudioVideoInput); ok {
		caps = append(caps, string(CapabilityAudioVideoInput))
	}

	if _, ok := device.(DeviceWithBlank); ok {
		caps = append(caps, string(CapabilityBlank))
	}

	if _, ok := device.(DeviceWithVolume); ok {
		caps = append(caps, string(CapabilityVolume))
	}

	if _, ok := device.(DeviceWithMute); ok {
		caps = append(caps, string(CapabilityMute))
	}

	if _, ok := device.(DeviceWithInfo); ok {
		caps = append(caps, string(CapabilityInfo))
	}

	return &Capabilities{
		Capabilities: caps,
	}, nil
}

func (g *grpcServer) GetPower(ctx context.Context, info *DeviceInfo) (*Power, error) {
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

func (g *grpcServer) SetPower(ctx context.Context, req *SetPowerRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetAudioInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
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

func (g *grpcServer) SetAudioInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetVideoInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
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

func (g *grpcServer) SetVideoInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetAudioVideoInputs(ctx context.Context, info *DeviceInfo) (*Inputs, error) {
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

func (g *grpcServer) SetAudioVideoInput(ctx context.Context, req *SetInputRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetBlank(ctx context.Context, info *DeviceInfo) (*Blank, error) {
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

func (g *grpcServer) SetBlank(ctx context.Context, req *SetBlankRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetVolumes(ctx context.Context, info *GetAudioInfo) (*Volumes, error) {
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

func (g *grpcServer) SetVolume(ctx context.Context, req *SetVolumeRequest) (*empty.Empty, error) {
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

func (g *grpcServer) GetMutes(ctx context.Context, info *GetAudioInfo) (*Mutes, error) {
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

func (g *grpcServer) SetMute(ctx context.Context, req *SetMuteRequest) (*empty.Empty, error) {
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
