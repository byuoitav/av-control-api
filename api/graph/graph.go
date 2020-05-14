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

	// debugging
	fmt.Printf("\ngraph:\n")

	// print all nodes
	fmt.Printf("nodes:\n")
	_nodes := g.Nodes()
	for _nodes.Next() {
		gn := _nodes.Node().(Node)
		fmt.Printf("\t%s\t%v\n", gn.Device.ID, gn.ID())
	}

	// print all edges
	fmt.Printf("edges\n")
	_edges := g.Edges()
	for _edges.Next() {
		ge := _edges.Edge().(Edge)
		fmt.Printf("\t%v|%v -> %v|%v\n", ge.Src.Device.ID, ge.SrcPort.Name, ge.Dst.Device.ID, ge.DstPort.Name)
	}

	fmt.Printf("\n")

	return g
}
