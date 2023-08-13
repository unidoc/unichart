package examples

import (
	"log"
	"os"

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
}
