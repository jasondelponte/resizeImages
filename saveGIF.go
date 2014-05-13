package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"io"
)

type PNGManipulator struct {
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

	return &PNGManipulator{gif: g}, nil
}

func (m *PNGManipulator) Format() string {
	return "gif"
}

func (m *PNGManipulator) Bounds() image.Rectangle {
	if len(m.gif.Image) == 0 {
		return image.Rectangle{}
	}

	return m.gif.Image[0].Bounds()
}

func (m *PNGManipulator) Resize(width, height uint) {
	// for i, v := range m.gif.Image {
	// 	m.gif.Image[i] = resize.Resize(width, height, v, resize.NearestNeighbor)
	// }
}

func (m *PNGManipulator) Save(writer io.Writer) error {
	return gif.EncodeAll(writer, m.gif)
}
