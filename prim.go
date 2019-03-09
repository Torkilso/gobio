package main

import (
	"github.com/alonsovidales/go_graph"
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// byWeight Used to sort the graph edges by weight
type byWeight []graphs.Edge

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight < a[j].Weight }

func PreparePrim(gr *graphs.Graph) (*graph.LabeledUndirected, []graph.LI) {
	labels := make([]graph.LI, len((*gr).RawEdges))
	var g graph.LabeledUndirected

	for i, edge := range (*gr).RawEdges {
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

	a := g.LabeledAdjacencyList
	f := graph.NewFromList(len(a))

	var leaves bits.Bits
	_, _ = g.Prim(actualStart, w, &f, labels, &leaves)

	res, _ := f.LabeledUndirected(labels, nil)
	edges := make([]graphs.Edge, 0)
	res.Edges(func(e graph.LabeledEdge) {
		edges = append(edges, graphs.Edge{uint64(e.Edge.N1), uint64(e.Edge.N2), float64(1)})
	})

	return edges
}
