package drivers

import (
	"errors"
	"sort"
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
	is := is.New(t)

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)
	is.True(len(r.List()) == 0)

	drivers := []string{
		"driver/0",
		"driver/1",
		"driver/2",
		"driver/3",
		"driver/4",
		"driver/5",
	}

	sort.Strings(drivers)

	for i, driver := range drivers {
		r.MustRegister(driver, &testDriver{})
		list := r.List()
		sort.Strings(list)
		is.Equal(list, drivers[:i+1])
	}
}

func TestMustRegister(t *testing.T) {
	is := is.New(t)
	defer func() {
		r := recover()
		is.True(r != nil)
		is.True(strings.Contains(r.(error).Error(), "already registered"))
	}()

	r, err := NewWithConfig(make(map[string]map[string]interface{}))
	is.NoErr(err)

	r.MustRegister("driver/name", &testDriver{})
	r.MustRegister("driver/name", &testDriver{})

}
