package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type SaveImage interface {
	Save(*os.File, image.Image) error
}

var supportedFormats = make(map[string]NewManipulator)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	supportedFormats["jpeg"] = NewJPEGManipulator
	supportedFormats["png"] = NewPNGManipulator
	supportedFormats["gif"] = NewGIFManipulator

	var imgPath string
	var targetWidth, targetHeight, targetPercent uint
	var ratioWidth, ratioHeight bool

	flag.StringVar(&imgPath, "p", "./", "Path to resize images from")
	flag.UintVar(&targetWidth, "w", 1152, "Resized width of image")
	flag.UintVar(&targetHeight, "h", 768, "Resized height of image")
	flag.BoolVar(&ratioWidth, "rw", false, "Use only width input, and height will be calculated maintaining ratio")
	flag.BoolVar(&ratioHeight, "rh", false, "Use only height input, and height will be calculated maintaining ratio")
	flag.UintVar(&targetPercent, "percent", 100, "Resizes the image by a percentage, maintianing ratio")
	flag.Parse()

	if ratioHeight && ratioWidth {
		log.Fatalln("Cannot use both calculated ratio for height and width")
	}
	if (ratioHeight || ratioWidth) && targetPercent != 100 {
		log.Fatalln("Cannot use both ratio Height/Width and percentage")
	}

	fileInfos, err := ioutil.ReadDir(imgPath)
	if err != nil {
		log.Fatalln("Failed to read directory", imgPath, "because", err)
	}

	for _, v := range fileInfos {
		// No need to resize the image if it was already done.
		if strings.Contains(v.Name(), "-resized-") || v.IsDir() {
			continue
		}

		manipulator, err := loadImage(filepath.Join(imgPath, v.Name()))
		if err != nil {
			log.Println("Image Open failed:", v.Name(), err)
			continue
		}

		// Resize the image
		if targetPercent != 100 {
			// Calculates the target width and height for the given percentage resize.
			ratioWidth = true
			targetWidth = uint(float64(manipulator.Bounds().Max.X) * (float64(targetPercent) / float64(100)))
		}

		if ratioWidth {
			// Calculates the height from a target width maintaing the image's ration;
			// e.g. height/width * targetWidth = matchingHeight
			targetHeight = uint((float64(manipulator.Bounds().Max.Y) / float64(manipulator.Bounds().Max.X)) * float64(targetWidth))
		} else if ratioHeight {
			// Calculates the width from a target height maintaing the image's ration;
			// e.g. width/height * targetHeight = matchingWidth
			targetWidth = uint((float64(manipulator.Bounds().Max.X) / float64(manipulator.Bounds().Max.Y)) * float64(targetHeight))
		}
		manipulator.Resize(targetWidth, targetHeight)

		// Save the image
		if err := saveImage(imgPath, v.Name(), manipulator); err != nil {
			log.Println("Failed to save resized image:", err)
			continue
		}

		log.Println("Resized", v.Name(), "to", manipulator.Bounds().Max.String(), "format", manipulator.Format())
	}
}

func loadImage(filePath string) (Manipulator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open image: %s", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode image: %s", err)
	}

	manipGenerator, supported := supportedFormats[format]
	if !supported {
		return nil, fmt.Errorf("format, %s, not supported", format)
	}

	return manipGenerator(file, format, img)
}

func saveImage(imgPath, origFilename string, manipulator Manipulator) error {
	filename := makeNewFilName(origFilename, manipulator.Bounds())
	file, err := os.Create(filepath.Join(imgPath, filename))
	if err != nil {
		log.Println("Failed to create resized file:", err)
	}
	return manipulator.Save(file)
}

func makeNewFilName(origName string, bounds image.Rectangle) string {
	extention := filepath.Ext(origName)
	filename := strings.TrimSuffix(origName, extention)
	return fmt.Sprintf("%s-resized-%dx%d%s", filename, bounds.Max.X, bounds.Max.Y, extention)
}
