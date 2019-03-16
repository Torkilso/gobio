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
	popSize          = 80
	folderId         = "178054"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	//defer profile.Start(profile.MemProfile).Stop()
	//runMultiObjective(folderId, generationsToRun, popSize)
	runSingleObjective(folderId, generationsToRun, popSize)

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

	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}
	for _, d := range optimalDir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Optimal_Segmentation_Files", d.Name()}...))
	}
	for _, file := range dataDir {
		if strings.Contains(file.Name(), "GT") {
			copyTo("./data/"+folderId+"/"+file.Name(), "./solutions/Optimal_Segmentation_Files/"+file.Name())
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

/*
func runSingleObjective(folderId string, generations, popSize int) {
	cleanTestingDirs(folderId)
	image := initialize(folderId)

	bestSolution := singleObjective(image, generations, popSize)

	filename := "./solutions/Student_Segmentation_Files/sol.jpg"

	drawSolutionSegmentsBorders(image, bestSolution, color.Black, filename)
}
*/
