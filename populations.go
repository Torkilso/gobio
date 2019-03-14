package main

import (
	"math/rand"
	"time"
)

func createInitialPopulation(image *Image, populationSize int) []*Solution {
	population := make([]*Solution, 0)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < populationSize; i++ {
		population = append(population, &Solution{
			connectivity:     float64(r1.Intn(1000)),
			deviation:        float64(r1.Intn(1000)),
			crowdingDistance: 0.0,
		})
	}

	return population
}

func (solution *Solution) dominate(opponent *Solution) bool {
	return solution.connectivity <= opponent.connectivity && solution.deviation <= opponent.deviation
}
