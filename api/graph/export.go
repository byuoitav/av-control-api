package graph

import (
	"bytes"

	"github.com/goccy/go-graphviz"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

func GraphToDot(g *simple.DirectedGraph) ([]byte, error) {
	return dot.Marshal(g, "", "", "")
}

func GraphToSvg(graph *simple.DirectedGraph) ([]byte, error) {
	b, err := dot.Marshal(graph, "", "", "")
	if err != nil {
		return nil, err
	}

	cGraph, err := graphviz.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	g := graphviz.New()

	var buf bytes.Buffer
	if err := g.Render(cGraph, graphviz.SVG, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
