package main

import (
	"github.com/alonsovidales/go_graph"
	"github.com/google/gxui"
	"github.com/google/gxui/themes/dark"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageDrawer func(img *image.RGBA)

func GenerateImageDrawer(driver gxui.Driver, rect image.Rectangle) ImageDrawer {
	theme := dark.CreateTheme(driver)
	window := theme.CreateWindow(rect.Max.X, rect.Max.Y, "Image viewer")
	window.OnClose(driver.Terminate)
	return func(rgba *image.RGBA) {
		f, err := os.Create("img.jpg")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		jpeg.Encode(f, rgba, nil)

		img := theme.CreateImage()
		window.RemoveAll()
		texture := driver.CreateTexture(rgba, 1.0)
		img.SetTexture(texture)
		window.AddChild(img)
	}
}

func SaveJPEG(img *Image) {
	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, (*img).toRGBA(), nil)
}

func SaveJPEGRaw(img *image.RGBA, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

func SavePNGRaw(img *image.RGBA, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func DrawImageBoundries(img *Image, gr *graphs.Graph, color color.Color) *image.RGBA {
	res := img.toRGBA()
	groups := gr.ConnectedComponents()
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

					res.Set(x1, y1, color)
					res.Set(x2, y2, color)
				}
			}
		}
	}
	return res
}

func drawSolutionSegmentsBorders(img *Image, solution *Solution, col color.Color, name string) {

	imgRect := image.Rect(0, 0, len(*img), len((*img)[0]))
	blank := image.NewRGBA(imgRect)

	for x := imgRect.Min.X ; x < imgRect.Max.X ; x++ {
		for y := imgRect.Min.Y ; y < imgRect.Max.Y ; y++ {
			blank.Set(x, y, color.White)
		}
	}

	gr := GenoToGraph(img, solution.genotype)
	groups := gr.ConnectedComponents()

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

	gr := GenoToGraph(img, solution.genotype)
	groups := gr.ConnectedComponents()

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
