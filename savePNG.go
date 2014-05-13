package main

import (
	"image"
	"image/png"
	"io"
)

func NewPNGManipulator(reader io.ReadSeeker, format string, image image.Image) (Manipulator, error) {
	return NewGenericManipulator(image, format, pngSaver), nil
}

func pngSaver(writer io.Writer, img image.Image) error {
	return png.Encode(writer, img)
}
