package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/gif"
	"io"
)

type GIFManipulator struct {
	gif *gif.GIF
}

func NewGIFManipulator(reader io.ReadSeeker, format string, image image.Image) (Manipulator, error) {
	if _, err := reader.Seek(0, 0); err != nil {
		return nil, err
	}

	g, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	return &GIFManipulator{gif: g}, nil
}

func (m *GIFManipulator) Format() string {
	return "gif"
}

func (m *GIFManipulator) Bounds() image.Rectangle {
	if len(m.gif.Image) == 0 {
		return image.Rectangle{}
	}

	return m.gif.Image[0].Bounds()
}

func (m *GIFManipulator) Resize(width, height uint) {
	for i, v := range m.gif.Image {
		resizedImg := resize.Resize(width, height, v, resize.NearestNeighbor)
		m.gif.Image[i] = imageToPalleted(resizedImg, m.gif.Image[i].Palette)
	}
}

func (m *GIFManipulator) Save(writer io.Writer) error {
	return gif.EncodeAll(writer, m.gif)
}

func imageToPalleted(img image.Image, palette color.Palette) *image.Paletted {
	palleted := image.NewPaletted(img.Bounds(), palette)

	for x := palleted.Bounds().Min.X; x < palleted.Bounds().Max.X; x++ {
		for y := palleted.Bounds().Min.Y; y < palleted.Bounds().Max.Y; y++ {
			palleted.Set(x, y, img.At(x, y))
		}
	}

	return palleted
}
