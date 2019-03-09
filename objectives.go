package main

import (
	"github.com/alonsovidales/go_graph"
	"image"
	"math"
)

var (
	maxDeviation    float64 = 100
	minDeviation    float64 = 0
	maxConnectivity float64 = 100
	minConnectivity float64 = 0
)

func distRGB(p1 *Pixel, p2 *Pixel) float64 {
	return math.Sqrt(math.Pow(float64(p1.r-p2.r), 2) + math.Pow(float64(p1.g-p2.g), 2) + math.Pow(float64(p1.b-p2.b), 2))
}

func deviation(img *image.Image, graph *graphs.Graph) float64 {
	connectedGroups := graph.ConnectedComponents()

	var dist float64
	rect := (*img).Bounds()
	for _, group := range connectedGroups {
		centroid := Centroid(img, group)

		for k := range group {
			x, y := Flatten(rect, int(k))
			dist += ColorDist((*img).At(x, y), centroid)
		}
	}
	return dist
}

func connectiviy() float64 {
	return 0.0
}
