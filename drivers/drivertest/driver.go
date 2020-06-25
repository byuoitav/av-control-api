package drivertest

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/byuoitav/av-control-api/drivers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type Driver struct {
	Devices map[string]drivers.Device
}

func (d *Driver) NewDeviceFunc() drivers.NewDeviceFunc {
	return func(ctx context.Context, addr string) (drivers.Device, error) {
		if dev, ok := d.Devices[addr]; ok {
			return dev, nil
		}

		return nil, fmt.Errorf("no device stored with address %s", addr)
	}
}

type Server struct {
	listener *bufconn.Listener
	Config   drivers.Server
}

func NewServer(newDev drivers.NewDeviceFunc) *Server {
	ts := NewUnstartedServer(newDev)
	ts.Start()
	return ts
}

func NewUnstartedServer(newDev drivers.NewDeviceFunc) *Server {
	return &Server{
		listener: bufconn.Listen(1024 * 1024),
		Config:   drivers.NewServer(newDev),
	}
}

func (s *Server) Start() {
	go s.Config.Serve(s.Listener())
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.Config.Stop(ctx)
	if err != nil {
		s.Listener().Close()
	}
}

func (s *Server) Listener() net.Listener {
	return s.listener
}

func (s *Server) GRPCClientConn(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, s.Listener().Addr().String(), grpc.WithContextDialer(bufConnDialer(s.listener)), grpc.WithInsecure())
}

func bufConnDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}
