package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/unidoc/unichart"
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/creator"
)

func main() {
	xVals, yVals, err := readData("data.csv")
	if err != nil {
		log.Fatalf("failed to read input data: %v", err)
	}

	priceSeries := dataset.TimeSeries{
		Name: "SPY",
		Style: render.Style{
			StrokeColor: render.ColorBlue,
		},
		XValues: xVals,
		YValues: yVals,
	}

	smaSeries := dataset.SMASeries{
		Name: "SPY - SMA",
		Style: render.Style{
			StrokeColor:     render.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &dataset.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: render.Style{
			StrokeColor: render.ColorLightGray,
			FillColor:   render.ColorLightGray,
		},
		InnerSeries: priceSeries,
	}

	chart := &unichart.Chart{
		XAxis: unichart.XAxis{
			TickPosition: unichart.TickPositionBetweenTicks,
		},
		YAxis: unichart.YAxis{
			Range: &sequence.ContinuousRange{
				Max: 220.0,
				Min: 180.0,
			},
		},
		Series: []dataset.Series{
			bbSeries,
			priceSeries,
			smaSeries,
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
		log.Fatalf("failed to write output file: %v", err)
	}
}

func readData(inputPath string) ([]time.Time, []float64, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	var (
		xVals []time.Time
		yVals []float64
	)

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		xVal, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, nil, err
		}
		xVals = append(xVals, xVal)

		yVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, nil, err
		}
		yVals = append(yVals, yVal)
	}

	return xVals, yVals, nil
}
