package main

import (
	"github.com/nfnt/resize"
	"image"
	"io"
)

type Manipulator interface {
	Format() string
	Bounds() image.Rectangle
	Resize(uint, uint)
	Save(io.Writer) error
}

type NewManipulator func(io.ReadSeeker, string, image.Image) (Manipulator, error)

type Saver func(io.Writer, image.Image) error

type GenericManipulator struct {
	image  image.Image
	saver  Saver
	format string
}

func NewGenericManipulator(image image.Image, format string, saver Saver) Manipulator {
	return &GenericManipulator{
		image:  image,
		saver:  saver,
		format: format,
	}
}

func (m *GenericManipulator) Format() string {
	return m.format
}

func (m *GenericManipulator) Bounds() image.Rectangle {
	return m.image.Bounds()
}

func (m *GenericManipulator) Resize(width, height uint) {
	m.image = resize.Resize(width, height, m.image, resize.NearestNeighbor)
}

func (m *GenericManipulator) Save(writer io.Writer) error {
	return m.saver(writer, m.image)
}
