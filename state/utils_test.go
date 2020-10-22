package state

import (
	"testing"

	avcontrol "github.com/byuoitav/av-control-api"
	"github.com/google/go-cmp/cmp"
)

func boolP(b bool) *bool {
	return &b
}

func stringP(s string) *string {
	return &s
}

type sortErrorsTest struct {
	name string
	in   []avcontrol.DeviceStateError
	out  []avcontrol.DeviceStateError
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
		in: []avcontrol.DeviceStateError{
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
		out: []avcontrol.DeviceStateError{
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
		in: []avcontrol.DeviceStateError{
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
		out: []avcontrol.DeviceStateError{
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
		in: []avcontrol.DeviceStateError{
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
		out: []avcontrol.DeviceStateError{
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
