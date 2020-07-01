package couch

import (
	"fmt"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	"golang.org/x/net/context"
)

type DataService struct {
	client       *kivik.Client
	database     string
	mappingDocID string
	environment  string
}

func New(ctx context.Context, addr string, opts ...Option) (*DataService, error) {
	options := options{
		scheme:       _defaultScheme,
		database:     _defaultDatabase,
		mappingDocID: _defaultMappingDocID,
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
		client:       client,
		database:     options.database,
		mappingDocID: options.mappingDocID,
		environment:  options.environment,
	}, nil
}
