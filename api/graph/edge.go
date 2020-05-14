package graph

import (
	"github.com/byuoitav/av-control-api/api"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

type Edge struct {
	Src     Node
	Dst     Node
	SrcPort *api.Port
	DstPort *api.Port
}

func (e Edge) From() graph.Node { return e.Src }

func (e Edge) To() graph.Node { return e.Dst }

func (e Edge) ReversedEdge() graph.Edge {
	return Edge{
		Src:     e.Dst,
		SrcPort: e.DstPort,
		Dst:     e.Src,
		DstPort: e.SrcPort,
	}
}

func (e Edge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		encoding.Attribute{Key: "taillabel", Value: e.SrcPort.Name},
		encoding.Attribute{Key: "headlabel", Value: e.DstPort.Name},
	}
}
