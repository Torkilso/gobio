package main

func main() {
	imagePath := "./Test Images Project 2/216066/Test image.jpg"
	image := readJPEGFile(imagePath)

	solutions := nsgaII(image, 100, 100)
	_ = solutions

}
