package main

import (
	"fmt"
	"log"
	"time"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/data/series"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	// Create chart component.
	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			TickPosition: unichart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typedDate := time.Unix(0, int64(v.(float64)))
				return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []series.Series{
			series.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			series.ContinuousSeries{
				YAxis:   series.YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}
	chart.SetHeight(300)

	// Create unipdf chart component.
	c := creator.New()
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
