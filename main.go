package main

import (
	"fmt"
	"math"
)

func main() {


	imagePath := "./data/216066/Test image_tn.jpg"
	image := readJPEGFile(imagePath)

	//solutions := nsgaII(&image, 100, 100)
	//_ = solutions

	runGenerations(&image)
}


func runGenerations(img *Image){

	pop := GeneratePopulation(img, 10)

	for i := 0 ; i < 100 ; i++ {
		pop = RunGeneration(img, pop)
		sol := pop.solutions[0]
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
		var best float64
		var avg float64
		for _, s := range pop.solutions {
			best = math.Max(best, s.weightedSum())
			avg += s.weightedSum()
		}
		fmt.Println("Gen", i, "Best", best, "Avg", avg / float64(len(pop.solutions)) )
	}
}
