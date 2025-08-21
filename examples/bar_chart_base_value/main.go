package main

import (
	"image/color"
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/examples"
	"github.com/unidoc/unichart/render"
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
	// Create chart component.
	profitStyle := render.Style{
		FillColor:   color.RGBA{R: 19, G: 193, B: 88, A: 255},
		StrokeColor: color.RGBA{R: 19, G: 193, B: 88, A: 255},
		StrokeWidth: 0,
	}

	lossStyle := render.Style{
		FillColor:   color.RGBA{R: 193, G: 19, B: 19, A: 255},
		StrokeColor: color.RGBA{R: 193, G: 19, B: 19, A: 255},
		StrokeWidth: 0,
	}

	chart := &unichart.BarChart{
		Title: "Bar Chart using BaseValue",
		Background: render.Style{
			Padding: render.Box{
				Top: 40,
			},
		},
		BarWidth: 30,
		YAxis: unichart.YAxis{
			Ticks: []unichart.Tick{
				{Value: -4.0, Label: "-4"},
				{Value: -2.0, Label: "-2"},
				{Value: 0, Label: "0"},
				{Value: 2.0, Label: "2"},
				{Value: 4.0, Label: "4"},
				{Value: 6.0, Label: "6"},
				{Value: 8.0, Label: "8"},
				{Value: 10.0, Label: "10"},
				{Value: 12.0, Label: "12"},
			},
		},
		UseBaseValue: true,
		BaseValue:    0.0,
		Bars: []dataset.Value{
			{Value: 10.0, Style: profitStyle, Label: "Profit"},
			{Value: 12.0, Style: profitStyle, Label: "More Profit"},
			{Value: 8.0, Style: profitStyle, Label: "Still Profit"},
			{Value: -4.0, Style: lossStyle, Label: "Loss!"},
			{Value: 3.0, Style: profitStyle, Label: "Phew Ok"},
			{Value: -2.0, Style: lossStyle, Label: "Oh No!"},
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
