package main

import (
	"embed"
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed soyken.png
var embedFS embed.FS

func main() {
	flag.Parse()

	path := flag.Arg(0)

	fInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if fInfo.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			if e.IsDir() {
				log.Printf("skipping %s, it is a directory\n", e.Name())
				continue
			}

			if strings.Contains(e.Name(), ".soy.") {
				continue
			}

			fullpath := filepath.Join(path, e.Name())

			if !isImage(fullpath) {
				log.Printf("skipping %s, it is not an image\n", e.Name())
				continue
			}

			err := soyifyImage(fullpath)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("soyken'd all images in %s successfully\n", path)
		return
	}

	if !isImage(path) {
		log.Fatalf("%s is not an image", path)
	}

	err = soyifyImage(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("soyken'd %s successfully", path)
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func outputPath(path string) string {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)

	base := filename[:len(filename)-len(ext)]

	outputFilename := base + ".soy" + ext

	return filepath.Join(dir, outputFilename)
}

func resizeImage(img image.Image, factor float64) image.Image {
	width := int(float64(img.Bounds().Dx()) * factor)
	height := int(float64(img.Bounds().Dy()) * factor)

	resizedImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origX := int(float64(x) / factor)
			origY := int(float64(y) / factor)
			resizedImage.Set(x, y, img.At(origX, origY))
		}
	}

	return resizedImage
}

func soyifyImage(path string) error {
	bgFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer bgFile.Close()

	fgFile, err := embedFS.Open("soyken.png")
	if err != nil {
		return err
	}
	defer fgFile.Close()

	bgImage, _, err := image.Decode(bgFile)
	if err != nil {
		return err
	}

	fgImage, _, err := image.Decode(fgFile)
	if err != nil {
		return err
	}

	bgWidth := bgImage.Bounds().Dx()
	bgHeight := bgImage.Bounds().Dy()
	maxHeight := int(0.8 * float64(bgHeight))
	scaleFactor := 1.0
	if fgImage.Bounds().Dy() > maxHeight {
		scaleFactor = float64(maxHeight) / float64(fgImage.Bounds().Dy())
	}

	fgResized := resizeImage(fgImage, scaleFactor)

	posX := bgWidth - fgResized.Bounds().Dx()
	posY := bgHeight - fgResized.Bounds().Dy()

	canvas := image.NewRGBA(bgImage.Bounds())

	draw.Draw(canvas, canvas.Bounds(), bgImage, image.Pt(0, 0), draw.Src)

	draw.Draw(canvas, fgResized.Bounds().Add(image.Pt(posX, posY)), fgResized, image.Pt(0, 0), draw.Over)

	outPath := outputPath(path)

	output, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer output.Close()

	err = jpeg.Encode(output, canvas, nil)
	if err != nil {
		return err
	}

	return nil
}
