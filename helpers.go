package main

import (
	"math"
)

func Flatten(width int, idx int) (x, y int) {
	x = idx % width
	y = (idx - x) / width
	return x, y
}

func ImageSize(img *Image) int {
	return len(*img) * len((*img)[0])
}

func Expand(width int, x, y int) int {
	return y*width + x
}
func Dist(img *Image, idx1, idx2 int) float64 {

	width := len(*img)
	x1, y1 := Flatten(width, idx1)
	x2, y2 := Flatten(width, idx2)

	fromPx := (*img)[x1][y1]
	toPx := (*img)[x2][y2]
	return math.Round(ColorDist(fromPx, toPx))
}

func ColorDist(p1 *Pixel, p2 *Pixel) float64 {
	r, g, b := float64(p1.r-p2.r), float64(p1.g - p2.g), float64(p1.b - p2.b)
	return math.Sqrt(r*r + g*g + b*b)
}

func Centroid(img *Image, group map[uint64]bool) *Pixel {
	var r int
	var g int
	var b int

	var count int

	width := len(*img)

	for k := range group {
		x, y := Flatten(width, int(k))
		px := (*img)[x][y]
		r1, g1, b1 := int(px.r), int(px.g), int(px.b)
		r += r1
		g += g1
		b += b1
		count++
	}
	r = r / count
	g = g / count
	b = b / count

	return &Pixel{int16(r), int16(g), int16(b)}
}
