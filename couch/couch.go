package couch

import (
	"fmt"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	"golang.org/x/net/context"
)

type DataService struct {
	client   *kivik.Client
	database string
}

func New(ctx context.Context, url string, opts ...Option) (*DataService, error) {
	client, err := kivik.New("couch", url)
	if err != nil {
		return nil, fmt.Errorf("unable to build client: %w", err)
	}

	return NewWithClient(ctx, client, opts...)
}

func NewWithClient(ctx context.Context, client *kivik.Client, opts ...Option) (*DataService, error) {
	options := options{
		database: _defaultDatabase,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	if options.authFunc != nil {
		if err := client.Authenticate(ctx, options.authFunc); err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
	}

	return &DataService{
		client:   client,
		database: options.database,
	}, nil
}
