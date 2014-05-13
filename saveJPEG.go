package main

import (
	"image"
	"image/jpeg"
	"io"
)

func NewJPEGManipulator(reader io.ReadSeeker, format string, image image.Image) (Manipulator, error) {
	return NewGenericManipulator(image, format, jpegSaver), nil
}

func jpegSaver(writer io.Writer, img image.Image) error {
	return jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
}
