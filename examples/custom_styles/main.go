package main

import (
	"image/color"
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	chart := &unichart.Chart{
		Background: render.Style{
			FillColor: render.ColorBlue,
		},
		Canvas: render.Style{
			FillColor: color.RGBA{R: 239, G: 239, B: 239, A: 255},
		},
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				Style: render.Style{
					StrokeColor: render.ColorRed,
					FillColor:   render.ColorRed,
				},
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
		log.Fatalf("failed to write pdf: %v", err)
	}
}