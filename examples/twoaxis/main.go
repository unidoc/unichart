package main

import (
	"fmt"
	"log"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	var timestamps []float64
	for i := 0; i < 5; i++ {
		date := time.Date(2000+i, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		timestamps = append(timestamps, float64(date.UnixNano()))
	}

	// Create chart component.
	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			TickPosition: unichart.TickPositionUnderTick,
			ValueFormatter: func(v interface{}) string {
				d := time.Unix(0, int64(v.(float64)))
				fmt.Println(int64(v.(float64)), d)
				return fmt.Sprintf("%02d-%02d\n%d", d.Month(), d.Day(), d.Year())
			},
		},
		Series: []dataset.Series{
			dataset.ContinuousSeries{
				XValues: timestamps,
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			dataset.ContinuousSeries{
				YAxis:   dataset.YAxisSecondary,
				XValues: timestamps,
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
		log.Fatalf("failed to write output file: %v", err)
	}
}
