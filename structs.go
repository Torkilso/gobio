package main

type Genotype []uint64

type Solution struct {
	genotype         Genotype
	deviation        float64
	connectivity     float64
	crowdingDistance float64
}

func (s *Solution) weightedSum() float64 {
	return s.deviation + s.connectivity
}

type Population struct {
	solutions []Solution
}

type Pixel struct {
	r int16
	g int16
	b int16
}

type Image [][]Pixel
