package couch

import "github.com/go-kivik/couchdb/v3"

const (
	_defaultDatabase = "av-control-api"
)

type options struct {
	authFunc interface{}
	database string
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
