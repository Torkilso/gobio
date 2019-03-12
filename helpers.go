package main

import (
	"math"
)

func Flatten(width int, idx int) (x, y int) {
	x = idx % width
	y = (idx - x) / width
	return x, y
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

	return ColorDist(&fromPx, &toPx)
}

func ColorDist(p1 *Pixel, p2 *Pixel) float64 {
	return math.Sqrt(math.Pow(float64(p1.r-p2.r), 2) + math.Pow(float64(p1.g-p2.g), 2) + math.Pow(float64(p1.b-p2.b), 2))
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
