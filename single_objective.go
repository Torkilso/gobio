package main

import (
	"fmt"
)

func singleObjective(image *Image, generations, populationSize int) *Solution {
	population := generatePopulation(image, populationSize)

	for i := 0; i < generations; i++ {
		population.evolveSingleObjective(image)
		sol := bestSolution(population)

		groups := GenoToConnectedComponents(sol.genotype)
		fmt.Println("Gen", i, "Best", sol.weightedSum(), "Segments", len(groups))
	}
	return bestSolution(population)
}
