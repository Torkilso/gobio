package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
)

func GoImageToImage(src image.Image) Image {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	pixels := make([][]*Pixel, width)

	for i := range pixels {
		pixels[i] = make([]*Pixel, height)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := src.At(i, j).RGBA()

			red := int16(r >> 8)
			green := int16(g >> 8)
			blue := int16(b >> 8)

			pixels[i][j] = &Pixel{
				r: red,
				g: green,
				b: blue,
			}
		}
	}
	return pixels
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

	return GoImageToImage(src)
}

func copyTo(sourceFile, destinationFile string) {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}

}
