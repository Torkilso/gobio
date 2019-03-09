package main

import (
	"github.com/alonsovidales/go_graph"
	"image"
)

func GetTargets(img *image.Image, idx int) []int {
	rect := (*img).Bounds()
	x, y := Flatten((*img).Bounds(), idx)

	nodes := make([]int, 0, 8)

	// TOP
	if y > 0 {
		nodes = append(nodes, Expand(rect, x, y-1))
	}
	// TOP-RIGHT
	if y > 0 && x < rect.Max.X {
		nodes = append(nodes, Expand(rect, x+1, y-1))
	}
	// RIGHT
	if x < rect.Max.X {
		nodes = append(nodes, Expand(rect, x+1, y))
	}
	//BOTTOM-RIGHT
	if x < rect.Max.X && y < rect.Max.Y {
		nodes = append(nodes, Expand(rect, x+1, y+1))
	}
	// BOTTOM
	if y < rect.Max.Y {
		nodes = append(nodes, Expand(rect, x, y+1))
	}
	// BOTTOM-LEFT
	if y < rect.Max.Y && x > 0 {
		nodes = append(nodes, Expand(rect, x-1, y+1))
	}
	// LEFT
	if x > 0 {
		nodes = append(nodes, Expand(rect, x-1, y))
	}
	// TOP-LEFT
	if y > 0 && x > 0 {
		nodes = append(nodes, Expand(rect, x-1, y-1))
	}
	return nodes
}

func GenerateGraph(img *image.Image) *graphs.Graph {
	rect := (*img).Bounds()
	numPixels := rect.Max.X * rect.Max.Y

	edges := make([]graphs.Edge, 0, numPixels*8)
	for i := 0; i < numPixels; i++ {
		targets := GetTargets(img, i)
		for _, target := range targets {
			edges = append(edges, graphs.Edge{uint64(i), uint64(target), Dist(img, i, target)})
		}
	}
	return graphs.GetGraph(edges, true)
}
