package main

import (
	"github.com/alonsovidales/go_graph"
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func LabeledUndirected(gr *graphs.Graph) (graph.LabeledUndirected, []graph.LI) {

	g := make(graph.LabeledAdjacencyList, len(gr.RawEdges))
	labels := make([]graph.LI, len(gr.RawEdges))

	for i := len(gr.RawEdges); i >= 0; i-- {
		edge := gr.RawEdges[i]
		l := graph.LI(i)
		labels[i] = l
		g[i] = append(g[i], graph.Half{graph.NI(edge.To), l})
		g[edge.From] = append(g[edge.From], graph.Half{graph.NI(i), l})
	}

	return graph.LabeledUndirected{g}, labels
}

func PreparePrim(gr *graphs.Graph) (*graph.LabeledUndirected, []graph.LI) {
	/*
		g, labels := LabeledUndirected(gr)
		return &g, labels
	*/
	labels := make([]graph.LI, len((*gr).RawEdges))
	var g graph.LabeledUndirected

	for i := len(gr.RawEdges) - 1; i >= 0; i-- {
		edge := gr.RawEdges[i]
		l := graph.LI(i)
		labels[i] = l
		g.AddEdge(graph.Edge{graph.NI(edge.From), graph.NI(edge.To)}, l)
	}

	return &g, labels

}

func Prim(start uint64, g *graph.LabeledUndirected, labels []graph.LI, gr *graphs.Graph) (mst []graphs.Edge) {

	actualStart := graph.NI(start)

	// weight function
	w := func(arcLabel graph.LI) float64 { return (*gr).RawEdges[int(arcLabel)].Weight }

	// get connected components

	//a := g.LabeledAdjacencyList
	f := graph.NewFromList(len(g.LabeledAdjacencyList))

	var leaves bits.Bits
	_, _ = g.Prim(actualStart, w, &f, labels, &leaves)


	res, _ := f.LabeledUndirected(labels, nil)
	edges := make([]graphs.Edge, 0, len(labels))
	res.Edges(func(e graph.LabeledEdge) {
		edges = append(edges, graphs.Edge{uint64(e.Edge.N1), uint64(e.Edge.N2), w(e.LI)})
	})

	return edges
}
