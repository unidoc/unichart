package main

import (
	"log"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	// Create chart component.
	chart := &unichart.Chart{
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			dataset.AnnotationSeries{
				Annotations: []dataset.Value2{
					{XValue: 1.0, YValue: 1.0, Label: "One"},
					{XValue: 2.0, YValue: 2.0, Label: "Two"},
					{XValue: 3.0, YValue: 3.0, Label: "Three"},
					{XValue: 4.0, YValue: 4.0, Label: "Four"},
					{XValue: 5.0, YValue: 5.0, Label: "Five"},
				},
			},
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
}
