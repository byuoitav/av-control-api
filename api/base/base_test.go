package base_test

import (
	"testing"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/matryer/is"
)

func TestActionByOrderLess(t *testing.T) {
	is := is.New(t)

	actions := []base.ActionStructure{
		base.ActionStructure{
			Action: "ValidAction",
			Device: base.Device{
				Type: base.DeviceType{
					Commands: map[string]base.Command{
						"ValidAction": base.Command{
							Order: 10,
						},
					},
				},
			},
		},
		base.ActionStructure{
			Action: "ValidAction2",
			Device: base.Device{
				Type: base.DeviceType{
					Commands: map[string]base.Command{
						"ValidAction2": base.Command{
							Order: 20,
						},
					},
				},
			},
		},
	}

	t.Run("Should return true when order of i is less than j", func(t *testing.T) {
		is := is.New(t)
		got := base.ActionByOrder(actions).Less(0, 1)
		is.Equal(got, true)
	})

	t.Run("Should return false when order of i is greater than j", func(t *testing.T) {
		is := is.New(t)
		got := base.ActionByOrder(actions).Less(1, 0)
		is.Equal(got, false)
	})

	t.Run("Should return false when order of i is equal to j", func(t *testing.T) {
		is := is.New(t)
		got := base.ActionByOrder(actions).Less(1, 1)
		is.Equal(got, false)
	})
}
