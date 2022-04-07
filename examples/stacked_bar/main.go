package main

import (
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	chart := &unichart.StackedBarChart{
		Title: "Test Stacked Bar Chart",
		Background: render.Style{
			Padding: render.Box{
				Top: 40,
			},
		},
		Bars: []unichart.StackedBar{
			{
				Name: "This is a very long string to test word break wrapping.",
				Values: []dataset.Value{
					{Value: 5, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 4, Label: "Gray"},
					{Value: 3, Label: "Orange"},
					{Value: 3, Label: "Test"},
					{Value: 2, Label: "??"},
					{Value: 1, Label: "!!"},
				},
			},
			{
				Name: "Test",
				Values: []dataset.Value{
					{Value: 10, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 1, Label: "Gray"},
				},
			},
			{
				Name: "Test 2",
				Values: []dataset.Value{
					{Value: 10, Label: "Blue"},
					{Value: 6, Label: "Green"},
					{Value: 4, Label: "Gray"},
				},
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
