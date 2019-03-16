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
	generationsToRun = 100
	popSize          = 20
	folderId         = "86016"
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
	runAndStoreImagesForTesting(folderId, generationsToRun, popSize)
	//runNSGAOnTestFolder("216066")
	//img := readJPEGFile("./testimages/Untitled2.jpg")
	//testMaxObjectives(&img)
	//perf()

}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 1, 2)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts, "final_pareto.png")

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
	imagePath := "./data/" + folderId + "/image.jpg"
	image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())
	setObjectivesMaxMinValues(&image)

	fmt.Println("Max conn =", maxConnectivity, "Max dev =", maxDeviation)

	solutions := nsgaII(&image, generations, popSize)

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts, "final_pareto.png")

	for id, solution := range solutions {
		segments := GenoToConnectedComponents(solutions[id].genotype)
		fmt.Println("Solution", id, ": weightedSum:", solution.weightedSum(), ", segments:", len(segments), ", c:", solution.connectivity, ", d:", solution.deviation)
	}

	dir, err := ioutil.ReadDir("./solutions/Student_Segmentation_Files")
	dir2, err2 := ioutil.ReadDir("./solutions/Solutions_With_Image")

	if err != nil {
		panic(err)
	}

	if err2 != nil {
		panic(err2)
	}

	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}

	for _, d := range dir2 {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Solutions_With_Image", d.Name()}...))
	}

	for i, s := range solutions {
		segments := GenoToConnectedComponents(solutions[i].genotype)
		if len(segments) > 2000 || len(segments) < 2 {
			continue
		}

		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)
		filename2 := fmt.Sprintf("./solutions/Solutions_With_Image/sol%d.jpg", i)
		//fmt.Println("Storing solution", s.weightedSum(), filename)
		drawSolutionSegmentsBorders(&image, s, color.Black, filename)
		drawSolutionSegmentsBordersWithImage(&image, s, color.RGBA{G: 255, A: 0xff}, filename2)
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
