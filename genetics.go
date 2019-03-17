package main

import (
	"fmt"
	"github.com/google/gxui/math"
	"math/rand"
	"sync"
	"time"
)

func tournamentWeighted(solutions []*Solution) int {

	idx1 := rand.Intn(len(solutions))
	idx2 := rand.Intn(len(solutions))

	if solutions[idx1].weightedSum() < solutions[idx2].weightedSum() {
		return idx1
	}
	return idx2
}

/**
1. Make a slice to hold visited, where value is segment number
2. Make a list of maps, which holds the segments
3. For each item in geno
	1. Check if it has been visited, if so, continue.
	2. Make a group to hold
*/

func GenoToConnectedComponents(geno []uint64) []map[uint64]bool {

	edges := make([][]uint64, len(geno))
	for i := range geno {
		edges[i] = []uint64{geno[i]}
	}
	for i := range geno {
		edges[geno[i]] = append(edges[geno[i]], uint64(i))
	}
	segments := make([]map[uint64]bool, 0)
	used := make([]bool, len(geno))
	for from := range geno {
		if !used[from] {
			used[from] = true
			next := make([]uint64, 0)
			group := map[uint64]bool{uint64(from): true}
			next = append(next, edges[from]...)
			for len(next) > 0 {
				current := next[0]
				if !used[current] {
					group[current] = true
					next = append(next, edges[current]...)
				}
				used[current] = true
				next = next[1:]
			}
			segments = append(segments, group)
		}
	}
	return segments
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

	channel := make(chan *Solution, n)
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func(index int) {
			start := r1.Intn(width * height)
			geno := Prim(uint64(start), primGraph, labels, imgAsGraph, width*height)

			groups := GenoToConnectedComponents(geno)
			d := deviation(img, groups)
			s2 := &Solution{geno, d, connectivity(img, groups), 0.0,0, len(groups)}
			channel <- s2
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
func (p *Population) evolveSingleObjective(img *Image) {
	size := len(*p)
	result := make([]*Solution, 0, size)

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(size + size/2)

	for i := 0; i < size; i += 2 {
		go func(index int) {
			p1Idx := tournamentWeighted(*p)
			p2Idx := tournamentWeighted(*p)

			p1 := (*p)[p1Idx]
			p2 := (*p)[p2Idx]

			leftChild, rightChild := crossover(img, p1, p2)

			if rand.Float32() < 0.2 {
				leftChild.mutate(img)
			}
			if rand.Float32() < 0.2 {
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

	*p = result
}

func (p *Population) evolveWithTournament(img *Image) {
	size := len(*p)

	parentsA := make([]*Solution, 0, size)
	parentsB := make([]*Solution, 0, size)

	for i := 0; i < size; i += 2 {
		if rand.Float32() > 0.7 {
			continue
		}
		p1 := tournamentNSGA(p, 4)
		p2 := tournamentNSGA(p, 4)

		parentsA = append(parentsA, p1)
		parentsB = append(parentsB, p2)
	}

	resultSize := len(parentsA)
	result := make([]*Solution, 0, resultSize*2)

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(resultSize + resultSize*2)

	for i := 0; i < resultSize; i++ {
		go func(index int) {

			leftChild, rightChild := crossover(img, parentsA[index], parentsB[index])

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

func tournamentNSGA(p *Population, k int) *Solution {
	bestIdx := -1
	bestFront := math.MaxInt
	bestCrowding := -1.0

	for i := 0 ; i <= k ; i++ {
		id := rand.Intn(len(*p))

		if (*p)[id].frontNumber < bestFront {
			bestIdx = id
			bestFront = (*p)[id].frontNumber
			bestCrowding = (*p)[id].crowdingDistance
		} else if (*p)[id].frontNumber == bestFront && (*p)[id].crowdingDistance > bestCrowding {
			bestIdx = id
			bestFront = (*p)[id].frontNumber
			bestCrowding = (*p)[id].crowdingDistance
		}
	}
	return (*p)[bestIdx]
}

func (s *Solution) mutate(img *Image) {

	index := rand.Intn(len(s.genotype))
	possibleValues := GetTargets(img, index, true)
	chosen := rand.Intn(len(possibleValues))

	s.genotype[uint64(index)] = uint64(possibleValues[chosen])

	groups := GenoToConnectedComponents(s.genotype)

	s.connectivity = connectivity(img, groups)
	s.deviation = deviation(img, groups)
	s.crowdingDistance = 0.0
	s.amountOfSegments = len(groups)
}

func (s *Solution) mutateMultiple(img *Image) {
	mutated := false

	/*groups := GenoToConnectedComponents(s.genotype)

	for _, group := range groups {
		if len(group) > 5 {
			continue
		}

		for i := range group {
			possibleValues := GetTargets(img, int(i))
			chosen := rand.Intn(len(possibleValues))
			s.genotype[uint64(i)] = uint64(possibleValues[chosen])

			mutated = true
		}
	}*/

	for i := range s.genotype {
		if rand.Float32() < 0.00001 {
			possibleValues := GetCloseTargetsWithSelf(img, i)
			chosen := rand.Intn(len(possibleValues))
			s.genotype[uint64(i)] = uint64(possibleValues[chosen])

			mutated = true
		}
	}

	if mutated {
		groups := GenoToConnectedComponents(s.genotype)
		s.connectivity = connectivity(img, groups)
		s.deviation = deviation(img, groups)
		s.crowdingDistance = 0.0
		s.amountOfSegments = len(groups)
	}
}

func crossover(img *Image, parent1, parent2 *Solution) (*Solution, *Solution) {

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

	if rand.Float32() < 0.2 {
		index := rand.Intn(len(offspring1))
		possibleValues := GetCloseTargetsWithSelf(img, index)
		chosen := rand.Intn(len(possibleValues))
		offspring1[uint64(index)] = uint64(possibleValues[chosen])
	}
	if rand.Float32() < 0.2 {
		index := rand.Intn(len(offspring2))
		possibleValues := GetCloseTargetsWithSelf(img, index)
		chosen := rand.Intn(len(possibleValues))
		offspring2[uint64(index)] = uint64(possibleValues[chosen])
	}

	groups1 := GenoToConnectedComponents(offspring1)
	groups2 := GenoToConnectedComponents(offspring2)

	/*
	if len(groups1) > 500 {
		sort.Slice(groups1, func(i, j int) bool {
			return len(groups1[i]) > len(groups1[j])
		})

		before := len(groups1)
		for _, group := range groups1 {
			for i := range group {
				for _, neighbour := range GetCloseTargetsWithSelf(img, int(i)) {
					if _, ok := group[uint64(neighbour)]; !ok { // Not in same segment
						offspring1[int(i)] = uint64(neighbour)
						break
					}
				}
			}
		}
		groups1 = GenoToConnectedComponents(offspring1)
		fmt.Println("Num segments", before, len(groups1))

	}
	if rand.Float32() < 0.1 {
		for _, group := range groups2 {
			if len(group) > 100 {
				continue
			}
			for i := range group {
				for _, neighbour := range GetCloseTargetsWithSelf(img, int(i)) {
					if _, ok := group[uint64(neighbour)]; !ok { // Not in same segment
						offspring2[int(i)] = uint64(neighbour)
						break
					}
				}
			}
		}
		groups2 = GenoToConnectedComponents(offspring2)

	}
	*/


	s1 := &Solution{
		offspring1,
		deviation(img, groups1),
		connectivity(img, groups1),
		0.0,
		0,
		len(groups2),
	}

	s2 := &Solution{
		offspring2,
		deviation(img, groups2),
		connectivity(img, groups2),
		0.0,
		0,
		len(groups2),
	}

	return s1, s2
}


func (p *Population) joinSegments(img *Image, segmentSizeThreshold int) {
	var wg sync.WaitGroup
	wg.Add(len(*p))

	for i := 0; i < len(*p); i++ {
		go func(index int) {
			if (*p)[index].amountOfSegments < 5000 {
				for {
					foundGroup := false
					groupsInner := GenoToConnectedComponents((*p)[index].genotype)

					for _, group := range groupsInner {
						if len(group) > segmentSizeThreshold {
							continue
						}
						for i := range group {
							possibleValues := GetTargets(img, int(i), true)
							chosen := rand.Intn(len(possibleValues))
							(*p)[index].genotype[uint64(i)] = uint64(possibleValues[chosen])
						}

						foundGroup = true
					}
					if !foundGroup {
						break
					}
				}

				newGroups := GenoToConnectedComponents((*p)[index].genotype)
				connectivity := connectivity(img, newGroups)

				(*p)[index].connectivity = connectivity
				(*p)[index].deviation = deviation(img, newGroups)
				(*p)[index].crowdingDistance = 0.0
				(*p)[index].amountOfSegments = len(newGroups)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func (p *Population) expandWithSolutions(img *Image, amount int) {

	parentsA := make([]*Solution, 0, amount/2)
	parentsB := make([]*Solution, 0, amount/2)

	for i := 0; i < amount; i += 2 {
		p1 := tournamentNSGA(p, 4)
		p2 := tournamentNSGA(p, 4)

		parentsA = append(parentsA, p1)
		parentsB = append(parentsB, p2)
	}

	resultSize := len(parentsA)
	result := make([]*Solution, 0, resultSize*2)

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(resultSize + resultSize*2)

	for i := 0; i < resultSize; i++ {
		go func(index int) {

			leftChild, rightChild := crossover(img, parentsA[index], parentsB[index])

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

func (p *Population) hasTakenLeapOfFaith() bool {
	for _, s := range *p {
		if s.amountOfSegments < 500 && s.amountOfSegments > 1 {
			return true
		}
	}
	return false
}
