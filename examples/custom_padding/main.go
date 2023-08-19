package main

import (
	"image/color"
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/examples"
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
	chart := &unichart.Chart{
		Background: render.Style{
			Padding: render.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: color.RGBA{R: 239, G: 239, B: 239, A: 239},
		},
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
				YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(100).WithMin(100).WithMax(512)}.Values(),
			},
		},
	}
	chart.SetHeight(400)

	// Create unipdf chart component.
	c := creator.New()
	chartComponent := creator.NewChart(chart)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	examples.RenderPDFToImage("output.pdf")
}
