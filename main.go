package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"image/color"
	"math/rand"
	"time"
)

func main() {

	imagePath := "./data/216066/Test image.jpg"
	//imagePath := "./testimages/Untitled.jpg"
	image := readJPEGFile(imagePath)

	//solutions := nsgaII(&image, 100, 100)
	//_ = solutions
	rand.Seed(time.Now().UTC().UnixNano())

	//runGenerations(&image)
	runNSGA(&image)
	//runSoniaMST(&image)
}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 10, 10)

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
	return
	//sol := BestSolution(pop)
	//graph := GenoToGraph(img, sol.genotype)

	//visualizeImageGraph("graph.png", img, graph)

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

func runSoniaMST(img *Image) {
	imgAsGraph := GenerateGraph(img)

	width := len(*img)
	height := len((*img)[0])

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	start := r1.Intn(width * height)
	startT := time.Now()

	primGraph, labels := PreparePrim(imgAsGraph)
	fmt.Println("Time to prepare", time.Now().Sub(startT))
	startT = time.Now()

	mst2 := Prim(uint64(start), primGraph, labels,imgAsGraph)
	fmt.Println("Time to prim", time.Now().Sub(startT))

	mstGraph := graphs.GetGraph(mst2, false)
	visualizeImageGraph("mst.png", img, mstGraph)

}