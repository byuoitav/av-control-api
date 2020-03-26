package state_test

import (
	"testing"

	"github.com/byuoitav/av-control-api/api/state"
	"github.com/matryer/is"
)

func TestReplaceParameters(t *testing.T) {
	is := is.New(t)

	t.Run("Should replace all parameters when all are provided", func(t *testing.T) {
		is := is.New(t)
		url := "http://localhost:8026/{{address}}/output/{{output}}/input/{{input}}"
		expected := "http://localhost:8026/10.0.0.1/output/2/input/1"
		params := map[string]string{
			"address": "10.0.0.1",
			"output":  "2",
			"input":   "1",
		}

		got, err := state.ReplaceParameters(url, params)
		is.NoErr(err)           // Expected to run without error
		is.Equal(got, expected) // Expected the url returned to match the expected url
	})
	t.Run("Should return an error when a parameter is unused", func(t *testing.T) {
		is := is.New(t)
		url := "http://localhost:8026/{{address}}/output/{{output}}/input/"
		params := map[string]string{
			"address": "10.0.0.1",
			"output":  "2",
			"input":   "1",
		}

		_, err := state.ReplaceParameters(url, params)
		is.True(err != nil) // Expected to get an error
	})
	t.Run("Should return an error when a template parameter is not a provided parameter", func(t *testing.T) {
		is := is.New(t)
		url := "http://localhost:8026/{{address}}/output/{{output}}/input/{{output}}/{{foobar}}"
		params := map[string]string{
			"address": "10.0.0.1",
			"output":  "2",
			"input":   "1",
		}

		_, err := state.ReplaceParameters(url, params)
		is.True(err != nil) // Expected to get an error
	})
}
