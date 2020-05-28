package graph

import (
	"testing"

	"github.com/byuoitav/av-control-api/api"
	"gonum.org/v1/gonum/graph/path"
)

func TestFindPathToNeighbor(t *testing.T) {
	room := []api.Device{
		api.Device{
			ID: "ITB-1101-HDMI2",
			Ports: []api.Port{
				api.Port{
					Endpoints: []api.DeviceID{"ITB-1101-D1"},
					Outgoing:  true,
					Type:      "video",
				},
				api.Port{
					Endpoints: []api.DeviceID{"ITB-1101-D2"},
					Outgoing:  true,
					Type:      "video",
				},
				api.Port{
					Endpoints: []api.DeviceID{"ITB-1101-D3"},
					Outgoing:  true,
					Type:      "video",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-D2",
			Ports: []api.Port{
				api.Port{
					Name:      "hdmi!1",
					Endpoints: []api.DeviceID{"ITB-1101-HDMI1"},
					Incoming:  true,
					Type:      "video",
				},
				api.Port{
					Name:      "hdmi!2",
					Endpoints: []api.DeviceID{"ITB-1101-HDMI2"},
					Incoming:  true,
					Type:      "video",
				},
				api.Port{
					Name:      "hdmi!3",
					Endpoints: []api.DeviceID{"ITB-1101-VIA1"},
					Incoming:  true,
					Type:      "video",
				},
			},
		},
	}

	g := NewGraph(room, "video")
	shortest := path.DijkstraAllPaths(g)

	pathEdges := PathFromTo(g, &shortest, room[0].ID, room[1].ID)
	if len(pathEdges) == 0 {
		t.Fatalf("no path found from %s to %s", room[0].ID, room[1].ID)
	}

	if len(pathEdges) != 1 {
		t.Fatalf("found %v edges between %s and %s, expected 1", len(pathEdges), room[0].ID, room[1].ID)
	}

	// make sure the edges are correct
	switch {
	case pathEdges[0].Src.Device.ID != room[0].ID:
		t.Fatalf("path was invalid at edge 0. srcDeviceID: %s, expected: %s\n", pathEdges[0].Src.Device.ID, room[0].ID)
	case pathEdges[0].SrcPort.Name != "":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].SrcPort.Name, "")
	case pathEdges[0].Dst.Device.ID != room[1].ID:
		t.Fatalf("path was invalid at edge 0. dstDeviceID: %s, expected: %s\n", pathEdges[0].Dst.Device.ID, room[1].ID)
	case pathEdges[0].DstPort.Name != "hdmi!2":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].DstPort.Name, "hdmi!2")
	}
}

func TestFindPathThroughOneSwitch(t *testing.T) {
	room := []api.Device{
		api.Device{
			ID: "ITB-1101-HDMI1",
			Ports: []api.Port{
				api.Port{
					Endpoints: []api.DeviceID{"ITB-1101-SW1"},
					Outgoing:  true,
					Type:      "video",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-HDMI2",
			Ports: []api.Port{
				api.Port{
					Endpoints: []api.DeviceID{"ITB-1101-SW1"},
					Outgoing:  true,
					Type:      "video",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-SW1",
			Ports: []api.Port{
				api.Port{
					Name:      "0",
					Endpoints: []api.DeviceID{"ITB-1101-HDMI1"},
					Incoming:  true,
					Type:      "video",
				},
				api.Port{
					Name:      "1",
					Endpoints: []api.DeviceID{"ITB-1101-HDMI2"},
					Incoming:  true,
					Type:      "video",
				},
				api.Port{
					Name:      "0",
					Endpoints: []api.DeviceID{"ITB-1101-D1"},
					Outgoing:  true,
					Type:      "video",
				},
				api.Port{
					Name:      "1",
					Endpoints: []api.DeviceID{"ITB-1101-D2"},
					Outgoing:  true,
					Type:      "video",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-D1",
			Ports: []api.Port{
				api.Port{
					Name:      "hdmi!2",
					Endpoints: []api.DeviceID{"ITB-1101-SW1"},
					Incoming:  true,
					Type:      "video",
				},
			},
		},
		api.Device{
			ID: "ITB-1101-D2",
			Ports: []api.Port{
				api.Port{
					Name:      "hdmi!2",
					Endpoints: []api.DeviceID{"ITB-1101-SW1"},
					Incoming:  true,
					Type:      "video",
				},
			},
		},
	}

	g := NewGraph(room, "video")
	shortest := path.DijkstraAllPaths(g)

	//
	//
	// HDMI1 -> D1
	pathEdges := PathFromTo(g, &shortest, room[0].ID, room[3].ID)
	if len(pathEdges) == 0 {
		t.Fatalf("no path found from %s to %s", room[0].ID, room[3].ID)
	}

	if len(pathEdges) != 2 {
		t.Fatalf("found %v edges between %s and %s, expected 2", len(pathEdges), room[0].ID, room[3].ID)
	}

	// check edge 0
	switch {
	case pathEdges[0].Src.Device.ID != room[0].ID:
		t.Fatalf("path was invalid at edge 0. srcDeviceID: %s, expected: %s\n", pathEdges[0].Src.Device.ID, room[0].ID)
	case pathEdges[0].SrcPort.Name != "":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].SrcPort.Name, "")
	case pathEdges[0].Dst.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 0. dstDeviceID: %s, expected: %s\n", pathEdges[0].Dst.Device.ID, room[2].ID)
	case pathEdges[0].DstPort.Name != "0":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].DstPort.Name, "0")
	}

	// check edge 1
	switch {
	case pathEdges[1].Src.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 1. srcDeviceID: %s, expected: %s\n", pathEdges[1].Src.Device.ID, room[0].ID)
	case pathEdges[1].SrcPort.Name != "0":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].SrcPort.Name, "0")
	case pathEdges[1].Dst.Device.ID != room[3].ID:
		t.Fatalf("path was invalid at edge 1. dstDeviceID: %s, expected: %s\n", pathEdges[1].Dst.Device.ID, room[3].ID)
	case pathEdges[1].DstPort.Name != "hdmi!2":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].DstPort.Name, "hdmi!2")
	}

	//
	//
	// HDMI2 -> D1
	pathEdges = PathFromTo(g, &shortest, room[1].ID, room[3].ID)
	if len(pathEdges) == 0 {
		t.Fatalf("no path found from %s to %s", room[1].ID, room[3].ID)
	}

	if len(pathEdges) != 2 {
		t.Fatalf("found %v edges between %s and %s, expected 2", len(pathEdges), room[1].ID, room[3].ID)
	}

	// check edge 0
	switch {
	case pathEdges[0].Src.Device.ID != room[1].ID:
		t.Fatalf("path was invalid at edge 0. srcDeviceID: %s, expected: %s\n", pathEdges[0].Src.Device.ID, room[1].ID)
	case pathEdges[0].SrcPort.Name != "":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].SrcPort.Name, "")
	case pathEdges[0].Dst.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 0. dstDeviceID: %s, expected: %s\n", pathEdges[0].Dst.Device.ID, room[2].ID)
	case pathEdges[0].DstPort.Name != "1":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].DstPort.Name, "1")
	}

	// check edge 1
	switch {
	case pathEdges[1].Src.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 1. srcDeviceID: %s, expected: %s\n", pathEdges[1].Src.Device.ID, room[0].ID)
	case pathEdges[1].SrcPort.Name != "0":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].SrcPort.Name, "0")
	case pathEdges[1].Dst.Device.ID != room[3].ID:
		t.Fatalf("path was invalid at edge 1. dstDeviceID: %s, expected: %s\n", pathEdges[1].Dst.Device.ID, room[3].ID)
	case pathEdges[1].DstPort.Name != "hdmi!2":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].DstPort.Name, "hdmi!2")
	}

	//
	//
	// HDMI1 -> D2
	pathEdges = PathFromTo(g, &shortest, room[0].ID, room[4].ID)
	if len(pathEdges) == 0 {
		t.Fatalf("no path found from %s to %s", room[0].ID, room[4].ID)
	}

	if len(pathEdges) != 2 {
		t.Fatalf("found %v edges between %s and %s, expected 2", len(pathEdges), room[0].ID, room[4].ID)
	}

	// check edge 0
	switch {
	case pathEdges[0].Src.Device.ID != room[0].ID:
		t.Fatalf("path was invalid at edge 0. srcDeviceID: %s, expected: %s\n", pathEdges[0].Src.Device.ID, room[0].ID)
	case pathEdges[0].SrcPort.Name != "":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].SrcPort.Name, "")
	case pathEdges[0].Dst.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 0. dstDeviceID: %s, expected: %s\n", pathEdges[0].Dst.Device.ID, room[2].ID)
	case pathEdges[0].DstPort.Name != "0":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].DstPort.Name, "0")
	}

	// check edge 1
	switch {
	case pathEdges[1].Src.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 1. srcDeviceID: %s, expected: %s\n", pathEdges[1].Src.Device.ID, room[0].ID)
	case pathEdges[1].SrcPort.Name != "1":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].SrcPort.Name, "1")
	case pathEdges[1].Dst.Device.ID != room[4].ID:
		t.Fatalf("path was invalid at edge 1. dstDeviceID: %s, expected: %s\n", pathEdges[1].Dst.Device.ID, room[4].ID)
	case pathEdges[1].DstPort.Name != "hdmi!2":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].DstPort.Name, "hdmi!2")
	}

	//
	//
	// HDMI2 -> D2
	pathEdges = PathFromTo(g, &shortest, room[1].ID, room[4].ID)
	if len(pathEdges) == 0 {
		t.Fatalf("no path found from %s to %s", room[1].ID, room[4].ID)
	}

	if len(pathEdges) != 2 {
		t.Fatalf("found %v edges between %s and %s, expected 2", len(pathEdges), room[1].ID, room[4].ID)
	}

	// check edge 0
	switch {
	case pathEdges[0].Src.Device.ID != room[1].ID:
		t.Fatalf("path was invalid at edge 0. srcDeviceID: %s, expected: %s\n", pathEdges[0].Src.Device.ID, room[1].ID)
	case pathEdges[0].SrcPort.Name != "":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].SrcPort.Name, "")
	case pathEdges[0].Dst.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 0. dstDeviceID: %s, expected: %s\n", pathEdges[0].Dst.Device.ID, room[2].ID)
	case pathEdges[0].DstPort.Name != "1":
		t.Fatalf("path was invalid at edge 0. srcPortName: %s, expected: %s\n", pathEdges[0].DstPort.Name, "1")
	}

	// check edge 1
	switch {
	case pathEdges[1].Src.Device.ID != room[2].ID:
		t.Fatalf("path was invalid at edge 1. srcDeviceID: %s, expected: %s\n", pathEdges[1].Src.Device.ID, room[0].ID)
	case pathEdges[1].SrcPort.Name != "1":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].SrcPort.Name, "1")
	case pathEdges[1].Dst.Device.ID != room[4].ID:
		t.Fatalf("path was invalid at edge 1. dstDeviceID: %s, expected: %s\n", pathEdges[1].Dst.Device.ID, room[4].ID)
	case pathEdges[1].DstPort.Name != "hdmi!2":
		t.Fatalf("path was invalid at edge 1. srcPortName: %s, expected: %s\n", pathEdges[1].DstPort.Name, "hdmi!2")
	}
}
