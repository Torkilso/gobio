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

	fmt.Println("Initiating NSGAII")
	fmt.Println("Generating", populationSize, "solutions")

	start := time.Now()

	parents := GeneratePopulation(image, populationSize)

	children := make([]*Solution, 0)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds to generate solutions")
	fmt.Println()
	fmt.Println("Evolving solutions")

	return parents

	start = time.Now()

	for t := 0; t < generations; t++ {
		fmt.Println("Generation:", t)
		startGeneration := time.Now()

		population := append(parents, children...)
		fronts := fastNonDominatedSort(population)
		newParents := make([]*Solution, 0)
		i := 0

		for len(newParents)+len(fronts[i]) <= populationSize {
			if len(fronts[i]) == 0 {
				break
			}

			crowdingDistanceAssignment(fronts[i], population)
			frontSolutions := make([]*Solution, 0)

			for _, id := range fronts[i] {
				frontSolutions = append(frontSolutions, population[id])
			}

			newParents = append(newParents, frontSolutions...)
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

		parents = append(newParents, lastFrontier...)

		children = createPopulationFromParents(image, parents)

		newParents = make([]*Solution, 0)
		i = 0

		fmt.Println("Used", time.Since(startGeneration).Seconds(), "seconds")
		fmt.Println()
	}

	fmt.Println("Used", time.Since(start).Seconds(), "seconds to evolve solutions")

	for _, sol := range children {

		if sol.deviation == 0 {
			fmt.Println(sol)
		}
	}

	return children
}
