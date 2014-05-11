package main

import (
	"image"
	"image/jpeg"
	"os"
)

type SaveJPEGImage struct{}

func (s SaveJPEGImage) Save(file *os.File, img image.Image) error {
	return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
}
