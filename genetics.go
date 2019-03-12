package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"math"
	"math/rand"
	"sort"
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

		result[i], result[i+1] = Crossover(img, p1, p2)

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

func GenoToGraph(img *Image, geno []uint64) *graphs.Graph {
	edges := make([]graphs.Edge, len(geno))
	for i := range geno {
		edges[i] = graphs.Edge{uint64(i), geno[i], Dist(img, i, int(geno[i]))}
	}

	return graphs.GetGraph(edges, true)
}

func GeneratePopulation(img *Image, n int) []*Solution {
	solutions := make([]*Solution, n)

	startT := time.Now()
	imgAsGraph := GenerateGraph(img)

	width := len(*img)
	height := len((*img)[0])

	primGraph, labels := PreparePrim(imgAsGraph)

	fmt.Println("Done generation setup", time.Now().Sub(startT))
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < n; i++ {
		start := r1.Intn(width * height)
		mst2 := Prim(uint64(start), primGraph, labels, imgAsGraph)

		mstGraph := graphs.GetGraph(mst2, true)

		solutions[i] = SolutionFromGenotypeNSGA(img, mstGraph)

		//visualizeImageGraph("sol.png", img, mstGraph)
	}
	//visualizeImageGraph("sol_act.png", img, GenoToGraph(img, solutions[1].genotype))

	return solutions
}

// runGeneration for NSGA
func createPopulationFromParents(img *Image, pop []*Solution) []*Solution {
	result := make([]*Solution, len(pop))

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < len(pop); i += 2 {
		p1Idx := r1.Intn(len(pop))
		p2Idx := r1.Intn(len(pop))

		p1 := pop[p1Idx]
		p2 := pop[p2Idx]

		result[i], result[i+1] = Crossover(img, p1, p2)

		/*if r1.Float32() < .1 {
			result[i] = Mutate(result[i].genotype, img)
		}

		if r1.Float32() < .1 {
			result[i+1] = Mutate(result[i+1].genotype, img)
		}*/
	}

	return result
}

func Mutate(genotype []uint64, img *Image) *Solution {
	for i := range genotype {
		if rand.Float32() < .2 {
			possibleValues := GetTargets(img, i)
			chosen := rand.Intn(len(possibleValues))
			genotype[i] = uint64(possibleValues[chosen])
		}
	}
	graph := GenoToGraph(img, genotype)
	return SolutionFromGenotypeNSGA(img, graph)
}

func Crossover(img *Image, parent1, parent2 *Solution) (*Solution, *Solution) {

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

	graph1 := GenoToGraph(img, offspring1)
	graph2 := GenoToGraph(img, offspring2)

	return SolutionFromGenotypeNSGA(img, graph1), SolutionFromGenotypeNSGA(img, graph2)
}

func SolutionFromGenotypeNSGA(img *Image, g *graphs.Graph) *Solution {
	groups := g.ConnectedComponents()

	deviation := deviation(img, groups)
	connectivity := connectiviy(img, groups)
	crowdingDistance := 0.0
	//visualizeImageGraph("graph.png", img, g)
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