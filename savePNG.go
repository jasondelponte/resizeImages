package main

import (
	"image"
	"image/png"
	"os"
)

type SavePNGImage struct{}

func (s SavePNGImage) Save(file *os.File, img image.Image) error {
	return png.Encode(file, img)
}
