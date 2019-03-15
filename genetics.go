package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func Tournament(solutions []*Solution, k int) int {

	bestIdx := -1
	bestCost := math.MaxFloat64
	for i := 0; i < k; i++ {
		idx := rand.Intn(len(solutions))
		c := solutions[i].weightedSum()
		if c < bestCost {
			bestIdx = idx
			bestCost = c
		}
	}
	return bestIdx
}

func RunGeneration(img *Image, solutions []*Solution) []*Solution {
	result := make([]*Solution, len(solutions))

	for i := 0; i < len(solutions); i += 2 {
		p1Idx := Tournament(solutions, 2)
		p2Idx := Tournament(solutions, 2)

		p1 := solutions[p1Idx]
		p2 := solutions[p2Idx]

		result[i], result[i+1] = crossover(img, p1, p2)

		/* DONT MUTATE YET
		if rand.Float32() < .1 {
			result[i] = Mutate(result[i].genotype, img)
		}
		if rand.Float32() < .1 {
			result[i+1] = Mutate(result[i+1].genotype, img)
		}*/
	}

	return result
}

func (gr *Graph) ConnectedComponents() (groups []map[uint64]bool) {
	var groupToUse uint64
	usedVertex := make(map[uint64]uint64)
	currentGroup := uint64(0)
	for v := range gr.VertexEdges {
		if _, used := usedVertex[v]; !used {
			group := make(map[uint64]bool)
			gr.dfs(v, group, nil, nil)
			found := false
		groupSearch:
			for k := range group {
				if g, used := usedVertex[k]; used {
					groupToUse = g
					found = true

					break groupSearch
				}
			}
			if !found {
				groupToUse = currentGroup
				currentGroup++
				groups = append(groups, make(map[uint64]bool))
			}

			for k := range group {
				usedVertex[k] = groupToUse
				groups[groupToUse][k] = true
			}
		}
	}

	return
}
func dfs(origin int, usedVertex []bool, geno[]uint64) {
	usedVertex[origin] = true

	if !usedVertex[origin] {
		dfs(int(geno[origin]), usedVertex, geno)
	}


}
func GenoToConnectedComponents(geno []uint64) []map[uint64]bool {

	reverseGeno := make([]uint64, len(geno))
	for i, g := range geno {
		reverseGeno[geno[i]] = g
	}
	segments := make([]map[uint64]bool, 0)
	used := make([]bool, len(geno))
	segmentsMap := make(map[uint64]int, len(geno))
	for from, to := range geno {
		if !used[from] {
			group := make([]bool, 0)
			current := from
			for used[]

		}
	}
}

func GraphToGeno(gr *graphs.Graph, size int) []uint64 {
	geno := make([]uint64, size)

	edgesForNode := make(map[uint64]map[uint64]bool)

	for _, edge := range gr.RawEdges {
		if edgesForNode[edge.From] == nil {
			edgesForNode[edge.From] = make(map[uint64]bool)
		}
		if edgesForNode[edge.To] == nil {
			edgesForNode[edge.To] = make(map[uint64]bool)
		}
		edgesForNode[edge.From][edge.To] = true
		edgesForNode[edge.To][edge.From] = true
	}

	assigned := 0

	for assigned < size-1 {
		noneFound := true
		for id, val := range edgesForNode {
			if len(val) == 1 {
				for key := range edgesForNode[id] {
					geno[id] = key
					noneFound = false
					assigned++

					delete(edgesForNode[key], id)
					delete(edgesForNode, id)
				}
			}
		}

		if noneFound {
			for id := range edgesForNode {
				if len(edgesForNode[id]) == 0 {
					geno[id] = id
					assigned++

					delete(edgesForNode, id)
				} else {
					for key := range edgesForNode[id] {
						geno[id] = key
						assigned++

						delete(edgesForNode[key], id)
						delete(edgesForNode, id)
					}
				}
			}
		}
	}

	for lastKey := range edgesForNode {
		geno[lastKey] = lastKey
	}

	return geno
}

func GetGraph(size int, edges []graphs.Edge, undirected bool) (ug *graphs.Graph) {
	var weight float64

	ug = &graphs.Graph{
		RawEdges:    edges,
		Vertices:    make(map[uint64]bool, size),
		VertexEdges: make(map[uint64]map[uint64]float64, size),
		Undirected:  undirected,
		NegEdges:    false,
	}

	for _, edge := range edges {
		weight = edge.Weight

		ug.Vertices[edge.From] = true
		ug.Vertices[edge.To] = true
		if _, ok := ug.VertexEdges[edge.From]; ok {
			ug.VertexEdges[edge.From][edge.To] = weight
		} else {
			ug.VertexEdges[edge.From] = map[uint64]float64{edge.To: weight}
		}
		if undirected {
			if _, ok := ug.VertexEdges[edge.To]; ok {
				ug.VertexEdges[edge.To][edge.From] = weight
			} else {
				ug.VertexEdges[edge.To] = map[uint64]float64{edge.From: weight}
			}
		}
	}

	return
}

func GenoToGraph(img *Image, geno []uint64, weight bool) *graphs.Graph {
	edges := make([]graphs.Edge, len(geno))
	var w float64
	for i := range geno {
		if weight {
			w = Dist(img, i, int(geno[i]))
		}
		edges[i] = graphs.Edge{uint64(i), geno[i], w}
	}

	return GetGraph(ImageSize(img), edges, true)
}

func generatePopulation(img *Image, n int) Population {
	solutions := make([]*Solution, 0, n)

	startT := time.Now()
	imgAsGraph := GenerateGraph(img)

	width := len(*img)
	height := len((*img)[0])

	primGraph, labels := PreparePrim(imgAsGraph)

	fmt.Println("Done generation setup", time.Now().Sub(startT))
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	startT = time.Now()

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func(index int) {
			start := r1.Intn(width * height)
			mst := Prim(uint64(start), primGraph, labels, imgAsGraph)

			mstGraph := graphs.GetGraph(mst, true)

			channel <- SolutionFromGenotypeNSGA(img, mstGraph)
			defer wg.Done()
		}(i)
	}
	go func() {
		for t := range channel {
			solutions = append(solutions, t)
			wg.Done()
		}
	}()
	wg.Wait()
	fmt.Println("Done generation creation", time.Now().Sub(startT))

	return solutions
}

func (p *Population) evolve(img *Image) {
	size := len(*p)
	result := make([]*Solution, 0, size)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(size + size/2)

	for i := 0; i < size; i += 2 {
		go func(index int) {
			p1Idx := r1.Intn(size)
			p2Idx := r1.Intn(size)

			p1 := (*p)[p1Idx]
			p2 := (*p)[p2Idx]

			leftChild, rightChild := crossover(img, p1, p2)

			if r1.Float32() < .1 {
				leftChild.mutate(img)
			}

			if r1.Float32() < .1 {
				rightChild.mutate(img)
			}

			channel <- leftChild
			channel <- rightChild
			wg.Done()
		}(i)
	}

	go func() {
		for t := range channel {
			result = append(result, t)
			wg.Done()
		}
	}()

	wg.Wait()
	close(channel)

	*p = append(*p, result...)
}

func (s *Solution) mutate(img *Image) {

	index := rand.Intn(len(s.genotype))
	possibleValues := GetTargets(img, index)
	chosen := rand.Intn(len(possibleValues))

	s.genotype[uint64(index)] = uint64(possibleValues[chosen])

	graph := GenoToGraph(img, s.genotype, false)
	groups := graph.ConnectedComponents()

	s.connectivity = connectivity(img, groups)
	s.deviation = deviation(img, groups)
	s.crowdingDistance = 0.0
}

func crossover(img *Image, parent1, parent2 *Solution) (*Solution, *Solution) {

	s := time.Now()
	st := time.Now()
	n := len((*parent1).genotype)

	offspring1 := make([]uint64, n)
	offspring2 := make([]uint64, n)

	for i := 0; i < n; i++ {
		// Update with 50% change
		if rand.Float32() < .5 {
			offspring1[i], offspring2[i] = (*parent2).genotype[i], (*parent1).genotype[i]
		} else {
			offspring1[i], offspring2[i] = (*parent1).genotype[i], (*parent2).genotype[i]
		}
	}
	fmt.Println("Time for creating offsprings", time.Now().Sub(st).String())
	st = time.Now()
	graph1 := GenoToGraph(img, offspring1, false)
	graph2 := GenoToGraph(img, offspring2, false)
	fmt.Println("Time for geno to graph", time.Now().Sub(st).String())
	st = time.Now()

	groups1 := graph1.ConnectedComponents()
	groups2 := graph2.ConnectedComponents()
	fmt.Println("Time for connected components", time.Now().Sub(st).String())
	st = time.Now()

	s1 := &Solution{
		offspring1, deviation(img, groups1), connectivity(img, groups1), 0.0,
	}

	s2 := &Solution{
		offspring2, deviation(img, groups2), connectivity(img, groups2), 0.0,
	}
	fmt.Println("Time for solutions components", time.Now().Sub(st).String())

	fmt.Println("Time for crossover", time.Now().Sub(s).String())
	return s1, s2
}

func SolutionFromGenotypeNSGA(img *Image, g *graphs.Graph) *Solution {
	groups := g.ConnectedComponents()

	deviation := deviation(img, groups)
	connectivity := connectivity(img, groups)
	crowdingDistance := 0.0
	sol := &Solution{GraphToGeno(g, ImageSize(img)), deviation, connectivity, crowdingDistance}

	return sol
}

type GenoVertices struct {
	edges int
	value uint64
}

func GraphToGeno2(gr *graphs.Graph, size int) []uint64 {
	geno := make([]uint64, size)

	edgesForNode := make([]int, len(gr.Vertices))

	vertices := make([]GenoVertices, len(gr.Vertices))

	for _, edge := range gr.RawEdges {
		edgesForNode[edge.From]++
		edgesForNode[edge.To]++
	}
	for i := range gr.Vertices {
		vertices[i] = GenoVertices{edgesForNode[i], i}
	}

	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].edges < vertices[j].edges
	})

	for _, v := range vertices {
		for to := range gr.VertexEdges[v.value] {
			geno[v.value] = to
			if len(gr.VertexEdges[to]) > 1 {
				delete(gr.VertexEdges[to], v.value)
			}
			break
		}
	}

	return geno
}
