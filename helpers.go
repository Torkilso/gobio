package main

import (
	"image"
	"image/color"
	"math"
)

func Flatten(rect image.Rectangle, idx int) (x, y int) {
	w := rect.Dx()
	x = idx % w
	y = (idx - x) / w
	return x, y
}

func Expand(rect image.Rectangle, x, y int) int {
	return y * rect.Dx() + x
}
func Dist(img *image.Image, idx1, idx2 int) float64 {
	rect := (*img).Bounds()
	x1, y1 := Flatten(rect, idx1)
	x2, y2 := Flatten(rect, idx2)

	fromPx := (*img).At(x1, y1)
	toPx := (*img).At(x2, y2)
	return ColorDist(fromPx, toPx)
}

func ColorDist(c1 color.Color, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	return math.Pow(float64(r1-r2), 2) + math.Pow(float64(g1-g2), 2) + math.Pow(float64(b1-b2), 2)

}



func Centroid(img *image.Image, group map[uint64]bool) color.Color {

	rect := (*img).Bounds()
	var r uint32
	var g uint32
	var b uint32

	var count uint32

	for k := range group {
		x, y := Flatten(rect, int(k))
		c := (*img).At(x, y)
		r1, g1, b1, _ := c.RGBA()
		r += r1
		g += g1
		b += b1
		count++
	}
	r = r / count
	g = g / count
	b = b / count

	return color.RGBA{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(0xFF)}
}