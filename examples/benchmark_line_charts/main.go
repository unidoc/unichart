package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	unichart "github.com/unidoc/unichart"
	"github.com/unidoc/unichart/data/series"
	"github.com/unidoc/unipdf/v3/creator"
)

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func main() {
	// Create chart component.
	numValues := 1024
	numSeries := 100
	timeSeries := make([]series.Series, numSeries)

	for i := 0; i < numSeries; i++ {
		xValues := make([]time.Time, numValues)
		yValues := make([]float64, numValues)

		for j := 0; j < numValues; j++ {
			xValues[j] = time.Now().AddDate(0, 0, (numValues-j)*-1)
			yValues[j] = random(float64(-500), float64(500))
		}

		timeSeries[i] = series.TimeSeries{
			Name:    fmt.Sprintf("aaa.bbb.hostname-%v.ccc.ddd.eee.fff.ggg.hhh.iii.jjj.kkk.lll.mmm.nnn.value", i),
			XValues: xValues,
			YValues: yValues,
		}
	}

	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			Name: "Time",
		},
		YAxis: unichart.YAxis{
			Name: "Value",
		},
		Series: timeSeries,
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
		log.Fatalf("failed to write pdf: %v", err)
	}
}
