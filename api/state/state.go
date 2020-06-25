package state

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/byuoitav/av-control-api/api"
	"github.com/byuoitav/av-control-api/drivers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type getSetter struct {
	log api.Logger

	drivers map[string]drivers.DriverClient
}

func New(ctx context.Context, ds api.DataService, log api.Logger) (api.StateGetSetter, error) {
	gs := &getSetter{
		log:     log,
		drivers: make(map[string]drivers.DriverClient),
	}

	mapping, err := ds.DriverMapping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get driver mapping: %w", err)
	}

	var grpcCreds credentials.TransportCredentials

	for driver, config := range mapping {
		var opts []grpc.DialOption

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

		gs.log.Debug(fmt.Sprintf("Setting up %#q driver", driver), zap.String("address", config.Address))

		conn, err := grpc.DialContext(ctx, config.Address, opts...)
		if err != nil {
			return nil, fmt.Errorf("unable to dial driver (%s): %w", config.Address, err)
		}

		gs.drivers[driver] = drivers.NewDriverClient(conn)
	}

	return gs, nil
}
