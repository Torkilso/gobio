package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"image/jpeg"
	"os"
)

func visualizeFronts(population []*Solution, fronts map[int][]int) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Fronts"
	p.X.Label.Text = "Deviation"
	p.Y.Label.Text = "Connectivity"

	for front, ids := range fronts {
		pts := make(plotter.XYs, len(ids))

		for i := range ids {
			pts[i].X = population[ids[i]].deviation
			pts[i].Y = population[ids[i]].connectivity
		}

		lpLine, lpPoints, err := plotter.NewLinePoints(pts)

		if err != nil {
			panic(err)
		}

		lpPoints.Shape = draw.CircleGlyph{}
		lpLine.Color = color.RGBA{A: 0}

		if front == 0 {
			lpPoints.Color = color.RGBA{B: 255, A: 255}
			//p.Legend.Add("Pareto front", lpLine, lpPoints)
		} else {
			lpPoints.Color = color.RGBA{R: 255, A: 255}
		}

		p.Add(lpLine, lpPoints)
	}

	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func readJPEGFile(path string) Image {
	infile, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}

	defer infile.Close()

	src, err := jpeg.Decode(infile)

	if err != nil {
		panic(err.Error())
	}

	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	pixels := make([][]Pixel, width)

	for i := range pixels {
		pixels[i] = make([]Pixel, height)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := src.At(i, j).RGBA()

			red := int16(r >> 8)
			green := int16(g >> 8)
			blue := int16(b >> 8)

			pixels[i][j] = Pixel{
				r: red,
				g: green,
				b: blue,
			}
		}
	}

	return pixels
}
