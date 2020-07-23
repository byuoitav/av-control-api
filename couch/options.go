package couch

import "github.com/go-kivik/couchdb/v3"

const (
	_defaultDatabase     = "av-control-api"
	_defaultMappingDocID = "#driverMapping"
)

type options struct {
	authFunc     interface{}
	database     string
	mappingDocID string
	environment  string
}

// Option configures how we create the DataService.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithBasicAuth(username, password string) Option {
	return optionFunc(func(o *options) {
		o.authFunc = couchdb.BasicAuth(username, password)
	})
}

func WithDatabase(database string) Option {
	return optionFunc(func(o *options) {
		o.database = database
	})
}

func WithMappingDocumentID(docID string) Option {
	return optionFunc(func(o *options) {
		o.mappingDocID = docID
	})
}

func WithEnvironment(env string) Option {
	return optionFunc(func(o *options) {
		o.environment = env
	})
}
