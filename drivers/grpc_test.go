package drivers

import (
	"context"
	"net"
	reflect "reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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
			caps, err := client.GetCapabilities(ctx, info)
			if err != nil {
				t.Fatalf("unable to get capabilities: %s", err)
			}

			expected := &Capabilities{
				Capabilities: []string{
					string(CapabilityPower),
					string(CapabilityAudioVideoInput),
					string(CapabilityBlank),
					string(CapabilityVolume),
					string(CapabilityMute),
				},
			}

			opts := cmp.Options{
				cmp.Exporter(func(t reflect.Type) bool {
					return true
				}),
				cmpopts.IgnoreTypes(protoimpl.MessageState{}, protoimpl.UnknownFields{}),
			}

			if diff := cmp.Diff(expected, caps, opts...); diff != "" {
				t.Fatalf("generated incorrect response (-want, +got):\n%s", diff)
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
