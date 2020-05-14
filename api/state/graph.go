package state

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/byuoitav/av-control-api/api"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type graphNode struct {
	id int64
	*api.Device
}

func (g graphNode) ID() int64 {
	return g.id
}

func (g graphNode) DOTID() string {
	return string(g.Device.ID)
}

type graphEdge struct {
	Src     graphNode
	Dst     graphNode
	SrcPort *api.Port
	DstPort *api.Port
}

func (e graphEdge) From() graph.Node { return e.Src }

func (e graphEdge) To() graph.Node { return e.Dst }

func (e graphEdge) ReversedEdge() graph.Edge {
	return graphEdge{
		Src:     e.Dst,
		SrcPort: e.DstPort,
		Dst:     e.Src,
		DstPort: e.SrcPort,
	}
}

func (g graphEdge) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		encoding.Attribute{Key: "taillabel", Value: g.SrcPort.Name},
		encoding.Attribute{Key: "headlabel", Value: g.DstPort.Name},
	}
}

type graphPath []graphEdge

func (p graphPath) String() string {
	var b strings.Builder
	for _, edge := range p {
		b.WriteString(fmt.Sprintf("%s|%s -> ", edge.Src.Device.ID, edge.SrcPort.Name))
	}

	b.WriteString(fmt.Sprintf("%s|%s", p[len(p)-1].Dst.Device.ID, p[len(p)-1].DstPort.Name))
	return b.String()
}

func graphNodeID(id string) int64 {
	sum := sha1.Sum([]byte(id))

	var i big.Int
	i.SetBytes(sum[:])
	return i.Int64()
}

func newDeviceGraph(devices []api.Device, portType string) *simple.DirectedGraph {
	var nodes []graphNode

	for d := range devices {
		nodes = append(nodes, graphNode{
			id:     graphNodeID(string(devices[d].ID)),
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

					g.SetEdge(graphEdge{
						Src:     nodes[src],
						Dst:     nodes[dst],
						SrcPort: &nodes[src].Ports[srcP],
						DstPort: &nodes[dst].Ports[dstP],
					})
				}
			}
		}
	}

	/*
		// debugging
		fmt.Printf("\ngraph:\n")

		// print all nodes
		fmt.Printf("nodes:\n")
		_nodes := g.Nodes()
		for _nodes.Next() {
			gn := _nodes.Node().(graphNode)
			fmt.Printf("\t%s\t%v\n", gn.Device.ID, gn.ID())
		}

		// print all edges
		fmt.Printf("edges\n")
		_edges := g.Edges()
		for _edges.Next() {
			ge := _edges.Edge().(graphEdge)
			fmt.Printf("\t%v|%v -> %v|%v\n", ge.Src.Device.ID, ge.SrcPort.Name, ge.Dst.Device.ID, ge.DstPort.Name)
		}

		fmt.Printf("\n")
	*/

	return g
}

func edgesBetween(g *simple.DirectedGraph, paths *path.AllShortest, src, dst api.DeviceID) graphPath {
	path, _, _ := paths.Between(graphNodeID(string(src)), graphNodeID(string(dst)))
	if len(path) == 0 {
		return nil
	}

	var edges []graphEdge

	for i := 0; i < len(path)-1; i++ {
		edge := g.Edge(path[i].ID(), path[i+1].ID())
		if edge == nil {
			return nil
		}

		edges = append(edges, edge.(graphEdge))
	}

	return edges
}

func exportGraphDot(g *simple.DirectedGraph, filePath string) error {
	b, err := dot.Marshal(g, "", "", "")
	if err != nil {
		return fmt.Errorf("unable to marshal graph: %w", err)
	}

	return ioutil.WriteFile(filePath, b, 0644)
}
