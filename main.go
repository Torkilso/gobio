package main

import (
	"fmt"
	"image/color"
	"log"
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
	solutions := nsgaII(img, 3, 10)

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts)

	graph := GenoToGraph(img, solutions[0].genotype)
	thisImg := img.toRGBA()

	imgCopy := GoImageToImage(thisImg)

	edgedImg := DrawImageBoundries(&imgCopy, graph, color.Black)
	SaveJPEGRaw(edgedImg, "edges.jpg")

}

func runGenerations(img *Image) {

	pop := GeneratePopulation(img, 100)

	for i := 0; i < 100; i++ {
		pop = RunGeneration(img, pop)
		sol := pop.BestSolution()
		graph := GenoToGraph(img, sol.genotype)
		groups := graph.ConnectedComponents()
		width := len(*img)

		thisImg := img.toRGBA()
		imgCopy := GoImageToImage(thisImg)
		for _, g := range groups {
			c := Centroid(img, g)
			log.Println("Centroid", c)
			for k := range g {
				x, y := Flatten(width, int(k))
				thisImg.Set(x, y, c.toRGBA())
			}
			SaveJPEGRaw(thisImg, "img.jpg")
		}

		edgedImg := DrawImageBoundries(&imgCopy, graph, color.Black)
		SaveJPEGRaw(edgedImg, "edges.jpg")

		fmt.Println("Gen", i, "Best", sol.weightedSum(), "Segments", len(groups))

	}
}
