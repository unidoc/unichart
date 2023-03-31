package main

import (
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}
}

func main() {
	// Create chart component.
	bar := &unichart.ProgressBar{
		Background: render.Style{
			FillColor:   render.ColorBlue,
			StrokeWidth: 1.0,
			StrokeColor: render.ColorRed,
		},
		Foreground: render.Style{
			FillColor: render.ColorAlternateGreen,
		},

		RoundedEdgeStart: true,
		RoundedEdgeEnd:   true,
	}
	bar.SetHeight(20)
	bar.SetProgress(0.68)

	// Create unipdf chart component.
	c := creator.New()
	chartComponent := creator.NewChart(bar)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}