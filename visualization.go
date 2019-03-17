package main

import (
	"github.com/alonsovidales/go_graph"
	"github.com/fogleman/gg"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

func visualizeImageGraph(filename string, img *Image, graph *graphs.Graph) {

	width := len(*img)
	imageWidth := 20 * width
	height := len((*img)[0])
	imageHeight := 20 * height

	dc := gg.NewContext(imageWidth, imageHeight)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			dc.DrawCircle(float64(20*x)+10, float64(20*y)+10, 5)
		}
	}
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(3)
	dc.Stroke()

	for _, edge := range graph.RawEdges {

		fromX := edge.From % uint64(width)
		fromY := edge.From / uint64(width)
		toX := edge.To % uint64(width)
		toY := edge.To / uint64(width)

		dc.DrawLine(float64(20*fromX+10), float64(20*fromY+10), float64(20*toX+10), float64(20*toY+10))
	}

	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.Stroke()

	dc.SavePNG(filename)
}

func createParetoPlotter() *plot.Plot {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Fronts"
	p.X.Label.Text = "Deviation"
	p.Y.Label.Text = "Connectivity"

	return p
}

func addParetoFrontToPlotter(p *plot.Plot, population []*Solution, fronts map[int][]int, generation int) {

	if generation % 10 != 0 {
		return
	}

	pts := make(plotter.XYs, len(fronts[0]))

	for i := range fronts[0] {
		pts[i].X = population[fronts[0][i]].deviation / (maxDeviation - minDeviation)
		pts[i].Y = population[fronts[0][i]].connectivity / (maxConnectivity - minConnectivity)
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(pts)

	if err != nil {
		panic(err)
	}

	lpPoints.Shape = draw.CircleGlyph{}
	lpLine.Color = color.RGBA{A: 0}

	lpPoints.Color = color.RGBA{B: uint8(255 * generation / generationsToRun), R: 255 - uint8(255*generation/generationsToRun), A: 255}

	p.Add(lpLine, lpPoints)
}

func saveParetoPlotter(p *plot.Plot, name string) {
	if err := p.Save(10*vg.Inch, 10*vg.Inch, name); err != nil {
		panic(err)
	}
}

func visualizeFronts(population []*Solution, fronts map[int][]int, name string) {
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

	if err := p.Save(10*vg.Inch, 10*vg.Inch, name); err != nil {
		panic(err)
	}
}
