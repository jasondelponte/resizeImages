package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	supportedFormats := make(map[string]struct{})
	supportedFormats["jpeg"] = struct{}{}
	supportedFormats["png"] = struct{}{}

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

		img, format, err := openImage(filepath.Join(imgPath, v.Name()))
		if err != nil {
			log.Println("Image Open failed:", v.Name(), err)
			continue
		}

		if _, supported := supportedFormats[format]; !supported {
			log.Println("Unsupported format:", format, v.Name())
			continue
		}

		if targetPercent != 100 {
			// Calculates the target width and height for the given percentage resize.
			ratioWidth = true
			targetWidth = uint(float64(img.Bounds().Max.X) * (float64(targetPercent) / float64(100)))
		}

		if ratioWidth {
			// Calculates the height from a target width maintaing the image's ration;
			// e.g. 1200/1600 * 400 = 300
			targetHeight = uint((float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X)) * float64(targetWidth))
		} else if ratioHeight {
			// Calculates the width from a target height maintaing the image's ration;
			// e.g. 1600/1200 * 300 = 400
			targetWidth = uint((float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y)) * float64(targetHeight))
		}

		resizedImg := resize.Resize(targetWidth, targetHeight, img, resize.NearestNeighbor)

		newFilename := makeNewFilName(v.Name(), resizedImg.Bounds())
		newFilePath := filepath.Join(imgPath, newFilename)
		if err := saveImage(newFilePath, resizedImg, format); err != nil {
			log.Println("Failed to save resized image", err)
		}

		log.Println("Resized", v.Name(), "to", newFilename, "format", format)
	}
}

func openImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to open image: %s", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to decode image: %s", err)
	}

	return img, format, nil
}

func saveImage(filePath string, img image.Image, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create resized image file: %s", err)
	}
	defer file.Close()

	if format == "jpeg" {
		if err = saveJPEGImage(file, img); err != nil {
			return fmt.Errorf("Failed writing image as JPEG to file: %s", err)
		}
	} else if format == "png" {
		if err = savePNGImage(file, img); err != nil {
			return fmt.Errorf("Failed writing image as PNG to file: %s", err)
		}
	}

	return nil
}

func makeNewFilName(origName string, bounds image.Rectangle) string {
	extention := filepath.Ext(origName)
	filename := strings.TrimSuffix(origName, extention)
	return fmt.Sprintf("%s-resized-%dx%d%s", filename, bounds.Max.X, bounds.Max.Y, extention)
}
