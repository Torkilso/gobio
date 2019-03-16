package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

var (
	generationsToRun = 120
	popSize = 80
	folderId         = "176035"
)

func main() {

	//imagePath := "./data/216066/Test_image.jpg"
	//imagePath := "./testimages/Untitled2.jpg"
	//image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())

	//defer profile.Start(profile.MemProfile).Stop()

	//runGenerations(&image)
	//runNSGA(&image)
	runAndStoreImagesForTesting(folderId, generationsToRun, popSize)
	//runNSGAOnTestFolder("216066")

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
	imagePath := "./data/" + folderId + "/Test image.jpg"
	image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())
	setObjectivesMaxMinValues(&image)

	fmt.Println("Max conn =", maxConnectivity, "Max dev =", maxDeviation, "Min edge =", minEdgeValues)

	solutions := nsgaII(&image, generations, popSize)

	for id, solution := range solutions {
		segments := GenoToConnectedComponents(solutions[id].genotype)
		fmt.Println("Solution", id, ": weightedSum:", solution.weightedSum(), ", segments:", len(segments), ", c:", solution.connectivity, ", d:", solution.deviation)
	}

	dir, err := ioutil.ReadDir("./solutions/Student_Segmentation_Files")
	optimalDir, err := ioutil.ReadDir("./solutions/Optimal_Segmentation_Files")
	dataDir, err := ioutil.ReadDir("./data/" + folderId + "/")

	if err != nil {
		panic(err)
	}

	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}
	for _, d := range optimalDir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Optimal_Segmentation_Files", d.Name()}...))
	}
	for _, file := range dataDir {
		if strings.Contains(file.Name(), "GT") {
			copyTo("./data/" + folderId + "/" + file.Name(), "./solutions/Optimal_Segmentation_Files/" + file.Name())
		}
	}

	for i, s := range solutions {
		segments := GenoToConnectedComponents(s.genotype)
		if len(segments) > 500 || len(segments) < 2 {
			continue
		}
		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)

		drawSolutionSegmentsBorders(&image, s, color.Black, filename)
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
