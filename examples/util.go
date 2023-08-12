package examples

import (
	"fmt"
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

	// Get total number of pages.
	numPages, err := reader.GetNumPages()
	if err != nil {
		log.Fatalf("Could not retrieve number of pages: %v", err)
	}

	// Render pages.
	device := render.NewImageDevice()
	device.OutputWidth = 2000
	for i := 1; i <= numPages; i++ {
		// Get page.
		page, err := reader.GetPage(i)
		if err != nil {
			log.Fatalf("Could not retrieve page: %v", err)
		}

		// Render page to PNG file.
		imgFileName := fmt.Sprintf("preview_%d.png", i)

		err = device.RenderToPath(page, imgFileName)
		if err != nil {
			log.Fatalf("Image rendering error: %v", err)
		}
	}
}
