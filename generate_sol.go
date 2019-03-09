package main

import (
	"github.com/alonsovidales/go_graph"
)

func GetTargets(img *Image, idx int) []int {
	width := len(*img) - 1
	height := len((*img)[0]) - 1
	x, y := Flatten(len(*img), idx)

	nodes := make([]int, 0, 8)

	nodes = append(nodes, idx)
	// TOP
	if y > 0 {
		nodes = append(nodes, Expand(width, x, y-1))
	}
	// TOP-RIGHT
	if y > 0 && x < width {
		nodes = append(nodes, Expand(width, x+1, y-1))
	}
	// RIGHT
	if x < width {
		nodes = append(nodes, Expand(width, x+1, y))
	}
	//BOTTOM-RIGHT
	if x < width && y < height {
		nodes = append(nodes, Expand(width, x+1, y+1))
	}
	// BOTTOM
	if y < height {
		nodes = append(nodes, Expand(width, x, y+1))

	}
	// BOTTOM-LEFT
	if y < height && x > 0 {
		nodes = append(nodes, Expand(width, x-1, y+1))

	}
	// LEFT
	if x > 0 {
		nodes = append(nodes, Expand(width, x-1, y))
	}
	// TOP-LEFT
	if y > 0 && x > 0 {
		nodes = append(nodes, Expand(width, x-1, y-1))
	}
	return nodes
}

func GenerateGraph(img *Image) *graphs.Graph {
	numPixels := len(*img) * len((*img)[0])
	edges := make([]graphs.Edge, 0, numPixels*8)
	for i := 0; i < numPixels; i++ {
		targets := GetTargets(img, i)
		for _, target := range targets {
			edges = append(edges, graphs.Edge{uint64(i), uint64(target), Dist(img, i, target)})
		}
	}
	return graphs.GetGraph(edges, true)
}
