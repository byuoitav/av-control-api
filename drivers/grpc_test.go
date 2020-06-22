package drivers

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/sync/errgroup"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type grpcDriverTest struct {
	name      string
	newDevice NewDeviceFunc
	test      func(context.Context, *testing.T, DriverClient)
}

var grpcDriverTests = []grpcDriverTest{
	grpcDriverTest{
		name: "TV/GetCapabilities",
		newDevice: func(context.Context, string) (Device, error) {
			return &mockTV{}, nil
		},
		test: func(ctx context.Context, t *testing.T, client DriverClient) {
			info := &DeviceInfo{}

			// try getting capabilities
			got, err := client.GetCapabilities(ctx, info)
			if err != nil {
				t.Fatalf("unable to get capabilities: %s", err)
			}

			want := &Capabilities{
				Capabilities: []string{
					string(CapabilityPower),
					string(CapabilityAudioVideoInput),
					string(CapabilityBlank),
					string(CapabilityVolume),
					string(CapabilityMute),
				},
			}

			opts := cmp.Options{
				cmpopts.IgnoreUnexported(Capabilities{}),
			}

			if diff := cmp.Diff(want, got, opts...); diff != "" {
				t.Fatalf("generated incorrect response (-want, +got):\n%s", diff)
			}
		},
	},
	grpcDriverTest{
		name: "TV/Power",
		newDevice: saveDevicesFunc(func(context.Context, string) (Device, error) {
			return &mockTV{}, nil
		}),
		test: func(ctx context.Context, t *testing.T, client DriverClient) {
			req := &SetPowerRequest{
				Info: &DeviceInfo{},
				Power: &Power{
					On: true,
				},
			}

			if _, err := client.SetPower(ctx, req); err != nil {
				t.Fatalf("unable to get set power: %s", err)
			}

			got, err := client.GetPower(ctx, req.GetInfo())
			if err != nil {
				t.Fatalf("unable to get power: %s", err)
			}

			opts := cmp.Options{
				cmpopts.IgnoreUnexported(Power{}),
			}

			if diff := cmp.Diff(req.GetPower(), got, opts...); diff != "" {
				t.Fatalf("generated incorrect response (-want, +got):\n%s", diff)
			}
		},
	},
	grpcDriverTest{
		name: "TV/CombineRequests",
		newDevice: saveDevicesFunc(func(context.Context, string) (Device, error) {
			return &mockTV{
				delay: time.Second,
			}, nil
		}),
		test: func(ctx context.Context, t *testing.T, client DriverClient) {
			req := &SetPowerRequest{
				Info: &DeviceInfo{},
				Power: &Power{
					On: true,
				},
			}

			var done1Time time.Time
			var done2Time time.Time

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				if _, err := client.SetPower(gctx, req); err != nil {
					return err
				}

				done1Time = time.Now()
				return nil
			})

			group.Go(func() error {
				if _, err := client.SetPower(gctx, req); err != nil {
					return err
				}

				done2Time = time.Now()
				return nil
			})

			if err := group.Wait(); err != nil {
				t.Fatalf("unable to set power: %s", err)
			}

			diff := done2Time.Sub(done1Time)
			if done1Time.After(done2Time) {
				diff = done1Time.Sub(done2Time)
			}

			if diff.Round(500*time.Millisecond) > 0 {
				t.Fatalf("server responded %v apart", diff)
			}
		},
	},
	grpcDriverTest{
		name: "TV/SeparateRequests",
		newDevice: saveDevicesFunc(func(context.Context, string) (Device, error) {
			return &mockTV{
				delay: time.Second,
			}, nil
		}),
		test: func(ctx context.Context, t *testing.T, client DriverClient) {
			var done1Time time.Time
			var done2Time time.Time

			group, gctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				req := &SetPowerRequest{
					sizeCache: 5000,
					Info:      &DeviceInfo{},
					Power: &Power{
						On: true,
					},
				}

				if _, err := client.SetPower(gctx, req); err != nil {
					return err
				}

				done1Time = time.Now()
				return nil
			})

			group.Go(func() error {
				req := &SetPowerRequest{
					Info: &DeviceInfo{},
					Power: &Power{
						On: false,
					},
				}

				if _, err := client.SetPower(gctx, req); err != nil {
					return err
				}

				done2Time = time.Now()
				return nil
			})

			if err := group.Wait(); err != nil {
				t.Fatalf("unable to set power: %s", err)
			}

			diff := done2Time.Sub(done1Time)
			if done1Time.After(done2Time) {
				diff = done1Time.Sub(done2Time)
			}

			if diff.Round(time.Second) != time.Second {
				t.Fatalf("different requests responded faster than they should have (%v)", diff)
			}
		},
	},
}

func bufConnDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestGRPC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tt := range grpcDriverTests {
		t.Run(tt.name, func(t *testing.T) {
			lis := bufconn.Listen(1024 * 1024)
			server := newGrpcServer(tt.newDevice)

			t.Cleanup(func() {
				if err := server.Stop(ctx); err != nil {
					t.Fatalf("unable to stop server: %s", err)
				}
			})

			go server.Serve(lis)

			conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithContextDialer(bufConnDialer(lis)), grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to dial server: %s", err)
			}

			t.Cleanup(func() {
				if err := conn.Close(); err != nil {
					t.Fatalf("unable to close client connection: %s", err)
				}
			})

			tt.test(ctx, t, NewDriverClient(conn))
		})
	}
}
