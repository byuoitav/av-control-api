package couch

import "github.com/go-kivik/couchdb/v4"

const (
	_defaultScheme = "https"
)

type options struct {
	authFunc interface{}
	scheme   string
}

// Option configures how we create the Camera.
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
