package main

import (
	"github.com/alonsovidales/go_graph"
)

func GetTargets(img *Image, idx int) []int {

	width := len(*img)
	height := len((*img)[0])
	x, y := Flatten(len(*img), idx)

	nodes := make([]int, 0, 9)

	nodes = append(nodes, idx)
	// TOP
	if y > 0 {
		//fmt.Println("TOP")

		nodes = append(nodes, Expand(width, x, y-1))
	}
	// TOP-RIGHT
	if y > 0 && x < width-1 {
		//fmt.Println("TOP-RIGHT")

		nodes = append(nodes, Expand(width, x+1, y-1))
	}
	// RIGHT
	if x < width-1 {
		//fmt.Println("RIGHT")

		nodes = append(nodes, Expand(width, x+1, y))
	}
	//BOTTOM-RIGHT
	if x < width-1 && y < height-1 {
		//	fmt.Println("BOTTOM-RIGHT")

		nodes = append(nodes, Expand(width, x+1, y+1))
	}
	// BOTTOM
	if y < height-1 {
		//fmt.Println("BOTTOM")

		nodes = append(nodes, Expand(width, x, y+1))

	}
	// BOTTOM-LEFT
	if y < height-1 && x > 0 {
		//fmt.Println("BOTTOM-LEFT")

		nodes = append(nodes, Expand(width, x-1, y+1))

	}
	// LEFT
	if x > 0 {
		//fmt.Println("LEFT")

		nodes = append(nodes, Expand(width, x-1, y))
	}
	// TOP-LEFT
	if y > 0 && x > 0 {
		//fmt.Println("TOP-LEFT")

		nodes = append(nodes, Expand(width, x-1, y-1))
	}
	return nodes
}

func GetCloseTargets(img *Image, idx int) []int {
	width := len(*img)
	height := len((*img)[0])
	x, y := Flatten(len(*img), idx)

	nodes := make([]int, 0, 4)

	// TOP
	if y > 0 {
		//fmt.Println("TOP")

		nodes = append(nodes, Expand(width, x, y-1))
	}
	// RIGHT
	if x < width-1 {
		//fmt.Println("RIGHT")

		nodes = append(nodes, Expand(width, x+1, y))
	}
	// BOTTOM
	if y < height-1 {
		//fmt.Println("BOTTOM")

		nodes = append(nodes, Expand(width, x, y+1))

	}
	// LEFT
	if x > 0 {
		//fmt.Println("LEFT")

		nodes = append(nodes, Expand(width, x-1, y))
	}
	return nodes
}

func GenerateGraph(img *Image) *graphs.Graph {
	numPixels := len(*img)*len((*img)[0]) - 1
	edges := make([]graphs.Edge, 0, numPixels*8)
	for i := 0; i < numPixels; i++ {
		targets := GetTargets(img, i)

		for _, target := range targets {
			edges = append(edges, graphs.Edge{uint64(i), uint64(target), Dist(img, i, target)})
		}
	}
	return graphs.GetGraph(edges, true)
}
