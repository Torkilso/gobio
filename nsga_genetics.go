package main

import (
	"github.com/alonsovidales/go_graph"
	"log"
	"math"
	"math/rand"
)

func GeneratePopulationNSGA(img *Image, n int) []*Solution {
	solutions := make([]*Solution, n)

	imgAsGraph := GenerateGraph(img)

	log.Println("Prim graph start")
	primGraph, labels := PreparePrim(imgAsGraph)
	log.Println("Prim graph end")

	for i := 0; i < n; i++ {
		start := rand.Intn(n)
		mstGraph := graphs.GetGraph(Prim(uint64(start), primGraph, labels, imgAsGraph), true)
		solutions[i] = SolutionFromGenotypeNSGA(img, mstGraph)
	}
	return solutions
}

// runGeneration for NSGA
func createPopulationFromParents(img *Image, pop []*Solution) []*Solution {
	result := make([]*Solution, len(pop))

	for i := 0; i < len(pop); i += 2 {
		p1Idx := TournamentNSGA(pop, 2)
		p2Idx := TournamentNSGA(pop, 2)

		p1 := pop[p1Idx]
		p2 := pop[p2Idx]

		result[i], result[i+1] = CrossoverNSGA(img, p1, p2)

		if rand.Float32() < .1 {
			result[i] = MutateNSGA(result[i].genotype, img)
		}

		if rand.Float32() < .1 {
			result[i+1] = MutateNSGA(result[i+1].genotype, img)
		}
	}

	return result
}

func TournamentNSGA(p []*Solution, k int) int {
	bestIdx := -1
	bestCost := math.MaxFloat64

	for i := 0; i < k; i++ {
		idx := rand.Intn(len(p))

		c := p[i].weightedSum()

		if c < bestCost {
			bestIdx = idx
			bestCost = c
		}
	}

	return bestIdx
}

func MutateNSGA(genotype []uint64, img *Image) *Solution {
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

func CrossoverNSGA(img *Image, parent1, parent2 *Solution) (*Solution, *Solution) {

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

	return &Solution{GraphToGeno(g), deviation, connectivity, crowdingDistance}
}
