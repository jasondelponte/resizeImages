package main

import (
	"image"
	"image/png"
	"os"
)

func savePNGImage(file *os.File, img image.Image) error {
	return png.Encode(file, img)
}
