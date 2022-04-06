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
	// Create chart component.
	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			Name: "The XAxis",
		},
		YAxis: unichart.YAxis{
			Name: "The YAxis",
		},
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				Style: render.Style{
					StrokeColor: color.RGBA{R: 0, G: 116, B: 217, A: 255},
					FillColor:   color.RGBA{R: 106, G: 195, B: 203, A: 255},
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
