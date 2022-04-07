package main

import (
	"image/color"
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	viridisByY := func(xr, yr sequence.Range, index int, x, y float64) color.Color {
		return render.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	chart := &unichart.Chart{
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				Style: render.Style{
					StrokeWidth:      -1,
					DotWidth:         5,
					DotColorProvider: viridisByY,
				},
				XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(0).WithEnd(127)}.Values(),
				YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(128).WithMin(0).WithMax(1024)}.Values(),
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
