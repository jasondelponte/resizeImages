package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var imgPath string
	var targetWidth, targetHeight uint

	flag.StringVar(&imgPath, "p", "./", "Path to resize images from")
	flag.UintVar(&targetWidth, "w", 1152, "Resized width of image")
	flag.UintVar(&targetHeight, "h", 768, "Resized height of image")
	flag.Parse()

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
			log.Println("Failed to resize", v.Name(), err)
			continue
		}

		if format != "jpeg" {
			log.Println("Unsupported format", format, v.Name())
			continue
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
		log.Println("Failed to open image")
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Println("Failed to decode image")
		return nil, "", err
	}

	return img, format, nil
}

func saveImage(filePath string, img image.Image, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to create resized image file")
		return err
	}
	defer file.Close()

	if format == "jpeg" {
		if err = saveJPEGImage(file, img); err != nil {
			log.Println("Failed writing image as JPEG to file")
			return err
		}
	}

	return nil
}

func makeNewFilName(origName string, bounds image.Rectangle) string {
	extention := filepath.Ext(origName)
	filename := strings.TrimSuffix(origName, extention)
	return fmt.Sprintf("%s-resized-%dx%d%s", filename, bounds.Max.X, bounds.Max.Y, extention)
}
