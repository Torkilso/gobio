package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
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

	channel := make(chan *Solution, n)
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func(index int) {
			start := r1.Intn(width * height)
			geno := Prim(uint64(start), primGraph, labels, imgAsGraph, width * height)

			//mstGraph := graphs.GetGraph(mst, true)
			//s := SolutionFromGenotypeNSGA(img, mstGraph)
			groups := GenoToConnectedComponents(geno)
			e, c := connectivityAndEdge(img, groups)
			d := deviation(img, groups)
			s2 := &Solution{geno, d, c, 0.0, e, 0, len(groups)}
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

	c, e := connectivityAndEdge(img, groups)
	s.connectivity = c
	s.deviation = deviation(img, groups)
	s.edgeValue = e
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
			possibleValues := GetTargets(img, i)
			chosen := rand.Intn(len(possibleValues))
			s.genotype[uint64(i)] = uint64(possibleValues[chosen])

			mutated = true
		}
	}

	if mutated {
		groups := GenoToConnectedComponents(s.genotype)
		c, e := connectivityAndEdge(img, groups)

		s.connectivity = c
		s.deviation = deviation(img, groups)
		s.edgeValue = e
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
		possibleValues := GetTargets(img, index)
		chosen := rand.Intn(len(possibleValues))
		offspring1[uint64(index)] = uint64(possibleValues[chosen])
	}
	if rand.Float32() < 0.2 {
		index := rand.Intn(len(offspring2))
		possibleValues := GetTargets(img, index)
		chosen := rand.Intn(len(possibleValues))
		offspring2[uint64(index)] = uint64(possibleValues[chosen])
	}


	groups1 := GenoToConnectedComponents(offspring1)
	groups2 := GenoToConnectedComponents(offspring2)

	c1, e1 := connectivityAndEdge(img, groups1)
	c2, e2 := connectivityAndEdge(img, groups2)
	s1 := &Solution{
		offspring1,
		deviation(img, groups1),
		c1,
		0.0,
		e1,
		0,
		len(groups2),
	}

	s2 := &Solution{
		offspring2,
		deviation(img, groups2),
		c2,
		0.0,
		e2,
		0,
		len(groups2),
	}
	return s1, s2
}

func SolutionFromGenotypeNSGA(img *Image, g *graphs.Graph) *Solution {
	groups := g.ConnectedComponents()

	deviation := deviation(img, groups)
	connectivity, edgeValue := connectivityAndEdge(img, groups)
	crowdingDistance := 0.0

	sol := &Solution{GraphToGeno(g, ImageSize(img)),
		deviation,
		connectivity,
		crowdingDistance,
		edgeValue,
		0,
		len(groups)}

	return sol
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
							possibleValues := GetTargets(img, int(i))
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
				connectivity, edgeValue := connectivityAndEdge(img, newGroups)

				(*p)[index].connectivity = connectivity
				(*p)[index].deviation = deviation(img, newGroups)
				(*p)[index].edgeValue = edgeValue
				(*p)[index].crowdingDistance = 0.0
				(*p)[index].amountOfSegments = len(newGroups)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func (p *Population) expandWithSolutions(img *Image, amount int) {

	size := len(*p)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)


	parentsA := make([]*Solution, 0, amount/2)
	parentsB := make([]*Solution, 0, amount/2)

	for i := 0; i < amount; i += 2 {
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
