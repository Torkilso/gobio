package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/themes/dark"
	"image"
	"image/jpeg"
	"os"
)

type ImageDrawer func (img *image.RGBA)


func GenerateImageDrawer(driver gxui.Driver, rect image.Rectangle)  ImageDrawer {
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

func SaveJPEGRaw(img *image.RGBA) {
	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

}
