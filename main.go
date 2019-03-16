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
	runMultiObjective(folderId, generationsToRun, popSize)
	//runNSGAOnTestFolder("216066")

}

func initialize(folderId string) *Image {

	imagePath := "./data/" + folderId + "/Test image.jpg"
	image := readJPEGFile(imagePath)

	rand.Seed(time.Now().UTC().UnixNano())
	setObjectivesMaxMinValues(&image)

	return &image
}

func cleanTestingDirs(folderId string) {
	dir, err := ioutil.ReadDir("./solutions/Student_Segmentation_Files")
	optimalDir, err := ioutil.ReadDir("./solutions/Optimal_Segmentation_Files")
	dataDir, err := ioutil.ReadDir("./data/" + folderId + "/")

	if err != nil {
		panic(err)
	}

	if err2 != nil {
		panic(err2)
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
}
func runMultiObjective(folderId string, generations, popSize int) {
	cleanTestingDirs(folderId)
	image := initialize(folderId)

	solutions := nsgaII(image, generations, popSize)

	for i, s := range solutions {
		segments := GenoToConnectedComponents(s.genotype)
		if len(segments) > 500 || len(segments) < 2 {
			continue
		}
		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)

		drawSolutionSegmentsBorders(image, s, color.Black, filename)
	}
}



func runSingleObjective(folderId string, generations, popSize int) {
	cleanTestingDirs(folderId)
	image := initialize(folderId)

	solutions := singleObjective(image, generations, popSize)

	for i, s := range solutions {
		segments := GenoToConnectedComponents(s.genotype)
		if len(segments) > 500 || len(segments) < 2 {
			continue
		}
		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)

		drawSolutionSegmentsBorders(image, s, color.Black, filename)
	}
	pop := generatePopulation(img, 2)

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
