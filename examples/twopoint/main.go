package main

import (
	"log"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	var b float64
	b = 1000

	ts1 := dataset.ContinuousSeries{
		Name:    "Time Series",
		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{1.0, 2.0, 30.0, 4.0, 50.0, 6.0, 7.0, 88.0},
	}

	ts2 := dataset.ContinuousSeries{
		Style: render.Style{
			StrokeColor: render.ColorRed,
		},

		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{15.0, 52.0, 30.0, 42.0, 50.0, 26.0, 77.0, 38.0},
	}

	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			Name:           "The XAxis",
			ValueFormatter: dataset.TimeMinuteValueFormatter,
		},
		YAxis: unichart.YAxis{
			Name: "The YAxis",
		},
		Series: []dataset.Series{
			ts1,
			ts2,
		},
	}
	chart.SetHeight(400)

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
		log.Fatalf("failed to write pdf: %v", err)
	}
}
