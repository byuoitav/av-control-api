package api

import "testing"

var deviceIDTests = []struct {
	id   string
	room string
}{
	{"ITB-1101-CP1", "ITB-1101"},
	{"EB-101-D1", "EB-101"},
	{"ITB-1101", "ITB-1101"},
	{"hello", "hello"},
	{"", ""},
}

func TestDeviceID(t *testing.T) {
	for _, tt := range deviceIDTests {
		actual := DeviceID(tt.id).Room()
		if actual != tt.room {
			t.Errorf("(%q).Room(): expected %q, got %q", tt.id, tt.room, actual)
		}
	}
}
