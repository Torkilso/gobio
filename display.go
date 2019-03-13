package main

import (
	"fmt"
	"github.com/alonsovidales/go_graph"
	"github.com/google/gxui"
	"github.com/google/gxui/themes/dark"
	"image"
	"image/color"
	"image/jpeg"
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

func drawImageSegmentsWithCentroidColor(img *Image, gr *graphs.Graph) *image.RGBA {
	res := img.toRGBA()
	groups := gr.ConnectedComponents()
	width := len(*img)
	for id, group := range groups {
		fmt.Println(id)
		centroid := Centroid(img, group)
		for k := range group {
			intK := int(k)

			x1, y1 := Flatten(width, intK)

			color := color.RGBA{R: uint8(centroid.r), G: uint8(centroid.g), B: uint8(centroid.b)}

			res.Set(x1, y1, color)

		}
	}
	return res
}
