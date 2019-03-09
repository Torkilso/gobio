package main

import (
	"fmt"
)

func main() {

	imagePath := "./data/216066/Test_image.jpg"
	image := readJPEGFile(imagePath)

	//solutions := nsgaII(&image, 100, 100)
	//_ = solutions

	//runGenerations(&image)
	nsgaII(&image, 10, 80)
}

func runGenerations(img *Image) {

	pop := GeneratePopulation(img, 10)

	for i := 0; i < 100; i++ {
		pop = RunGeneration(img, pop)
		sol := pop.BestSolution()
		graph := GenoToGraph(img, sol.genotype)
		groups := graph.ConnectedComponents()
		width := len(*img)

		thisImg := img.toRGBA()
		for _, g := range groups {
			c := Centroid(img, g)
			for k := range g {
				x, y := Flatten(width, int(k))
				thisImg.Set(x, y, c.toRGBA())
			}
			SaveJPEGRaw(thisImg)
		}

		fmt.Println("Gen", i, "Best", sol.weightedSum(), "Segments", len(groups))
	}
}
