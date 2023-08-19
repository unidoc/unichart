package examples

import (
	"image"
	"image/png"

	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

func RenderPDFToImage(filename string) {
	// Create reader.
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open pdf file: %v", err)
	}
	defer file.Close()

	reader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatalf("Could not create reader: %v", err)
	}

	// Render pages.
	device := render.NewImageDevice()
	device.OutputWidth = 2000

	// Get page.
	page, err := reader.GetPage(1)
	if err != nil {
		log.Fatalf("Could not retrieve page: %v", err)
	}

	// Render page to PNG file.
	err = device.RenderToPath(page, "preview.png")
	if err != nil {
		log.Fatalf("Image rendering error: %v", err)
	}

	cropImage("preview.png")
}

func cropImage(imagePath string) {
	// Open the input image
	img, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// Find the bounding box of the non-empty (foreground) region
	boundingBox := findBoundingBox(img)

	// Crop the image using the bounding box
	croppedImg := imaging.Crop(img, boundingBox)

	// Save the cropped image
	outFile, err := os.Create(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, croppedImg)
	if err != nil {
		log.Fatal(err)
	}
}

func findBoundingBox(img image.Image) image.Rectangle {
	bounds := img.Bounds()
	minX, maxX := bounds.Dx(), 0
	minY, maxY := bounds.Dy(), 0
	boundPadding := 100

	// Iterate over each pixel to find the bounding box
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelColor := img.At(x, y)
			r, g, b, a := pixelColor.RGBA()

			if r != 65535 || g != 65535 || b != 65535 || a != 65535 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Define the bounding box based on the min and max values
	boundingBox := image.Rect(minX-boundPadding, minY-boundPadding, maxX+boundPadding+1, maxY+boundPadding+1)
	return boundingBox
}
