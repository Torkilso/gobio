package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"math"
	"math/rand"
	"sync"
	"time"
)

func tournamentWeighted(solutions []*Solution, k int) int {

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

	channel := make(chan *Solution)
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func(index int) {
			start := r1.Intn(width * height)
			mst := Prim(uint64(start), primGraph, labels, imgAsGraph)

			mstGraph := graphs.GetGraph(mst, true)

			channel <- SolutionFromGenotypeNSGA(img, mstGraph)
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
			p1Idx := tournamentWeighted(*p, 2)
			p2Idx := tournamentWeighted(*p, 2)

			p1 := (*p)[p1Idx]
			p2 := (*p)[p2Idx]

			leftChild, rightChild := crossover(img, p1, p2)

			leftChild.mutateMultiple(img)
			rightChild.mutateMultiple(img)


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

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	parentsA := make([]*Solution, 0, size)
	parentsB := make([]*Solution, 0, size)

	for i := 0; i < size; i += 2 {
		if rand.Float32() > 0.7 {
			continue
		}
		p1Idx := r1.Intn(size)
		p2Idx := r1.Intn(size)
		p3Idx := r1.Intn(size)
		p4Idx := r1.Intn(size)

		p1 := tournamentNSGA(p1Idx, p2Idx, p)
		p2 := tournamentNSGA(p3Idx, p4Idx, p)

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

			leftChild.mutateMultiple(img)
			rightChild.mutateMultiple(img)

			/*if r1.Float32() < .1 {
				leftChild.mutate(img)
			}

			if r1.Float32() < .1 {
				rightChild.mutate(img)
			}*/

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

func tournamentNSGA(id1, id2 int, p *Population) *Solution {
	if (*p)[id1].frontNumber < (*p)[id2].frontNumber {
		return (*p)[id1]
	} else if (*p)[id1].frontNumber == (*p)[id2].frontNumber {
		if (*p)[id1].crowdingDistance > (*p)[id2].crowdingDistance {
			return (*p)[id1]
		}
	}

	return (*p)[id2]
}

func (s *Solution) mutate(img *Image) {

	index := rand.Intn(len(s.genotype))
	possibleValues := GetTargets(img, index)
	chosen := rand.Intn(len(possibleValues))

	s.genotype[uint64(index)] = uint64(possibleValues[chosen])

	groups := GenoToConnectedComponents(s.genotype)

	s.connectivity = connectivity(img, groups)
	s.deviation = deviation(img, groups)
	s.crowdingDistance = 0.0
}
func (s *Solution) mutateMultiple(img *Image) {
	mutated := false
	for i := range s.genotype {
		if rand.Float32() < 0.00001 {
			possibleValues := GetTargets(img, i)
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

	groups1 := GenoToConnectedComponents(offspring1)
	groups2 := GenoToConnectedComponents(offspring2)

	s1 := &Solution{
		offspring1, deviation(img, groups1), connectivity(img, groups1), 0.0, edgeValues(img, groups1), 0,
	}

	s2 := &Solution{
		offspring2, deviation(img, groups2), connectivity(img, groups2), 0.0, edgeValues(img, groups2), 0,
	}
	return s1, s2
}

func SolutionFromGenotypeNSGA(img *Image, g *graphs.Graph) *Solution {
	groups := g.ConnectedComponents()

	deviation := deviation(img, groups)
	connectivity := connectivity(img, groups)
	edgeValue := edgeValues(img, groups)
	crowdingDistance := 0.0
	sol := &Solution{GraphToGeno(g, ImageSize(img)), deviation, connectivity, crowdingDistance, edgeValue, 0}

	return sol
}
