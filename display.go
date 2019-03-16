package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)


func SavePNGRaw(img *image.RGBA, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}


func drawSolutionSegmentsBorders(img *Image, solution *Solution, col color.Color, name string) {

	imgRect := image.Rect(0, 0, len(*img), len((*img)[0]))
	blank := image.NewRGBA(imgRect)

	for x := imgRect.Min.X; x < imgRect.Max.X; x++ {
		for y := imgRect.Min.Y; y < imgRect.Max.Y; y++ {
			blank.Set(x, y, color.White)
		}
	}

	groups := GenoToConnectedComponents(solution.genotype)

	width := len(*img)
	for _, group := range groups {
		for k := range group {
			intK := int(k)
			for _, neighbour := range GetTargets(img, intK) {
				if _, ok := group[uint64(neighbour)]; ok { // Same segments
				} else {
					// Two neighbours are not in same segment.
					// Draw edge in this and neighbour
					x1, y1 := Flatten(width, intK)
					x2, y2 := Flatten(width, neighbour)

					blank.Set(x1, y1, col)
					blank.Set(x2, y2, col)
				}
			}
		}
	}
	SavePNGRaw(blank, name)
}

func drawSolutionSegmentsBordersWithImage(img *Image, solution *Solution, col color.Color, name string) {
	imgRect := image.Rect(0, 0, len(*img), len((*img)[0]))
	blank := image.NewRGBA(imgRect)

	for x := imgRect.Min.X; x < imgRect.Max.X; x++ {
		for y := imgRect.Min.Y; y < imgRect.Max.Y; y++ {

			R := (*img)[x][y].r
			G := (*img)[x][y].g
			B := (*img)[x][y].b

			imageColor := color.RGBA{R: uint8(R), G: uint8(G), B: uint8(B), A: 0xff}
			blank.Set(x, y, imageColor)
		}
	}


	groups := GenoToConnectedComponents(solution.genotype)

	width := len(*img)
	for _, group := range groups {
		for k := range group {
			intK := int(k)
			for _, neighbour := range GetTargets(img, intK) {
				if _, ok := group[uint64(neighbour)]; ok { // Same segments
				} else {
					// Two neighbours are not in same segment.
					// Draw edge in this and neighbour
					x1, y1 := Flatten(width, intK)
					x2, y2 := Flatten(width, neighbour)

					blank.Set(x1, y1, col)
					blank.Set(x2, y2, col)
				}
			}
		}
	}
	SavePNGRaw(blank, name)
}

func drawSolutionSegmentsWithCentroidColor(img *Image, solution *Solution, name string) {

	imgRect := image.Rect(0, 0, len(*img), len((*img)[0]))
	blank := image.NewRGBA(imgRect)

	groups := GenoToConnectedComponents(solution.genotype)

	width := len(*img)
	for _, group := range groups {
		//fmt.Println("segment id: ", id, " group length:", len(group))
		centroid := Centroid(img, group)
		//fmt.Println("centroid", centroid)
		for k := range group {
			intK := int(k)

			x1, y1 := Flatten(width, intK)

			//fmt.Println(x1, y1)
			centroidColor := color.RGBA{R: uint8(centroid.r), G: uint8(centroid.g), B: uint8(centroid.b), A: 0xFF}

			blank.Set(x1, y1, centroidColor)

		}
	}
	SavePNGRaw(blank, name)
}
