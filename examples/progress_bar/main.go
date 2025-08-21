package main

import (
	"image/color"
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/examples"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/model"
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
	c := creator.New()

	linearProgressBars(c)
	circularProgressBars(c)

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	examples.RenderPDFToImage("output.pdf")
}

func linearProgressBars(c *creator.Creator) {
	labelFont, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Println(err)
	}

	// Linear progress bars
	addLinearProgressBar(c, 0.32, "", true, true, 20,
		render.Style{
			FillColor:   render.ColorAlternateLightGray,
			StrokeWidth: 1.0,
			StrokeColor: render.ColorRed,
		},
		render.Style{
			FillColor: render.ColorAlternateGreen,
		},
		render.Style{},
	)

	addLinearProgressBar(c, 0.68, "", false, true, 20,
		render.Style{
			FillColor: render.ColorBlue,
		},
		render.Style{
			FillColor: render.ColorRed,
		},
		render.Style{},
	)

	addLinearProgressBar(c, 0.96, "", true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 221, G: 254, B: 218, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 144, G: 205, B: 62, A: 255},
		},
		render.Style{},
	)

	addLinearProgressBar(c, 0.7432, "74.32%", true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 248, G: 239, B: 223, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 209, G: 145, B: 85, A: 255},
		},
		render.Style{
			FontSize: 12,
			Font:     labelFont,
		},
	)

	addLinearProgressBar(c, 0.32, "", true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 254, G: 246, B: 245, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
		},
		render.Style{},
	)

	addLinearProgressBarWithCustomInfo(c, 0.32, "32%", true, true, 16,
		render.Style{
			FillColor:   color.RGBA{R: 254, G: 246, B: 245, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 196, G: 196, B: 196, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
		},
		render.Style{
			FontSize: 12,
			Font:     labelFont,
		},
		func(r render.Renderer, x int) int {
			r.MoveTo(x, -14)
			r.LineTo(x, 15)
			r.SetStrokeWidth(1.0)
			r.SetStrokeColor(color.RGBA{R: 214, G: 214, B: 214, A: 255})
			r.Stroke()

			render.Text.Draw(r, "Custom top label here", x+10, 0,
				render.Style{
					Font:      model.DefaultFont(),
					FontSize:  14,
					FontColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
				},
			)

			return 20
		},
		func(r render.Renderer, x int) int {
			r.MoveTo(x, 40)
			r.LineTo(x, 70)
			r.SetStrokeWidth(1.0)
			r.SetStrokeColor(color.RGBA{R: 214, G: 214, B: 214, A: 255})
			r.Stroke()

			render.Text.Draw(r, "Custom bottom label here", x+10, 70,
				render.Style{
					Font:      model.DefaultFont(),
					FontSize:  14,
					FontColor: color.RGBA{R: 143, G: 47, B: 26, A: 255},
				},
			)

			return 30
		},
	)

	addLinearProgressBar(c, 0.32, "", false, false, 25,
		render.Style{
			FillColor:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			StrokeWidth: 1.0,
			StrokeColor: color.RGBA{R: 218, G: 218, B: 218, A: 255},
		},
		render.Style{
			FillColor: color.RGBA{R: 242, G: 172, B: 59, A: 255},
		},
		render.Style{},
	)
}

func addLinearProgressBar(c *creator.Creator, progress float64, label string, roundStart bool, roundEnd bool,
	height int, bgStyle render.Style, fgStyle render.Style, labelStyle render.Style) {
	addLinearProgressBarWithCustomInfo(c, progress, label, roundStart, roundEnd, height, bgStyle, fgStyle, labelStyle, nil, nil)
}

func addLinearProgressBarWithCustomInfo(c *creator.Creator, progress float64, label string, roundStart bool, roundEnd bool,
	height int, bgStyle render.Style, fgStyle render.Style, labelStyle render.Style,
	customTopInfo func(r render.Renderer, x int) int, customBottomInfo func(r render.Renderer, x int) int) {

	// Create chart component.
	linear := &unichart.LinearProgressBar{
		BackgroundStyle: bgStyle,
		ForegroundStyle: fgStyle,
		LabelStyle:      labelStyle,

		RoundedEdgeStart: roundStart,
		RoundedEdgeEnd:   roundEnd,

		CustomTopInfo:    customTopInfo,
		CustomBottomInfo: customBottomInfo,
	}
	linear.SetHeight(height)
	linear.SetProgress(progress)
	linear.SetLabel(label)

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
