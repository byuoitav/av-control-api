package inputgraph

import (
	"testing"

	"github.com/byuoitav/av-control-api/api/base"
	"github.com/byuoitav/common/log"
)

var i1 = base.Device{
	ID: "i1",
}
var i2 = base.Device{
	ID: "i2",
}
var i3 = base.Device{
	ID: "i3",
}
var i4 = base.Device{
	ID: "i1",
}
var i5 = base.Device{
	ID: "i5",
}
var i6 = base.Device{
	ID: "i6",
}

var a = base.Device{
	ID: "a",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "i1",
			ID:                "in1",
			DestinationDevice: "a",
		},
		base.Port{
			SourceDevice:      "i2",
			ID:                "in2",
			DestinationDevice: "a",
		},
		base.Port{
			SourceDevice:      "i2",
			ID:                "in2",
			DestinationDevice: "a",
		},
		base.Port{
			SourceDevice:      "i2",
			ID:                "in2",
			DestinationDevice: "a",
		},
		base.Port{
			SourceDevice:      "a",
			ID:                "out1",
			DestinationDevice: "c",
		},
	},
}

var b = base.Device{
	ID: "b",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "i3",
			ID:                "in1",
			DestinationDevice: "b",
		},
		base.Port{
			SourceDevice:      "i4",
			ID:                "in2",
			DestinationDevice: "b",
		},
		base.Port{
			SourceDevice:      "i5",
			ID:                "in3",
			DestinationDevice: "b",
		},
		base.Port{
			SourceDevice:      "b",
			ID:                "out1",
			DestinationDevice: "c",
		},
		base.Port{
			SourceDevice:      "b",
			ID:                "out2",
			DestinationDevice: "d",
		},
	},
}

var c = base.Device{
	ID: "c",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "a",
			ID:                "in1",
			DestinationDevice: "c",
		},
		base.Port{
			SourceDevice:      "b",
			ID:                "in2",
			DestinationDevice: "c",
		},
		base.Port{
			SourceDevice:      "c",
			ID:                "out1",
			DestinationDevice: "o1",
		},
		base.Port{
			SourceDevice:      "c",
			ID:                "out2",
			DestinationDevice: "o2",
		},
		base.Port{
			SourceDevice:      "c",
			ID:                "out3",
			DestinationDevice: "o3",
		},
	},
}
var d = base.Device{
	ID: "d",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "b",
			ID:                "in1",
			DestinationDevice: "d",
		},
		base.Port{
			SourceDevice:      "d",
			ID:                "out1",
			DestinationDevice: "o4",
		},
		base.Port{
			SourceDevice:      "d",
			ID:                "out2",
			DestinationDevice: "o5",
		},
	},
}
var o1 = base.Device{
	ID: "o1",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "c",
			ID:                "in1",
			DestinationDevice: "o1",
		},
	},
}
var o2 = base.Device{
	ID: "o2",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "c",
			ID:                "in1",
			DestinationDevice: "o2",
		},
	},
}
var o3 = base.Device{
	ID: "o3",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "c",
			ID:                "in1",
			DestinationDevice: "o3",
		},
	},
}
var o4 = base.Device{
	ID: "o4",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "d",
			ID:                "in1",
			DestinationDevice: "o4",
		},
	},
}
var o5 = base.Device{
	ID: "o5",
	Ports: []base.Port{
		base.Port{
			SourceDevice:      "d",
			ID:                "in1",
			DestinationDevice: "o5",
		},
		base.Port{
			SourceDevice:      "i6",
			ID:                "in2",
			DestinationDevice: "o5",
		},
	},
}

var Devices = []base.Device{a, b, c, d, i1, i2, i3, i4, i5, o1, o2, o3, o4, o5, i6}

func TestGraphBuilding(t *testing.T) {

	debug = false

	graph, err := BuildGraph(Devices, "video")
	if err != nil {
		log.L.Errorf("error: %v", err.Error())
		t.FailNow()
	}

	if debug {
		log.L.Errorf("%+v", graph.AdjacencyMap)
	}
}

func TestReachability(t *testing.T) {

	graph, err := BuildGraph(Devices, "video")
	if err != nil {
		log.L.Infof("error: %v", err.Error())
		t.FailNow()
	}

	debug = true
	ok, ret, _ := CheckReachability("o3", "i1", graph)
	if !ok {
		t.FailNow()
	}

	if debug {
		for _, v := range ret {
			log.L.Infof("%v", v.ID)
		}
	}
	debug = false

	ok, _, _ = CheckReachability("o5", "i1", graph)
	if ok {
		t.FailNow()
	}
	ok, _, _ = CheckReachability("o5", "i6", graph)
	if !ok {
		t.FailNow()
	}
	ok, _, _ = CheckReachability("o3", "i3", graph)
	if !ok {
		t.FailNow()
	}
}
