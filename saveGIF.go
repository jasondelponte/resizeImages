package main

import (
	"image"
	"image/gif"
	"os"
)

type SaveGIFImage struct{}

func (s SaveGIFImage) Save(file *os.File, img image.Image) error {
	return gif.Encode(file, img, nil)
}
