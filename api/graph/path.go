package graph

import (
	"fmt"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
)

type Path []Edge

func (p Path) String() string {
	var b strings.Builder
	for _, edge := range p {
		b.WriteString(fmt.Sprintf("%s|%s -> ", edge.Src.Device.ID, edge.SrcPort.Name))
	}

	b.WriteString(fmt.Sprintf("%s|%s", p[len(p)-1].Dst.Device.ID, p[len(p)-1].DstPort.Name))
	return b.String()
}

func PathFromTo(g *simple.DirectedGraph, paths *path.AllShortest, src, dst api.DeviceID) Path {
	var path Path

	p, _, _ := paths.Between(NodeID(src), NodeID(dst))
	if len(path) == 0 {
		return nil
	}

	for i := 0; i < len(p)-1; i++ {
		edge := g.Edge(p[i].ID(), p[i+1].ID())
		if edge == nil {
			return nil
		}

		path = append(path, edge.(Edge))
	}

	return path
}

func PathToEnd(g *simple.DirectedGraph, src api.DeviceID) Path {
	var path Path

	start := g.Node(NodeID(src))
	if start == nil {
		return nil
	}

	search := traverse.DepthFirst{
		Traverse: func(edge graph.Edge) bool {
			path = append(path, edge.(Edge))
			return true
		},
	}

	search.Walk(g, start, func(node graph.Node) bool {
		return false
	})

	return path
}
