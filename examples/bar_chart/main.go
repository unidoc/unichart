package main

import (
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
	// Create chart component.
	chart := &unichart.BarChart{
		Title: "Test Bar Chart",
		Background: render.Style{
			Padding: render.Box{
				Top: 40,
			},
		},
		BarWidth: 40,
		Bars: []dataset.Value{
			{Value: 5.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}
	chart.SetHeight(400)

	// Set Y-axis custom range.
	var max float64
	for _, bar := range chart.Bars {
		if max < bar.Value {
			max = bar.Value
		}
	}

	rng := &sequence.ContinuousRange{}
	rng.SetMin(0)
	rng.SetMax(max)
	chart.YAxis.Range = rng

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
