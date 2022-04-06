package main

import (
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	mainSeries := dataset.ContinuousSeries{
		Name: "A test series",
		// Generates a []float64 from 1.0 to 100.0 in 1.0 step increments, or 100 elements.
		XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		// Generates a []float64 randomly from 0 to 100 with 100 elements.
		YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(),
	}

	// NOTE: we create a LinearRegressionSeries series by assigning the inner series.
	linRegSeries := &dataset.LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	chart := &unichart.Chart{
		Series: []dataset.Series{
			mainSeries,
			linRegSeries,
		},
	}
	chart.SetHeight(250)

	// Create unipdf chart component.
	c := creator.New()
	chartComponent := creator.NewChart(chart)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write pdf: %v", err)
	}
}
