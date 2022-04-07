package main

import (
	"log"
	"time"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	now := time.Now()

	chart := &unichart.Chart{
		Series: []dataset.Series{
			dataset.TimeSeries{
				XValues: []time.Time{
					now.AddDate(0, 0, -10),
					now.AddDate(0, 0, -9),
					now.AddDate(0, 0, -8),
					now.AddDate(0, 0, -7),
					now.AddDate(0, 0, -6),
					now.AddDate(0, 0, -5),
					now.AddDate(0, 0, -4),
					now.AddDate(0, 0, -3),
					now.AddDate(0, 0, -2),
					now.AddDate(0, 0, -1),
					now,
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
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
}
