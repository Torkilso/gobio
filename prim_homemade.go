package main

import (
	"container/heap"
	"fmt"
	"github.com/alonsovidales/go_graph"
	"math"
)

/*

func Prim3(start uint64, gr *graphs.Graph) []graphs.Edge {

	nodes := (*gr).Vertices

	dst := make([]graphs.Edge, 0, len(nodes)-1)

	if len(nodes) == 0 {
		return nil
	}

	q := &primQueue{
		indexOf: make(map[int]int, len(nodes)-1),
		nodes:   make([]graphs.Edge, 0, len(nodes)-1),
	}
	start = 0
	for u := range nodes {
		if u == start { continue }
		heap.Push(q, graphs.Edge{u, math.MaxUint64, math.Inf(1)})
	}
	u := start

	//heap.Push(q, graphs.Edge{u, u, 0})
	fmt.Println("Edges", (*gr).VertexEdges[u])
	fmt.Println("QUeue before", len(q.nodes), q.nodes)
	to := uint64(0)
	toValue := math.MaxFloat64
	for v, w := range (*gr).VertexEdges[u] {
		if toValue > w {
			to = v
			toValue = w
		}
	}
	dst = append(dst, graphs.Edge{u, to, toValue})


	for q.Len() > 0 {

		e := heap.Pop(q).(graphs.Edge)
		fmt.Println("Edge", e)
		if e.To != math.MaxUint64  {
			dst = append(dst, e)
		}else {
			fmt.Println("Edge", e)
			//panic(e)
		}

		u = e.From
		for n, w := range (*gr).VertexEdges[u] {
			fmt.Println("Checking from to", u, n)
			if key, ok := q.key(n); ok {
				if w < key {
					q.update(u, n, w)
				}
			}
		}
		dst = append(dst, e)

	}

	fmt.Println("Len", len(dst), dst)
	return dst
}
*/
func Prim3(start uint64, gr *graphs.Graph) []graphs.Edge {
	nodes := (*gr).Vertices

	start = 0
	// Choose a random start vertices

	result := make([]graphs.Edge, 0, len(nodes)-1)
	//resultTaken := make(map[uint64]bool)

	q := &primQueue{
		indexOf: make(map[int]int, len(nodes)-1),
		nodes:   make([]graphs.Edge, 0, len(nodes)-1),
	}

	// Add all first edges to queue
	for v, w := range (*gr).VertexEdges[start] {
		heap.Push(q, graphs.Edge{start, v, w})
	}

	// Take the best edge
	edge := heap.Pop(q).(graphs.Edge)

	result = append(result, edge)

	next := edge.To

	for v, w := range (*gr).VertexEdges[next] {
		heap.Push(q, graphs.Edge{next, v, w})
	}

	edge = heap.Pop(q).(graphs.Edge)

	result = append(result, edge)

	fmt.Println("Result", result)
	return result
}

// primQueue is a Prim's priority queue. The priority queue is a
// queue of edge From nodes keyed on the minimum edge weight to
// a node in the set of nodes already connected to the minimum
// spanning forest.
type primQueue struct {
	indexOf map[int]int
	nodes   []graphs.Edge
}

func (q *primQueue) Less(i, j int) bool {
	return q.nodes[i].Weight < q.nodes[j].Weight
}

func (q *primQueue) Swap(i, j int) {
	q.indexOf[int(q.nodes[i].From)] = j
	q.indexOf[int(q.nodes[j].From)] = i
	q.nodes[i], q.nodes[j] = q.nodes[j], q.nodes[i]
}

func (q *primQueue) Len() int {
	return len(q.nodes)
}

func (q *primQueue) Push(x interface{}) {
	n := x.(graphs.Edge)
	q.indexOf[int(n.From)] = len(q.nodes)
	q.nodes = append(q.nodes, n)

}

func (q *primQueue) Pop() interface{} {
	n := q.nodes[len(q.nodes)-1]
	q.nodes = q.nodes[:len(q.nodes)-1]
	delete(q.indexOf, int(n.From))
	fmt.Println("Popped", n.From)

	return n
}

// key returns the key for the node u and whether the node is
// in the queue. If the node is not in the queue, key is returned
// as +Inf.
func (q *primQueue) key(u uint64) (key float64, ok bool) {

	i, ok := q.indexOf[int(u)]
	fmt.Println("Key", u, q.nodes[i])
	if !ok {
		return math.Inf(1), false
	}
	//fmt.Println("Key", u, i)
	return q.nodes[i].Weight, ok
}

// update updates u's position in the queue with the new closest
// MST-connected neighbour, v, and the key weight between u and v.
func (q *primQueue) update(u, v uint64, key float64) {
	i, ok := q.indexOf[int(v)]
	if !ok {
		fmt.Println("Could not update", v, q.indexOf)
		return
	}
	fmt.Println( "FROM", q.nodes[i].From,"TO", v, "Weight", key, q.nodes[i] )
	q.nodes[i].To = v
	q.nodes[i].Weight = key
	fmt.Println("After",q.nodes[i] )
	heap.Fix(q, i)
}
