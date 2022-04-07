package main

import (
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	chart := &unichart.Chart{
		YAxis: unichart.YAxis{
			Range: &sequence.ContinuousRange{
				Min: 0.0,
				Max: 4.0,
			},
			Ticks: []unichart.Tick{
				{Value: 0.0, Label: "0.00"},
				{Value: 2.0, Label: "2.00"},
				{Value: 4.0, Label: "4.00"},
				{Value: 6.0, Label: "6.00"},
				{Value: 8.0, Label: "Eight"},
				{Value: 10.0, Label: "Ten"},
			},
		},
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
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
