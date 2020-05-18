package graph

import (
	"fmt"

	"github.com/byuoitav/av-control-api/api"
	"gonum.org/v1/gonum/graph/simple"
)

func NewGraph(devices []api.Device, portType string) *simple.DirectedGraph {
	var nodes []Node

	for d := range devices {
		nodes = append(nodes, Node{
			id:     NodeID(devices[d].ID),
			Device: &devices[d],
		})
	}

	// create the graph
	g := simple.NewDirectedGraph()

	for src := range nodes {
		for srcP := range nodes[src].Ports {
			if !nodes[src].Ports[srcP].Outgoing || nodes[src].Ports[srcP].Type != portType {
				continue
			}

			// find the endpoint of this port
			for dst := range nodes {
				for dstP := range nodes[dst].Ports {
					if !nodes[dst].Ports[dstP].Incoming || nodes[dst].Ports[dstP].Type != portType {
						continue
					}

					// make sure they are both pointing to eachother
					if nodes[src].Ports[srcP].Endpoint != nodes[dst].Device.ID || nodes[dst].Ports[dstP].Endpoint != nodes[src].Device.ID {
						continue
					}

					g.SetEdge(Edge{
						Src:     nodes[src],
						Dst:     nodes[dst],
						SrcPort: &nodes[src].Ports[srcP],
						DstPort: &nodes[dst].Ports[dstP],
					})
				}
			}
		}
	}

	return g
}

func Transpose(g *simple.DirectedGraph) *simple.DirectedGraph {
	t := simple.NewDirectedGraph()

	edges := g.Edges()
	for edges.Next() {
		t.SetEdge(edges.Edge().ReversedEdge())
	}

	return t
}

func printGraph(g *simple.DirectedGraph) {
	fmt.Printf("\ngraph:\n")

	fmt.Printf("nodes:\n")
	nodes := g.Nodes()
	for nodes.Next() {
		gn := nodes.Node().(Node)
		fmt.Printf("\t%s\t%v\n", gn.Device.ID, gn.ID())
	}

	fmt.Printf("edges\n")
	edges := g.Edges()
	for edges.Next() {
		ge := edges.Edge().(Edge)
		fmt.Printf("\t%v|%v -> %v|%v\n", ge.Src.Device.ID, ge.SrcPort.Name, ge.Dst.Device.ID, ge.DstPort.Name)
	}

	fmt.Printf("\n")
}
