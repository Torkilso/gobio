package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"image/color"
	"math/rand"
	"time"
)

func main() {

	//imagePath := "./data/216066/Test_image.jpg"
	//imagePath := "./testimages/Untitled2.jpg"
	//image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())

	//runGenerations(&image)
	//runNSGA(&image)
	//runSoniaMST(&image)
	runNSGAOnTestFolder("216066")
}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 1, 2)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts)

	for id, solution := range solutions {
		graph := GenoToGraph(img, solutions[0].genotype)
		segments := graph.ConnectedComponents()
		fmt.Println("Solution", id, ": segments:", len(segments), ", c:", solution.connectivity, ", d:", solution.deviation)
	}
	graph := GenoToGraph(img, solutions[0].genotype)

	//visualizeImageGraph("mstgraph.png", img, graph)

	edgedImg := DrawImageBoundries(img, graph, color.Black)
	SaveJPEGRaw(edgedImg, "edges.jpg")
}

func runNSGAOnTestFolder(folderId string) {
	imagePath := "./data/" + folderId + "/colors.jpg"
	image := readJPEGFile(imagePath)
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	// Set max and min connectivity and deviation
	setObjectivesMaxMinValues(&image)

	solutions := nsgaII(&image, 100, 80)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fmt.Println()
	fmt.Println("Solutions:")

	for id, s := range solutions {
		graph := GenoToGraph(&image, s.genotype)
		segments := graph.ConnectedComponents()
		fmt.Println("Solution", id, ": segments:", len(segments), ", c:", s.connectivity, ", d:", s.deviation)

		if len(segments) > 1 {

			drawSolutionSegmentsBorders(&image, s, color.Black, "border.png")
			drawSolutionSegmentsWithCentroidColor(&image, s, "segments.png")
		}
	}
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

	mst2 := Prim(uint64(start), primGraph, labels, imgAsGraph)
	fmt.Println("Time to prim", time.Now().Sub(startT))

	mstGraph := graphs.GetGraph(mst2, false)
	visualizeImageGraph("mst.png", img, mstGraph)

}
