package main

import (
	"image/color"
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
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
	var (
		colorWhite          = color.RGBA{R: 241, G: 241, B: 241, A: 255}
		colorMariner        = color.RGBA{R: 60, G: 100, B: 148, A: 255}
		colorLightSteelBlue = color.RGBA{R: 182, G: 195, B: 220, A: 255}
		colorPoloBlue       = color.RGBA{R: 126, G: 155, B: 200, A: 255}
		colorSteelBlue      = color.RGBA{R: 73, G: 120, B: 177, A: 255}

		barWidth = 80
	)

	// Create chart component.
	chart := &unichart.StackedBarChart{
		Title: "Quarterly Sales",
		Background: render.Style{
			Padding: render.Box{Top: 50},
		},
		BarSpacing: 20,
		Bars: []unichart.StackedBar{
			{
				Name:  "Q1",
				Width: barWidth,
				Values: []dataset.Value{
					{
						Label: "32K",
						Value: 32,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "46K",
						Value: 46,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "48K",
						Value: 48,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "42K",
						Value: 42,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontColor:   colorWhite,
						},
					},
				},
			},
			{
				Name:  "Q2",
				Width: barWidth,
				Values: []dataset.Value{
					{
						Label: "45K",
						Value: 45,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "60K",
						Value: 60,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "62K",
						Value: 62,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "53K",
						Value: 53,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontColor:   colorWhite,
						},
					},
				},
			},
			{
				Name:  "Q3",
				Width: barWidth,
				Values: []dataset.Value{
					{
						Label: "54K",
						Value: 54,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "58K",
						Value: 58,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "55K",
						Value: 55,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "47K",
						Value: 47,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontColor:   colorWhite,
						},
					},
				},
			},
			{
				Name:  "Q4",
				Width: barWidth,
				Values: []dataset.Value{
					{
						Label: "46K",
						Value: 46,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "70K",
						Value: 70,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "74K",
						Value: 74,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontColor:   colorWhite,
						},
					},
					{
						Label: "60K",
						Value: 60,
						Style: render.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontColor:   colorWhite,
						},
					},
				},
			},
		},
	}
	chart.SetHeight(500)

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
