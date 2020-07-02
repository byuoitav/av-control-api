package state

import (
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"github.com/google/go-cmp/cmp"
)

type sortErrorsTest struct {
	name string
	in   []api.DeviceStateError
	out  []api.DeviceStateError
}

type containsStringTest struct {
	name     string
	slice    []string
	s        string
	expected bool
}

var sortErrorsTests = []sortErrorsTest{
	{
		name: "Empty",
	},
	{
		name: "ID",
		in: []api.DeviceStateError{
			{
				ID: "ITB-1101-D3",
			},
			{
				ID: "ITB-1101-D1",
			},
			{
				ID: "ITB-1101-D9",
			},
			{
				ID: "ITB-1101-D2",
			},
		},
		out: []api.DeviceStateError{
			{
				ID: "ITB-1101-D1",
			},
			{
				ID: "ITB-1101-D2",
			},
			{
				ID: "ITB-1101-D3",
			},
			{
				ID: "ITB-1101-D9",
			},
		},
	},
	{
		name: "IDAndField",
		in: []api.DeviceStateError{
			{
				ID:    "ITB-1101-D2",
				Field: "volumes.aux",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "volumes.aux",
			},
		},
		out: []api.DeviceStateError{
			{
				ID:    "ITB-1101-D1",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "volumes.aux",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "volumes.aux",
			},
		},
	},
	{
		name: "IDAndFieldAndError",
		in: []api.DeviceStateError{
			{
				ID: "ITB-1101-D3",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "blank",
				Error: "invalid option",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "blank",
				Error: "can't blank",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "volumes.aux",
			},
		},
		out: []api.DeviceStateError{
			{
				ID:    "ITB-1101-D1",
				Field: "blank",
				Error: "can't blank",
			},
			{
				ID:    "ITB-1101-D1",
				Field: "blank",
				Error: "invalid option",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "mutes.aux",
			},
			{
				ID:    "ITB-1101-D2",
				Field: "volumes.aux",
			},
			{
				ID: "ITB-1101-D3",
			},
		},
	},
}

var containsStringTests = []containsStringTest{
	{
		name:     "Contains",
		slice:    []string{"1", "out", "AudioBlock", "3"},
		s:        "AudioBlock",
		expected: true,
	},
	{
		name:     "NotContains",
		slice:    []string{"1", "out", "AudioBlock", "3"},
		s:        "2",
		expected: false,
	},
}

func TestSortErrors(t *testing.T) {
	for _, tt := range sortErrorsTests {
		t.Run(tt.name, func(t *testing.T) {
			sortErrors(tt.in)

			if diff := cmp.Diff(tt.out, tt.in); diff != "" {
				t.Fatalf("incorrect order/contents (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	for _, tt := range containsStringTests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsString(tt.slice, tt.s)
			if got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
