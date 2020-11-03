package couch

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

// I'm not even sure this is possible...
// func TestAuth(t *testing.T) {
// 	is := is.New(t)

// 	client, _, err := kivikmock.New()
// 	is.NoErr(err)

// 	_, err = NewWithClient(context.Background(), client)
// 	is.NoErr(err)
// }

func TestNew(t *testing.T) {
	is := is.New(t)

	_, err := New(context.TODO(), "")
	is.Equal("unable to build client: no URL specified", err.Error())
}
