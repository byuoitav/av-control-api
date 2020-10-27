package drivers

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestRegisterSuccess(t *testing.T) {
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	err = r.Register("driver/name", &testDriver{})
	is.NoErr(err)
}

func TestRegisterEmptyName(t *testing.T) {
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	err = r.Register("", &testDriver{})
	is.Equal(err.Error(), "driver must have a name")
}

func TestRegisterDuplicate(t *testing.T) {
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	err = r.Register("driver/name", &testDriver{})
	is.NoErr(err)

	err = r.Register("driver/name", &testDriver{})
	is.True(strings.Contains(err.Error(), "already registered"))
}

func TestRegisterParseError(t *testing.T) {
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	parseError := errors.New("parse config error")

	err = r.Register("driver/name", &testDriver{
		parseConfigErr: func() error {
			return parseError
		},
	})
	is.True(errors.Is(err, parseError))
}

func TestGet(t *testing.T) {
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	driver := &testDriver{}
	err = r.Register("driver/0", driver)
	is.NoErr(err)

	err = r.Register("driver/1", &testDriver{
		parseConfigErr: func() error {
			return errors.New("parse error")
		},
	})
	is.True(err != nil)

	d := r.Get("driver/0")
	is.True(d.(*deviceCache).Driver == driver)

	d = r.Get("driver/1")
	is.True(d == nil)
}

func TestList(t *testing.T) {
}

func TestMustRegister(t *testing.T) {
}
