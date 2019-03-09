package main

import (
	"github.com/alonsovidales/go_graph"
	"log"
	"math"
	"math/rand"
)



func GraphToGeno(gr *graphs.Graph) []uint64 {
	l := len((*gr).RawEdges)
	geno := make([]uint64, l + 1)
	for _, edge := range (*gr).RawEdges {
		geno[edge.From] = edge.To
	}
	geno[l] = uint64(l) // Points to itself
	return geno
}

func GenoToGraph(img *Image, geno []uint64) *graphs.Graph {
	edges := make([]graphs.Edge, len(geno) - 1)
	for i := range edges {
		edges[i] = graphs.Edge{uint64(i), geno[i], Dist(img, i, int(geno[i]))}
	}
	return graphs.GetGraph(edges, true)
}
func GeneratePopulation(img *Image, n int) *Population {
	solutions := make([]Solution, n)

	imgAsGraph := GenerateGraph(img)

	log.Println("Prim graph start")
	primGraph, labels := PreparePrim(imgAsGraph)
	log.Println("Prim graph end")

	for i := 0 ; i < n ; i++ {
		start := rand.Intn(n)
		mstGraph := graphs.GetGraph(Prim(uint64(start), primGraph, labels, imgAsGraph), true)
		solutions[i] = SolutionFromGenotype(img, mstGraph)
	}
	return &Population{solutions}
}

func Tournament(p *Population, k int) int {

	bestIdx := -1
	bestCost := math.MaxFloat64
	for i := 0; i < k; i++ {
		idx := rand.Intn(len((*p).solutions))
		c := (*p).solutions[i].weightedSum()
		if c < bestCost {
			bestIdx = idx
			bestCost = c
		}
	}
	return bestIdx

}

func RunGeneration(img *Image, pop *Population) *Population {
	result := make([]Solution, len((*pop).solutions))

	for i := 0 ; i < len((*pop).solutions); i += 2 {
		p1Idx := Tournament(pop, 2)
		p2Idx := Tournament(pop, 2)

		p1 := (*pop).solutions[p1Idx]
		p2 := (*pop).solutions[p2Idx]

		result[i], result[i+1] = Crossover(img, &p1, &p2)

		if rand.Float32() < .1 {
			result[i] = Mutate(result[i].genotype, img)
		}
		if rand.Float32() < .1 {
			result[i+1] = Mutate(result[i+1].genotype, img)
		}

	}
	return &Population{result}
}

func Mutate(genotype []uint64, img *Image) Solution {
	for i := range genotype {
		if rand.Float32() < .2 {
			possibleValues := GetTargets(img, i)
			chosen := rand.Intn(len(possibleValues))
			genotype[i] = uint64(possibleValues[chosen])
		}
	}
	graph := GenoToGraph(img, genotype)
	return SolutionFromGenotype(img, graph)
}


func SolutionFromGenotype(img *Image, g *graphs.Graph) Solution {
	groups := g.ConnectedComponents()
	deviation := deviation(img, groups)
	connectivity := connectiviy(img, groups)
	crowdingDistance := 0.0
	return Solution{GraphToGeno(g), deviation, connectivity, crowdingDistance }

}
func Crossover(img *Image, parent1, parent2 *Solution) (Solution, Solution) {

	n := len((*parent1).genotype)

	offspring1 := make([]uint64, n)
	offspring2 := make([]uint64, n)

	for i := 0; i < n; i++ {
		// Update with 50% change
		if rand.Float32() < .5 {
			offspring1[i], offspring2[i] = (*parent2).genotype[i], (*parent1).genotype[i]
		}else {
			offspring1[i], offspring2[i] = (*parent1).genotype[i], (*parent2).genotype[i]
		}
	}
	graph1 := GenoToGraph(img, offspring1)
	graph2 := GenoToGraph(img, offspring2)

	return SolutionFromGenotype(img, graph1), SolutionFromGenotype(img, graph2)

}


func createPopulationFromParents(parents []*Solution) []*Solution {
	return parents
}
