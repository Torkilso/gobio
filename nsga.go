package main

import (
	"fmt"
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

func nsgaII(image *Image, generations, populationSize int) []*Solution {
	start := time.Now()

	fmt.Println("Initiating NSGAII")
	fmt.Println("Generating", populationSize, "solutions")

	population := generatePopulation(image, populationSize)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds to generate solutions")
	fmt.Println("\nEvolving solutions")

	start = time.Now()

	p := createParetoPlotter()

	for t := 0; t < generations; t++ {

		startGeneration := time.Now()

		population.evolve(image)


		/*fmt.Println("Solutions in generation population")
		for id, sol := range population {
			graph := GenoToGraph(image, sol.genotype)
			segments := graph.ConnectedComponents()
			fmt.Println("Solution", id, ": segments:", len(segments), ", c:", sol.connectivity, ", d:", sol.deviation)
		}*/

		fronts := fastNonDominatedSort(population)

		fmt.Println("Generation:", t, "Best before:", BestSolution(population).weightedSum(), "Num fronts:", len(fronts))


		addParetoFrontToPlotter(p, population, fronts, t)


		newParents := make([]*Solution, 0)
		i := 0

		fmt.Println("Adding fronts", len(newParents), len(fronts[i]), populationSize)
		for len(newParents)+len(fronts[i]) <= populationSize {
			//fmt.Println("Best now", BestSolution(newParents).weightedSum(), len(fronts[i]))

			if len(fronts[i]) == 0 {
				fmt.Println("Len(fronts[i])", fronts[i])
				break
			}

			crowdingDistanceAssignment(fronts[i], population)
			frontSolutions := make([]*Solution, len(fronts[i]))

			for i, id := range fronts[i] {
				frontSolutions[i] = population[id]
			}

			newParents = append(newParents, frontSolutions...)
			//fmt.Println("Best now", BestSolution(newParents).weightedSum())
			i++
		}

		lastFrontier := make([]*Solution, 0)

		if len(fronts[i]) > 0 {
			crowdingDistanceAssignment(fronts[i], population)

			for _, id := range fronts[i] {
				if len(lastFrontier)+len(newParents) < populationSize {
					lastFrontier = append(lastFrontier, population[id])
				} else {
					break
				}
			}
		}

		population = append(newParents, lastFrontier...)

		//fmt.Println("Best from new:", BestSolution(population).weightedSum())

		//children = createPopulationFromParents(image, parents)

		i = 0

		fmt.Println("Used", time.Since(startGeneration).Seconds(), "seconds for generation")
		fmt.Println()
		saveParetoPlotter(p, "pareto.png")

	}

	saveParetoPlotter(p, "pareto.png")

	fmt.Println("Used", time.Since(start).Seconds(), "seconds to evolve solutions")

	return population
}
