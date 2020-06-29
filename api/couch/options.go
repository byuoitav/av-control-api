package couch

import "github.com/go-kivik/couchdb/v3"

const (
	_defaultScheme       = "https"
	_defaultDatabase     = "av-control-api"
	_defaultMappingDocID = "#driverMapping"
	_defaultEnvironment  = "default"
)

type options struct {
	authFunc     interface{}
	scheme       string
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

func WithInsecure() Option {
	return optionFunc(func(o *options) {
		o.scheme = "http"
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
