package state

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

type getSetter struct {
	logger *zap.Logger

	drivers      map[string]drivers.DriverClient
	driverStates []func() (string, string)
}

func New(ctx context.Context, ds api.DataService, logger *zap.Logger) (api.StateGetSetter, error) {
	gs := &getSetter{
		logger:  logger,
		drivers: make(map[string]drivers.DriverClient),
	}

	mapping, err := ds.DriverMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get driver mapping: %w", err)
	}

	var grpcCreds credentials.TransportCredentials

	for driver, config := range mapping {
		opts := []grpc.DialOption{
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.4,
					Jitter:     0.2,
					MaxDelay:   30 * time.Second,
				},
				MinConnectTimeout: 7500 * time.Millisecond,
			}),
		}

		if config.SSL {
			if grpcCreds == nil {
				pool, err := x509.SystemCertPool()
				if err != nil {
					return nil, fmt.Errorf("unable to get system cert pool: %v", err)
				}

				grpcCreds = credentials.NewClientTLSFromCert(pool, "")
			}

			opts = append(opts, grpc.WithTransportCredentials(grpcCreds))
		} else {
			opts = append(opts, grpc.WithInsecure())
		}

		gs.logger.Debug(fmt.Sprintf("Setting up %#q driver", driver), zap.String("address", config.Address))

		conn, err := grpc.DialContext(ctx, config.Address, opts...)
		if err != nil {
			return nil, fmt.Errorf("unable to dial driver (%s): %w", config.Address, err)
		}

		gs.drivers[driver] = drivers.NewDriverClient(conn)

		saved := driver
		gs.driverStates = append(gs.driverStates, func() (string, string) {
			return saved, conn.GetState().String()
		})
	}

	return gs, nil
}

func (gs *getSetter) DriverStates(ctx context.Context) (map[string]string, error) {
	states := make(map[string]string)

	for i := range gs.driverStates {
		driver, state := gs.driverStates[i]()
		states[driver] = state
	}

	return states, nil
}
