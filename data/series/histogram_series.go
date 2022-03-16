package series

import (
	"fmt"
	"math"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/render"
)

// HistogramSeries is a special type of series that draws as a histogram.
// Some peculiarities; it will always be lower bounded at 0 (at the very least).
// This may alter ranges a bit and generally you want to put a histogram series on it's own y-axis.
type HistogramSeries struct {
	Name        string
	Style       render.Style
	YAxis       YAxisType
	InnerSeries data.ValuesProvider
}

// GetName implements Series.GetName.
func (hs HistogramSeries) GetName() string {
	return hs.Name
}

// GetStyle implements Series.GetStyle.
func (hs HistogramSeries) GetStyle() render.Style {
	return hs.Style
}

// GetYAxis returns which yaxis the series is mapped to.
func (hs HistogramSeries) GetYAxis() YAxisType {
	return hs.YAxis
}

// Len implements BoundedValuesProvider.Len.
func (hs HistogramSeries) Len() int {
	return hs.InnerSeries.Len()
}

// GetValues implements ValuesProvider.GetValues.
func (hs HistogramSeries) GetValues(index int) (x, y float64) {
	return hs.InnerSeries.GetValues(index)
}

// GetBoundedValues implements BoundedValuesProvider.GetBoundedValue
func (hs HistogramSeries) GetBoundedValues(index int) (x, y1, y2 float64) {
	vx, vy := hs.InnerSeries.GetValues(index)

	x = vx

	if vy > 0 {
		y1 = vy
		return
	}

	y2 = vy
	return
}

// Render implements Series.Render.
func (hs HistogramSeries) Render(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, defaults render.Style) {
	style := hs.Style.InheritFrom(defaults)
	drawHistogramSeries(r, canvasBox, xrange, yrange, style, hs)
}

// Validate validates the series.
func (hs HistogramSeries) Validate() error {
	if hs.InnerSeries == nil {
		return fmt.Errorf("histogram series requires InnerSeries to be set")
	}
	return nil
}

// drawHistogramSeries draws a value provider as boxes from 0.
func drawHistogramSeries(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, style render.Style, vs data.ValuesProvider, barWidths ...int) {
	if vs.Len() == 0 {
		return
	}

	//calculate bar width?
	seriesLength := vs.Len()
	barWidth := int(math.Floor(float64(xrange.GetDomain()) / float64(seriesLength)))
	if len(barWidths) > 0 {
		barWidth = barWidths[0]
	}

	cb := canvasBox.Bottom
	cl := canvasBox.Left

	//foreach datapoint, draw a box.
	for index := 0; index < seriesLength; index++ {
		vx, vy := vs.GetValues(index)
		y0 := yrange.Translate(0)
		x := cl + xrange.Translate(vx)
		y := yrange.Translate(vy)

		render.Box{
			Top:    cb - y0,
			Left:   x - (barWidth >> 1),
			Right:  x + (barWidth >> 1),
			Bottom: cb - y,
		}.Draw(r, style)
	}
}
