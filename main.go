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
	generationsToRun = 300
		popSize          = 200
	folderId         = "176035"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	//defer profile.Start(profile.MemProfile).Stop()
	runMultiObjective(folderId, generationsToRun, popSize)
	//runSingleObjective(folderId, generationsToRun, popSize)

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
	greenBordersDir, err := ioutil.ReadDir("./solutions/Solutions_With_Image")

	if err != nil {
		panic(err)
	}

	for _, d := range dir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}
	for _, d := range optimalDir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Optimal_Segmentation_Files", d.Name()}...))
	}
	for _, d := range greenBordersDir {
		_ = os.RemoveAll(path.Join([]string{"./solutions/Solutions_With_Image", d.Name()}...))
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

	joinSegmentsStart := time.Now()
	solutions.joinSegments(image, 100)
	fmt.Print("Used ", time.Since(joinSegmentsStart).Seconds(), " to join segments\n\n")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts, "final_pareto.png")

	fmt.Println("\nSolutions:")

	for i, s := range solutions {
		segments := GenoToConnectedComponents(s.genotype)
		fmt.Println("segments:", len(segments))

		if len(segments) > 50 || len(segments) < 2 {
			continue
		}

		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)
		filenameGreen := fmt.Sprintf("./solutions/Solutions_With_Image/sol%d.jpg", i)
		filenameFill := fmt.Sprintf("./solutions/Solutions_With_Fill/sol%d.jpg", i)


		drawSolutionSegmentsBorders(image, s, color.Black, filename)
		drawSolutionSegmentsBordersWithImage(image, s, color.RGBA{G: 255, A: 0xff}, filenameGreen)
		drawSolutionSegmentsWithCentroidColor(image, s, filenameFill)

	}
}

func runSingleObjective(folderId string, generations, popSize int) {
	cleanTestingDirs(folderId)
	image := initialize(folderId)

	bestSolution := singleObjective(image, generations, popSize)

	filename := "./solutions/Student_Segmentation_Files/sol.jpg"

	drawSolutionSegmentsBorders(image, bestSolution, color.Black, filename)
}
