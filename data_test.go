package avcontrol

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var portNamesTests = []struct {
	name  string
	ports PortConfigs
	names []string
}{
	{
		name:  "Empty",
		ports: PortConfigs{},
		names: nil,
	},
	{
		name: "Normal",
		ports: PortConfigs{
			{"1", "audio"},
			{"2", "video"},
			{"3", "audioVideo"},
			{"4", ""},
		},
		names: []string{"1", "2", "3", "4"},
	},
}

func TestPortNames(t *testing.T) {
	for _, tt := range portNamesTests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.names, tt.ports.Names()); diff != "" {
				t.Fatalf("got incorrect names (-want, +got):\n%s", diff)
			}
		})
	}
}

var portTypeTests = []struct {
	name string
	in   PortConfigs
	typ  string
	out  PortConfigs
}{
	{
		typ: "audio",
		in: PortConfigs{
			{"1", "audio"},
			{"2", "video"},
			{"3", "audioVideo"},
			{"4", ""},
		},
		out: PortConfigs{
			{"1", "audio"},
			{"3", "audioVideo"},
		},
	},
	{
		typ: "video",
		in: PortConfigs{
			{"1", "audio"},
			{"2", "video"},
			{"3", "audio-video"},
			{"4", ""},
		},
		out: PortConfigs{
			{"2", "video"},
			{"3", "audio-video"},
		},
	},
}

func TestPortTypes(t *testing.T) {
	for _, tt := range portTypeTests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.out, tt.in.OfType(tt.typ)); diff != "" {
				t.Fatalf("got incorrect names (-want, +got):\n%s", diff)
			}
		})
	}
}
