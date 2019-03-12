package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
)

func main() {

	imagePath := "./testimages/Untitled.jpg"
	//imagePath := "./data/216066/Test_image.jpg"
	image := readJPEGFile(imagePath)

	//solutions := nsgaII(&image, 100, 100)
	//_ = solutions
	rand.Seed(time.Now().UTC().UnixNano())

	//runGenerations(&image)
	runNSGA(&image)
}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 100, 80)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts)

	graph := GenoToGraph(img, solutions[0].genotype)
	segments := graph.ConnectedComponents()
	fmt.Println("Amount of segments:", len(segments))
	visualizeImageGraph("mstgraph.png", img, graph)

	//thisImg := img.toRGBA()

	//imgCopy := GoImageToImage(thisImg)

	//edgedImg := DrawImageBoundries(&imgCopy, graph, color.Black)
	//SaveJPEGRaw(edgedImg, "edges.jpg")
}

func runGenerations(img *Image) {

	pop := GeneratePopulation(img, 2)
	sol := BestSolution(pop)
	graph := GenoToGraph(img, sol.genotype)

	visualizeImageGraph("graph.png", img, graph)

	for i := 0; i < 1; i++ {
		pop = RunGeneration(img, pop)
		sol := BestSolution(pop)

		graph := GenoToGraph(img, sol.genotype)
		groups := graph.ConnectedComponents()
		width := len(*img)

		thisImg := img.toRGBA()
		imgCopy := GoImageToImage(thisImg)
		fmt.Println("Number of groups", len(groups), "Pixels", len(*img)*len((*img)[0]), groups, graph.RawEdges)
		for _, g := range groups {
			c := Centroid(img, g)
			fmt.Println("Centroid", c)
			for k := range g {
				x, y := Flatten(width, int(k))
				thisImg.Set(x, y, c.toRGBA())
			}
		}

		edgedImg := DrawImageBoundries(&imgCopy, graph, color.Black)
		SaveJPEGRaw(edgedImg, "edges.jpg")
		visualizeImageGraph("graph.png", img, graph)
		SaveJPEGRaw(thisImg, "img.jpg")

		fmt.Println("Gen", i, "Best", sol.weightedSum(), "Segments", len(groups))
	}
}
