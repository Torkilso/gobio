package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	generationsToRun = 200
	folderId         = "216066"
)

func main() {

	//imagePath := "./data/216066/Test_image.jpg"
	//imagePath := "./testimages/Untitled2.jpg"
	//image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())

	//defer profile.Start(profile.MemProfile).Stop()

	//runGenerations(&image)
	//runNSGA(&image)
	//runSoniaMST(&image)
	runAndStoreImagesForTesting(folderId, generationsToRun, 80)
	//runNSGAOnTestFolder("216066")
	//img := readJPEGFile("./testimages/Untitled2.jpg")
	//testMaxObjectives(&img)
	//perf()


}

func perf() {
	imagePath := "./data/216066/Test image.jpg"
	img := readJPEGFile(imagePath)
	pop := generatePopulation(&img, 1)
	s := time.Now()
	gr := GenoToGraph(&img, pop[0].genotype, false)
	seg1 := gr.ConnectedComponents()
	fmt.Println("Time graph", time.Now().Sub(s).String())
	s = time.Now()

	seg2:= GenoToConnectedComponents(pop[0].genotype)
	fmt.Println("Time homemade", time.Now().Sub(s).String())

	fmt.Println("Lens", len(seg1), len(seg2))


//	fmt.Println(segm)
//	fmt.Println(segm2)
	//crossover(&img, pop[0], pop[1])

}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 1, 2)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts)

	for id, solution := range solutions {
		graph := GenoToGraph(img, solutions[0].genotype, true)
		segments := graph.ConnectedComponents()
		fmt.Println("Solution", id, ": segments:", len(segments), ", c:", solution.connectivity, ", d:", solution.deviation)
	}
	graph := GenoToGraph(img, solutions[0].genotype, true)

	//visualizeImageGraph("mstgraph.png", img, graph)

	edgedImg := DrawImageBoundries(img, graph, color.Black)
	SaveJPEGRaw(edgedImg, "edges.jpg")
}

func runAndStoreImagesForTesting(folderId string, generations, popSize int) {
	imagePath := "./data/" + folderId + "/Test_image.jpg"
	image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())
	setObjectivesMaxMinValues(&image)

	fmt.Println("Max conn =", maxConnectivity, "Max dev =", maxDeviation)

	solutions := nsgaII(&image, generations, popSize)

	for id, solution := range solutions {
		graph := GenoToGraph(&image, solutions[id].genotype, false)
		segments := graph.ConnectedComponents()
		fmt.Println("Solution", id, ": weightedSum:", solution.weightedSum(), ", segments:", len(segments), ", c:", solution.connectivity, ", d:", solution.deviation)
	}

	dir, err := ioutil.ReadDir("./solutions/Student_Segmentation_Files")

	if err != nil {
		panic(err)
	}

	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}

	for i, s := range solutions {
		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)
		//fmt.Println("Storing solution", s.weightedSum(), filename)
		drawSolutionSegmentsBorders(&image, s, color.Black, filename)
	}
}

func runNSGAOnTestFolder(folderId string) {
	imagePath := "./data/" + folderId + "/colors.jpg"
	image := readJPEGFile(imagePath)
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	// Set max and min connectivity and deviation
	setObjectivesMaxMinValues(&image)

	solutions := nsgaII(&image, 10, 30)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")
	fmt.Println("\nSolutions:")

	for id, s := range solutions {
		graph := GenoToGraph(&image, s.genotype, true)
		segments := graph.ConnectedComponents()
		fmt.Println("Solution", id, ": segments:", len(segments), ", c:", s.connectivity, ", d:", s.deviation)

		if len(segments) > 1 {

			drawSolutionSegmentsBorders(&image, s, color.Black, "border.png")
			drawSolutionSegmentsWithCentroidColor(&image, s, "segments.png")
		}
	}
}

func runGenerations(img *Image) {

	pop := generatePopulation(img, 2)
	return
	//sol := BestSolution(pop)
	//graph := GenoToGraph(img, sol.genotype)

	//visualizeImageGraph("graph.png", img, graph)

	for i := 0; i < 1; i++ {
		pop = RunGeneration(img, pop)
		sol := BestSolution(pop)

		graph := GenoToGraph(img, sol.genotype, false)
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
