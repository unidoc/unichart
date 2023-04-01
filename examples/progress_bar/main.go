package main

import (
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
	// Create chart component.
	linear := &unichart.LinearProgressChart{
		BackgroundStyle: render.Style{
			FillColor:   render.ColorAlternateLightGray,
			StrokeWidth: 1.0,
			StrokeColor: render.ColorRed,
		},
		ForegroundStyle: render.Style{
			FillColor: render.ColorAlternateGreen,
		},

		RoundedEdgeStart: true,
		RoundedEdgeEnd:   true,
	}
	linear.SetHeight(20)
	linear.SetProgress(0.32)

	// Create unipdf chart component.
	c := creator.New()
	chartComponent := creator.NewChart(linear)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	// Create chart component.
	linear = &unichart.LinearProgressChart{
		BackgroundStyle: render.Style{
			FillColor: render.ColorBlue,
		},
		ForegroundStyle: render.Style{
			FillColor: render.ColorRed,
		},

		RoundedEdgeStart: false,
		RoundedEdgeEnd:   true,
	}
	linear.SetHeight(20)
	linear.SetProgress(0.68)

	// Create unipdf chart component.
	chartComponent = creator.NewChart(linear)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	labelFont, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		log.Println(err)
	}

	circular := &unichart.CircularProgressChart{
		BackgroundStyle: render.Style{
			StrokeWidth: 10.0,
			StrokeColor: render.ColorAlternateLightGray,
		},
		ForegroundStyle: render.Style{
			StrokeWidth: 10.0,
			StrokeColor: render.ColorAlternateGreen,
		},
		LabelStyle: render.Style{
			FontSize: 20,
			Font:     labelFont,
		},

		Reversed: true,
	}
	circular.SetSize(50)
	circular.SetProgress(0.68)

	chartComponent = creator.NewChart(circular)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	circular = &unichart.CircularProgressChart{
		BackgroundStyle: render.Style{
			StrokeWidth: 20.0,
			StrokeColor: render.ColorYellow,
		},
		ForegroundStyle: render.Style{
			StrokeWidth: 20.0,
			StrokeColor: render.ColorAlternateGreen,
		},
		LabelStyle: render.Style{
			FontSize: 30,
			Font:     labelFont,
		},
	}
	circular.SetSize(100)
	circular.SetProgress(0.3)
	circular.SetLabel("30%")

	chartComponent = creator.NewChart(circular)

	// Draw chart component.
	if err := c.Draw(chartComponent); err != nil {
		log.Fatalf("failed to draw chart: %v", err)
	}

	// Save output file.
	if err := c.WriteToFile("output.pdf"); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}
