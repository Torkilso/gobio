package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"math"
	"sort"
	"time"
)

type SearchHelper struct {
	dominates         []int
	dominatedByAmount int
}

func fastNonDominatedSort(population []*Solution) map[int][]int {

	fronts := make(map[int][]int)
	searchHelperMap := make(map[int]*SearchHelper)

	for i, solution := range population {
		searchHelper := SearchHelper{
			dominatedByAmount: 0,
		}

		searchHelperMap[i] = &searchHelper

		for j, opponent := range population {
			if i == j {
				continue
			}

			if solution.dominate(opponent) {
				searchHelper.dominates = append(searchHelperMap[i].dominates, j)
			} else if opponent.dominate(solution) {
				searchHelper.dominatedByAmount++
			}
		}

		if searchHelper.dominatedByAmount == 0 {
			fronts[0] = append(fronts[0], i)
		}
	}

	frontRank := 0

	for len(fronts[frontRank]) != 0 {
		newFront := make([]int, 0)

		for _, frontSolution := range fronts[frontRank] {
			for _, solution := range searchHelperMap[frontSolution].dominates {
				searchHelperMap[solution].dominatedByAmount--
				if searchHelperMap[solution].dominatedByAmount == 0 {
					newFront = append(newFront, solution)
				}
			}
			population[frontSolution].frontNumber = frontRank
		}

		frontRank++
		if len(newFront) > 0 {
			fronts[frontRank] = newFront
		}
	}

	return fronts
}

func crowdingDistanceAssignment(ids []int, population []*Solution) {
	size := len(ids)

	// for deviation
	sort.Slice(ids, func(i, j int) bool {
		return population[ids[i]].deviation < population[ids[j]].deviation
	})

	population[ids[0]].crowdingDistance = math.Inf(1)
	population[ids[size-1]].crowdingDistance = math.Inf(1)

	for i := 1; i < size-1; i++ {
		population[ids[i]].crowdingDistance = (population[ids[i+1]].deviation - population[ids[i-1]].deviation) / (maxDeviation - minDeviation)
	}

	// for connectivity
	sort.Slice(ids, func(i, j int) bool {
		return population[ids[i]].connectivity < population[ids[j]].connectivity
	})

	for i := 1; i < size-1; i++ {
		population[ids[i]].crowdingDistance = population[ids[i]].crowdingDistance + (population[ids[i+1]].connectivity-population[ids[i-1]].connectivity)/(maxConnectivity-minConnectivity)
	}
}

func (p *Population) sortAndSelectParetoSolutions(populationSize, generation int, plotter *plot.Plot) {
	fronts := fastNonDominatedSort(*p)
	fmt.Println("\nGeneration:", generation, "Best before:", BestSolution(*p).weightedSum(), "Num fronts:", len(fronts))
	addParetoFrontToPlotter(plotter, *p, fronts, generation)

	newParents := make([]*Solution, 0)
	i := 0

	for len(newParents)+len(fronts[i]) < populationSize {
		crowdingDistanceAssignment(fronts[i], *p)
		frontSolutions := make([]*Solution, len(fronts[i]))

		for i, id := range fronts[i] {
			frontSolutions[i] = (*p)[id]
		}

		newParents = append(newParents, frontSolutions...)
		i++
	}

	lastFrontier := make([]*Solution, 0)

	if len(fronts[i]) > 0 {
		crowdingDistanceAssignment(fronts[i], *p)
		sort.Slice(fronts[i], func(j, k int) bool {
			return (*p)[fronts[i][j]].crowdingDistance > (*p)[fronts[i][k]].crowdingDistance
		})

		for _, id := range fronts[i] {
			if len(lastFrontier)+len(newParents) < populationSize {
				lastFrontier = append(lastFrontier, (*p)[id])
			} else {
				break
			}
		}
	}

	*p = append(newParents, lastFrontier...)
}

func nsgaII(image *Image, generations, populationSize int) []*Solution {
	start := time.Now()

	fmt.Println("Initiating NSGAII")
	fmt.Println("Generating", populationSize, "solutions")

	population := generatePopulation(image, populationSize)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds to generate solutions")

	start = time.Now()

	p := createParetoPlotter()
	population.sortAndSelectParetoSolutions(populationSize, 0, p)

	for t := 0; t < generations; t++ {
		startGeneration := time.Now()

		fmt.Println("Evolving")
		population.evolveWithTournament(image)
		fmt.Println("sortAndSelectParetoSolutions")

		population.sortAndSelectParetoSolutions(populationSize, t, p)

		//fmt.Println("Best from new:", BestSolution(population).weightedSum())

		fmt.Println("Used", time.Since(startGeneration).Seconds(), "seconds for generation")
		saveParetoPlotter(p, "pareto.png")
	}

	saveParetoPlotter(p, "pareto.png")

	fmt.Println("\nUsed", time.Since(start).Seconds(), "seconds to evolve solutions")

	fronts := fastNonDominatedSort(population)

	bestFronts := make([]*Solution, len(fronts[0]))
	for i, f := range fronts[0] {
		bestFronts[i] = population[f]
	}
	return bestFronts
}
