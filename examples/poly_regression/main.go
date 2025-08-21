package main

import (
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/examples"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
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
	mainSeries := dataset.ContinuousSeries{
		Name: "A test series",
		// Generates a []float64 from 1.0 to 100.0 in 1.0 step increments, or 100 elements.
		XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		// Generates a []float64 randomly from 0 to 100 with 100 elements.
		YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(),
	}

	polyRegSeries := &dataset.PolynomialRegressionSeries{
		Degree:      3,
		InnerSeries: mainSeries,
	}

	chart := &unichart.Chart{
		Series: []dataset.Series{
			mainSeries,
			polyRegSeries,
		},
	}
	chart.SetHeight(300)

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
