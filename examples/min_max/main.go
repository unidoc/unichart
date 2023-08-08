package main

import (
	"log"
	"os"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
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
	mainSeries := dataset.ContinuousSeries{
		Name:    "A test series",
		XValues: sequence.Wrapper{Sequence: sequence.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		YValues: sequence.Wrapper{Sequence: sequence.NewRandomSequence().WithLen(100).WithMin(50).WithMax(150)}.Values(),
	}

	minSeries := &dataset.MinSeries{
		Style: render.Style{
			StrokeColor:     render.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	maxSeries := &dataset.MaxSeries{
		Style: render.Style{
			StrokeColor:     render.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	chart := &unichart.Chart{
		YAxis: unichart.YAxis{
			Name: "Random Values",
			Range: &sequence.ContinuousRange{
				Min: 25,
				Max: 175,
			},
		},
		XAxis: unichart.XAxis{
			Name: "Random Other Values",
		},
		Series: []dataset.Series{
			mainSeries,
			minSeries,
			maxSeries,
			dataset.LastValueAnnotationSeries(minSeries),
			dataset.LastValueAnnotationSeries(maxSeries),
		},
	}

	chart.Elements = []render.Renderable{
		unichart.Legend(chart),
	}
	chart.SetHeight(500)

	// Create unipdf chart component.
	c := creator.New()
	c.SetPageSize(creator.PageSizeA3)
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
