package main

import (
	"image"
	"image/color"
	"math"
)

type Genotype []uint64

type Solution struct {
	genotype         Genotype
	deviation        float64
	connectivity     float64
	crowdingDistance float64
}

func (s *Solution) weightedSum() float64 {
	return s.deviation + s.connectivity*10000000
}

func BestSolution(solutions []*Solution) *Solution {
	bestFitness := math.MaxFloat64
	bestIdx := -1
	for i, s := range solutions {
		f := s.weightedSum()
		if f < bestFitness {
			bestIdx = i
			bestFitness = f
		}
	}
	return solutions[bestIdx]
}

type Pixel struct {
	r int16
	g int16
	b int16
}

func (px *Pixel) toRGBA() color.RGBA {
	return color.RGBA{uint8(px.r), uint8(px.g), uint8(px.b), 0xFF}
}

type Image [][]Pixel

func (img *Image) toRGBA() *image.RGBA {
	width := len(*img)
	height := len((*img)[0])
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := range *img {
		for j, px := range (*img)[i] {
			rgba.Set(i, j, px.toRGBA())
		}
	}
	return rgba
}
