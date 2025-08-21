package main

import (
	"log"
	"os"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/examples"
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
	now := time.Date(2023, 8, 17, 0, 0, 0, 0, time.UTC)

	chart := &unichart.Chart{
		Series: []dataset.Series{
			dataset.TimeSeries{
				XValues: []time.Time{
					now.AddDate(0, 0, -10),
					now.AddDate(0, 0, -9),
					now.AddDate(0, 0, -8),
					now.AddDate(0, 0, -7),
					now.AddDate(0, 0, -6),
					now.AddDate(0, 0, -5),
					now.AddDate(0, 0, -4),
					now.AddDate(0, 0, -3),
					now.AddDate(0, 0, -2),
					now.AddDate(0, 0, -1),
					now,
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
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

	examples.RenderPDFToImage("output.pdf")
}
