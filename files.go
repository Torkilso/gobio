package main

import (
	"image/jpeg"
	"os"
)

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
