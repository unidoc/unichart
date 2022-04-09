package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	xvalues, yvalues, err := readData("requests.csv")
	if err != nil {
		log.Fatalf("failed to read input data: %v", err)
	}

	mainSeries := dataset.TimeSeries{
		Name: "Prod Request Timings",
		Style: render.Style{
			StrokeColor: render.ColorBlue,
			FillColor:   render.ColorAlternateGreen,
		},
		XValues: xvalues,
		YValues: yvalues,
	}

	linreg := &dataset.LinearRegressionSeries{
		Name: "Linear Regression",
		Style: render.Style{
			StrokeColor:     render.ColorAlternateBlue,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	sma := &dataset.SMASeries{
		Name: "SMA",
		Style: render.Style{
			StrokeColor:     render.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	chart := &unichart.Chart{
		Background: render.Style{
			Padding: render.Box{
				Top: 50,
			},
		},
		YAxis: unichart.YAxis{
			Name: "Elapsed Millis",
			TickStyle: render.Style{
				TextRotationDegrees: 45.0,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d ms", int(v.(float64)))
			},
		},
		XAxis: unichart.XAxis{
			ValueFormatter: dataset.TimeHourValueFormatter,
			GridMajorStyle: render.Style{
				StrokeColor: render.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			GridLines: releases(),
		},
		Series: []dataset.Series{
			mainSeries,
			linreg,
			dataset.LastValueAnnotationSeries(linreg),
			sma,
			dataset.LastValueAnnotationSeries(sma),
		},
	}

	chart.Elements = []render.Renderable{
		unichart.LegendThin(chart),
	}
	chart.SetHeight(450)

	// Create unipdf chart component.
	c := creator.New()
	c.SetPageSize(creator.PageSize{420 * creator.PPMM, 297 * creator.PPMM})
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

func readData(inputPath string) ([]time.Time, []float64, error) {
	var xvalues []time.Time
	var yvalues []float64

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		year := parseInt(record[0])
		month := parseInt(record[1])
		day := parseInt(record[2])
		hour := parseInt(record[3])
		elapsedMillis := parseFloat64(record[4])

		xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC))
		yvalues = append(yvalues, elapsedMillis)
	}

	return xvalues, yvalues, nil
}

func releases() []unichart.GridLine {
	return []unichart.GridLine{
		{Value: timeToFloat64(time.Date(2016, 8, 1, 9, 30, 0, 0, time.UTC))},
		{Value: timeToFloat64(time.Date(2016, 8, 2, 9, 30, 0, 0, time.UTC))},
		{Value: timeToFloat64(time.Date(2016, 8, 2, 15, 30, 0, 0, time.UTC))},
		{Value: timeToFloat64(time.Date(2016, 8, 4, 9, 30, 0, 0, time.UTC))},
		{Value: timeToFloat64(time.Date(2016, 8, 5, 9, 30, 0, 0, time.UTC))},
		{Value: timeToFloat64(time.Date(2016, 8, 6, 9, 30, 0, 0, time.UTC))},
	}
}

func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func timeToFloat64(t time.Time) float64 {
	return float64(t.UnixNano())
}
