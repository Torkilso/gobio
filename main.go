package main

import (
	"fmt"
	_image "image"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

func main() {

	//imagePath := "./data/216066/Test_image.jpg"
	//imagePath := "./testimages/Untitled2.jpg"
	//image := readJPEGFile(imagePath)

	//solutions := nsgaII(&image, 100, 100)
	//_ = solutions
	rand.Seed(time.Now().UTC().UnixNano())

	//runGenerations(&image)
	//runNSGA(&image)
	//runSoniaMST(&image)
	runNSGAOnTestFolder("216066", 30, 50)
}

func runNSGA(img *Image) {

	start := time.Now()

	solutions := nsgaII(img, 10, 4)

	fmt.Println("Used", time.Since(start).Seconds(), "seconds in total")

	fronts := fastNonDominatedSort(solutions)
	visualizeFronts(solutions, fronts)

	best := BestSolution(solutions)

	graph := GenoToGraph(img, best.genotype)
	segments := graph.ConnectedComponents()
	fmt.Println("Amount of segments:", len(segments))

	//visualizeImageGraph("mstgraph.png", img, graph)

	//thisImg := img.toRGBA()

	//imgCopy := GoImageToImage(thisImg)

	edgedImg := DrawImageBoundries(img, graph, color.Black)
	SaveJPEGRaw(edgedImg, "edges.jpg")
}

func runNSGAOnTestFolder(folderId string, generations, popSize int) {
	imagePath := "./data/" + folderId + "/Test image.jpg"
	image := readJPEGFile(imagePath)
	rand.Seed(time.Now().UTC().UnixNano())

	solutions := nsgaII(&image, generations, popSize)


	dir, err := ioutil.ReadDir("./solutions/Student_Segmentation_Files")
	if err != nil {
		panic(err)
	}
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"./solutions/Student_Segmentation_Files", d.Name()}...))
	}

	for i, s := range solutions {
		white := _image.NewRGBA(_image.Rect(0, 0, len(image), len(image[0])))
		gr := GenoToGraph(&image, s.genotype)
		goImage := GoImageToImageRGBA(white)
		edgedImg := DrawImageBoundries(&goImage, gr, color.Black)
		filename := fmt.Sprintf("./solutions/Student_Segmentation_Files/sol%d.jpg", i)
		fmt.Println("Filename", filename)
		SaveJPEGRaw(edgedImg, filename)
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

