package main

import (
	"image/color"
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	c := creator.New()

	linearProgressBars(c)
	circularProgressBars(c)

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}

func linearProgressBars(c *creator.Creator) {
	// Linear progress bars
	addLinearProgressBar(c, 0.32, true, true, 20,
		render.Style{
			FillColor:   render.ColorAlternateLightGray,
			StrokeWidth: 1.0,
			StrokeColor: render.ColorRed,
		},
		render.Style{
			FillColor: render.ColorAlternateGreen,
		},
	)

	addLinearProgressBar(c, 0.68, false, true, 20,
		render.Style{
			FillColor: render.ColorBlue,
		},
		render.Style{
			FillColor: render.ColorRed,
		},
	)

	addLinearProgressBar(c, 0.96, true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 221, G: 254, B: 218, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 144, G: 205, B: 62, A: 255},
		},
	)

	addLinearProgressBar(c, 0.7432, true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 248, G: 239, B: 223, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 209, G: 145, B: 85, A: 255},
		},
	)

	addLinearProgressBar(c, 0.32, true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 254, G: 246, B: 245, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
		},
	)
}

func addLinearProgressBar(c *creator.Creator, progress float64, roundStart bool, roundEnd bool,
	height int, bgStyle render.Style, fgStyle render.Style) {
	// Create chart component.
	linear := &unichart.LinearProgressBar{
		BackgroundStyle: bgStyle,
		ForegroundStyle: fgStyle,

		RoundedEdgeStart: roundStart,
		RoundedEdgeEnd:   roundEnd,
	}
	linear.SetHeight(height)
	linear.SetProgress(progress)

	// Create unipdf chart component.
	chartComponent := creator.NewChart(linear)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}
}

func circularProgressBars(c *creator.Creator) {
	// Circular progress bars
	labelFont, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Println(err)
	}

	addCircularProgressBar(c, 0.68, 50, "", false,
		render.Style{
			StrokeWidth: 10.0,
			StrokeColor: render.ColorAlternateLightGray,
		},
		render.Style{
			StrokeWidth: 10.0,
			StrokeColor: render.ColorAlternateGreen,
		},
		render.Style{
			FontSize: 20,
			Font:     labelFont,
		},
	)

	addCircularProgressBar(c, 0.96, 80, "A", true,
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 221, G: 254, B: 218, A: 255},
		},
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 144, G: 205, B: 62, A: 255},
		},
		render.Style{
			FontSize: 30,
			Font:     labelFont,
		},
	)

	addCircularProgressBar(c, 0.6934, 80, "C", true,
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 248, G: 239, B: 223, A: 255},
		},
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 209, G: 145, B: 85, A: 255},
		},
		render.Style{
			FontSize: 30,
			Font:     labelFont,
		},
	)

	addCircularProgressBar(c, 0.4906, 80, "F", true,
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 254, G: 246, B: 245, A: 255},
		},
		render.Style{
			StrokeWidth: 15.0,
			StrokeColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
		},
		render.Style{
			FontSize: 30,
			Font:     labelFont,
		},
	)
}

func addCircularProgressBar(c *creator.Creator, progress float64, size int, label string, reversed bool,
	bgStyle render.Style, fgStyle render.Style, labelStyle render.Style) {
	circular := &unichart.CircularProgressBar{
		BackgroundStyle: bgStyle,
		ForegroundStyle: fgStyle,
		LabelStyle:      labelStyle,

		Reversed: reversed,
	}
	circular.SetSize(size)
	circular.SetProgress(progress)
	circular.SetLabel(label)

	chartComponent := creator.NewChart(circular)
	chartComponent.SetMargins(0, 0, 20, 20)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}
}
