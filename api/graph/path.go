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
	for i := range p {
		switch i {
		case 0:
			b.WriteString(fmt.Sprintf("%s|%s -> ", p[i].Src.Device.ID, p[i].SrcPort.Name))
		case len(p) - 1:
			b.WriteString(fmt.Sprintf("%s|%s|%s ->", p[i-1].DstPort.Name, p[i].Src.Device.ID, p[i].SrcPort.Name))
			b.WriteString(fmt.Sprintf("%s|%s", p[i].DstPort.Name, p[i].Dst.Device.ID))
		default:
			b.WriteString(fmt.Sprintf("%s|%s|%s ->", p[i-1].DstPort.Name, p[i].Src.Device.ID, p[i].SrcPort.Name))
		}
	}

	return b.String()
}

func PathFromTo(g *simple.DirectedGraph, paths *path.AllShortest, src, dst api.DeviceID) Path {
	var path Path

	p, _, _ := paths.Between(NodeID(src), NodeID(dst))
	if len(p) == 0 {
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

func LeavesFrom(g *simple.DirectedGraph, src api.DeviceID) []Node {
	var leaves []Node

	start := g.Node(NodeID(src))
	if start == nil {
		return nil
	}

	search := traverse.DepthFirst{
		Visit: func(node graph.Node) {
			if g.From(node.ID()).Len() == 0 {
				leaves = append(leaves, node.(Node))
			}
		},
	}

	search.Walk(g, start, func(graph.Node) bool {
		return false
	})

	return leaves
}

func Leaves(g *simple.DirectedGraph) []Node {
	var leaves []Node

	search := traverse.DepthFirst{
		Visit: func(node graph.Node) {
			if g.From(node.ID()).Len() == 0 {
				leaves = append(leaves, node.(Node))
			}
		},
	}

	nodes := g.Nodes()
	for nodes.Next() {
		search.Walk(g, nodes.Node(), func(graph.Node) bool {
			return false
		})
	}

	return leaves
}
