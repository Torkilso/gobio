package main

import (
	"image"
	"image/color"
)

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

func (px *Pixel) toRGBA() color.RGBA {
	return color.RGBA{uint8(px.r), uint8(px.g), uint8(px.b), 0xFF}
}
type Image [][]Pixel


func (img *Image) toRGBA() *image.RGBA{
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