package couch

import (
	"errors"
	"fmt"

	_ "github.com/go-kivik/couchdb/v4"
	kivik "github.com/go-kivik/kivik/v4"
	"golang.org/x/net/context"
)

type DataService struct {
	client *kivik.Client
}

func New(ctx context.Context, addr string, opts ...Option) (*DataService, error) {
	options := options{
		scheme: _defaultScheme,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	addr = fmt.Sprintf("%s://%s", options.scheme, addr)

	client, err := kivik.New("couch", addr)
	if err != nil {
		return nil, fmt.Errorf("unable to build client: %w", err)
	}

	if options.authFunc != nil {
		if err := client.Authenticate(ctx, options.authFunc); err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
	}

	return &DataService{
		client: client,
	}, nil
}

// Healthy is a healthcheck for the database
func (d *DataService) Healthy(ctx context.Context) error {
	alive, err := d.client.Ping(ctx)
	switch {
	case err != nil:
		return err
	case !alive:
		return errors.New("not healthy")
	}

	return nil
}
