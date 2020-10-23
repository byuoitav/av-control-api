package couch

import (
	"testing"

	"github.com/matryer/is"
)

func TestOptions(t *testing.T) {
	is := is.New(t)

	opts := []Option{
		WithBasicAuth("user", "pass"),
		WithDatabase("db"),
	}

	options := options{
		database: _defaultDatabase,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	is.True(options.authFunc != nil)
	is.Equal(options.database, "db")
}
