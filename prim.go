package main

import (
	"container/heap"
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

func Prim(start uint64, g *graph.LabeledUndirected, labels []graph.LI, gr *graphs.Graph, size int) (geno []uint64) {

	actualStart := graph.NI(start)

	// weight function
	w := func(arcLabel graph.LI) float64 { return (*gr).RawEdges[int(arcLabel)].Weight }

	// get connected components

	//a := g.LabeledAdjacencyList
	f := graph.NewFromList(len(g.LabeledAdjacencyList))

	var leaves bits.Bits
	_, _ = g.Prim(actualStart, w, &f, labels, &leaves)

	geno = make([]uint64, size)
	for i := 0; i < size; i++ {
		if int(start) == i {
			geno[i] = uint64(i)
			continue
		}
		geno[i] = uint64(f.Paths[i].From)

	}
	geno[start] = uint64(start)

	return geno
}

func Prim2(g *graph.LabeledUndirected, start graph.NI, w graph.WeightFunc, f *graph.FromList, labels []graph.LI, componentLeaves *bits.Bits) (numSpanned int, dist float64, geno []uint64) {
	al := g.LabeledAdjacencyList
	geno = make([]uint64, len(labels))
	if len(f.Paths) != len(al) {
		*f = graph.NewFromList(len(al))
	}
	if f.Leaves.Num != len(al) {
		f.Leaves = bits.New(len(al))
	}
	b := make([]prNode, len(al)) // "best"
	for n := range b {
		b[n].nx = graph.NI(n)
		b[n].fx = -1
	}
	rp := f.Paths
	var frontier prHeap
	rp[start] = graph.PathEnd{From: -1, Len: 1}
	numSpanned = 1
	fLeaves := &f.Leaves
	fLeaves.SetBit(int(start), 1)
	if componentLeaves != nil {
		if componentLeaves.Num != len(al) {
			*componentLeaves = bits.New(len(al))
		}
		componentLeaves.SetBit(int(start), 1)
	}
	for a := start; ; {
		for _, nb := range al[a] {
			if rp[nb.To].Len > 0 {
				continue // already in MST, no action
			}
			switch bp := &b[nb.To]; {
			case bp.fx == -1: // new node for frontier
				bp.from = fromHalf{From: a, Label: nb.Label}
				bp.wt = w(nb.Label)
				heap.Push(&frontier, bp)
			case w(nb.Label) < bp.wt: // better arc
				bp.from = fromHalf{From: a, Label: nb.Label}
				bp.wt = w(nb.Label)

				heap.Fix(&frontier, bp.fx)
			}
		}
		if len(frontier) == 0 {
			break // done
		}
		bp := heap.Pop(&frontier).(*prNode)
		a = bp.nx
		geno[int(a)] = uint64(bp.from.From)

		rp[a].Len = rp[bp.from.From].Len + 1
		rp[a].From = bp.from.From
		if len(labels) != 0 {
			labels[a] = bp.from.Label
		}
		dist += bp.wt
		fLeaves.SetBit(int(bp.from.From), 0)
		fLeaves.SetBit(int(a), 1)
		if componentLeaves != nil {
			componentLeaves.SetBit(int(bp.from.From), 0)
			componentLeaves.SetBit(int(a), 1)
		}
		numSpanned++
	}
	return
}

type prNode struct {
	nx   graph.NI
	from fromHalf
	wt   float64 // p.Weight(from.Label)
	fx   int
}

type prHeap []*prNode

func (h prHeap) Len() int           { return len(h) }
func (h prHeap) Less(i, j int) bool { return h[i].wt < h[j].wt }
func (h prHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].fx = i
	h[j].fx = j
}
func (p *prHeap) Push(x interface{}) {
	nd := x.(*prNode)
	nd.fx = len(*p)
	*p = append(*p, nd)
}
func (p *prHeap) Pop() interface{} {
	r := *p
	last := len(r) - 1
	*p = r[:last]
	return r[last]
}

type fromHalf struct {
	From  graph.NI
	Label graph.LI
}
