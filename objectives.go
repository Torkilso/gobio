package main

import (
	"github.com/alonsovidales/go_graph"
)

var (
	maxDeviation    float64 = 100
	minDeviation    float64 = 0
	maxConnectivity float64 = 100
	minConnectivity float64 = 0
)


func deviation(img *Image, graph *graphs.Graph) float64 {
	connectedGroups := graph.ConnectedComponents()

	var dist float64
	width := len(*img)
	for _, group := range connectedGroups {
		centroid := Centroid(img, group)

		for k := range group {
			x, y := Flatten(width, int(k))
			dist += ColorDist(&(*img)[x][y], centroid)
		}
	}
	return dist
}


func connectiviy() float64 {
	return 0.0
}
