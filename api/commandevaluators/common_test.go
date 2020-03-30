package commandevaluators_test

import (
	"testing"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/av-control-api/api/commandevaluators"
	"github.com/matryer/is"
)

func TestCheckCommands(t *testing.T) {
	is := is.New(t)

	t.Run("Should find the command if one exists that matches", func(t *testing.T) {
		is := is.New(t)
		c := map[string]base.Command{
			"ValidCommand": base.Command{
				Order: 100,
			},
		}
		exists, got := commandevaluators.CheckCommands(c, "ValidCommand")
		is.True(exists)                  // Expected command to be found
		is.Equal(got, c["ValidCommand"]) // Expected the correct command to be returned
	})

	t.Run("Should return false if the command does not exist", func(t *testing.T) {

	})
}
