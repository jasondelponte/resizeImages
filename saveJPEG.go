package main

import (
	"image"
	"image/jpeg"
	"os"
)

func saveJPEGImage(file *os.File, img image.Image) error {
	return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
}
