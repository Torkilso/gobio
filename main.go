package main

func main() {
	imagePath := "./Test Images Project 2/216066/Test image.jpg"
	image := readJPEGFile(imagePath)

	solutions := nsgaII(&image, 100, 100)
	_ = solutions

}
/*

func runGenerations(){
	pop := GeneratePopulation(&src, 4)

	for i := 0 ; i < 10 ; i++ {
		pop = RunGeneration(&src, pop)
		fmt.Println("Generation done")
		rgba2 := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		fmt.Println("Generating", src.Bounds().Max, b.Dx(), b.Dy())
		sol := pop.solutions[0]
		groups := sol.graph.ConnectedComponents()
		fmt.Println("Num groups", len(groups))
		maxX := 0
		maxY := 0
		for _, g := range groups {
			c := Centroid(&src, g)
			for k := range g {
				x, y := Flatten(rect, int(k))
				rgba2.Set(x, y, c)
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}

		fmt.Println("Drawing", maxX, maxY)
		imageDrawer(rgba2)
	}
	imageDrawer(rgba)

}
*/