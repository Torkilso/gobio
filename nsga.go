package main

import (
	"math"
	"sort"
)

type SearchHelper struct {
	dominates         []int
	dominatedByAmount int
}

func fastNonDominatedSort(population []*Solution) map[int][]int {

	fronts := make(map[int][]int)
	SearchHelperMap := make(map[int]*SearchHelper)

	for i, solution := range population {
		searchHelper := SearchHelper{
			dominatedByAmount: 0,
		}

		SearchHelperMap[i] = &searchHelper

		for j, opponent := range population {
			if i == j {
				continue
			}

			if solution.dominate(opponent) {
				searchHelper.dominates = append(SearchHelperMap[i].dominates, j)
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
			for _, solution := range SearchHelperMap[frontSolution].dominates {
				SearchHelperMap[solution].dominatedByAmount--
				if SearchHelperMap[solution].dominatedByAmount == 0 {
					newFront = append(newFront, solution)
				}
			}
		}

		frontRank++
		fronts[frontRank] = newFront
	}

	return fronts
}

func crowdingDistanceAssignment(ids []int, population []*Solution) {
	size := len(ids)

	// for deviation
	sort.Slice(ids, func(i, j int) bool {
		return population[ids[i]].deviation > population[ids[j]].deviation
	})

	population[ids[0]].crowdingDistance = math.Inf(1)
	population[ids[size-1]].crowdingDistance = math.Inf(1)

	for i := 1; i < size-1; i++ {
		population[ids[i]].crowdingDistance = (population[ids[i+1]].deviation - population[ids[i-1]].deviation) / (maxDeviation - minDeviation)
	}

	// for connectivity
	sort.Slice(ids, func(i, j int) bool {
		return population[ids[i]].connectivity > population[ids[j]].connectivity
	})

	for i := 1; i < size-1; i++ {
		population[ids[i]].crowdingDistance = population[ids[i]].crowdingDistance + (population[ids[i+1]].connectivity-population[ids[i-1]].connectivity)/(maxConnectivity-minConnectivity)
	}
}

func nsgaII(image *Image, generations, populationSize int) []*Solution {

	parents := createInitialPopulation(*image, populationSize)
	children := make([]*Solution, 0)

	for t := 0; t < generations; t++ {

		population := append(parents, children...)
		fronts := fastNonDominatedSort(population)
		newParents := make([]*Solution, 0)
		i := 0

		for len(newParents)+len(fronts[i]) <= populationSize {

			crowdingDistanceAssignment(fronts[i], population)
			frontSolutions := make([]*Solution, 0)

			for _, id := range fronts[i] {
				frontSolutions = append(frontSolutions, population[id])
			}

			newParents = append(newParents, frontSolutions...)
			i++
		}

		lastFrontier := make([]*Solution, 0)
		crowdingDistanceAssignment(fronts[i], population)

		for _, id := range fronts[i] {
			if len(newParents) < populationSize {
				lastFrontier = append(lastFrontier, population[id])
			}
		}

		parents = append(newParents, lastFrontier...)
		children = createPopulationFromParents(parents)
	}

	return children
}
