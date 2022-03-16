package series

import (
	"fmt"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// Interface Assertions.
var (
	_ Series                   = (*LinearSeries)(nil)
	_ data.FirstValuesProvider = (*LinearSeries)(nil)
	_ data.LastValuesProvider  = (*LinearSeries)(nil)
)

// LinearSeries is a series that plots a line in a given domain.
type LinearSeries struct {
	Name  string
	Style render.Style
	YAxis YAxisType

	XValues     []float64
	InnerSeries data.LinearCoefficientProvider

	m     float64
	b     float64
	stdev float64
	avg   float64
}

// GetName returns the name of the time series.
func (ls LinearSeries) GetName() string {
	return ls.Name
}

// GetStyle returns the line style.
func (ls LinearSeries) GetStyle() render.Style {
	return ls.Style
}

// GetYAxis returns which YAxis the series draws on.
func (ls LinearSeries) GetYAxis() YAxisType {
	return ls.YAxis
}

// Len returns the number of elements in the series.
func (ls LinearSeries) Len() int {
	return len(ls.XValues)
}

// GetEndIndex returns the effective limit end.
func (ls LinearSeries) GetEndIndex() int {
	return len(ls.XValues) - 1
}

// GetValues gets a value at a given index.
func (ls *LinearSeries) GetValues(index int) (x, y float64) {
	if ls.InnerSeries == nil || len(ls.XValues) == 0 {
		return
	}
	if ls.IsZero() {
		ls.computeCoefficients()
	}
	x = ls.XValues[index]
	y = (ls.m * ls.normalize(x)) + ls.b
	return
}

// GetFirstValues computes the first linear regression value.
func (ls *LinearSeries) GetFirstValues() (x, y float64) {
	if ls.InnerSeries == nil || len(ls.XValues) == 0 {
		return
	}
	if ls.IsZero() {
		ls.computeCoefficients()
	}
	x, y = ls.GetValues(0)
	return
}

// GetLastValues computes the last linear regression value.
func (ls *LinearSeries) GetLastValues() (x, y float64) {
	if ls.InnerSeries == nil || len(ls.XValues) == 0 {
		return
	}
	if ls.IsZero() {
		ls.computeCoefficients()
	}
	x, y = ls.GetValues(ls.GetEndIndex())
	return
}

// Render renders the series.
func (ls *LinearSeries) Render(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, defaults render.Style) {
	drawLineSeries(r, canvasBox, xrange, yrange, ls.Style.InheritFrom(defaults), ls)
}

// Validate validates the series.
func (ls LinearSeries) Validate() error {
	if ls.InnerSeries == nil {
		return fmt.Errorf("linear regression series requires InnerSeries to be set")
	}
	return nil
}

// IsZero returns if the linear series has computed coefficients or not.
func (ls LinearSeries) IsZero() bool {
	return ls.m == 0 && ls.b == 0
}

// computeCoefficients computes the `m` and `b` terms in the linear formula given by `y = mx+b`.
func (ls *LinearSeries) computeCoefficients() {
	ls.m, ls.b, ls.stdev, ls.avg = ls.InnerSeries.Coefficients()
}

func (ls *LinearSeries) normalize(xvalue float64) float64 {
	if ls.avg > 0 && ls.stdev > 0 {
		return (xvalue - ls.avg) / ls.stdev
	}
	return xvalue
}

// drawLineSeries draws a line series with a renderer.
func drawLineSeries(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, style render.Style, vs data.ValuesProvider) {
	if vs.Len() == 0 {
		return
	}

	cb := canvasBox.Bottom
	cl := canvasBox.Left

	v0x, v0y := vs.GetValues(0)
	x0 := cl + xrange.Translate(v0x)
	y0 := cb - yrange.Translate(v0y)

	yv0 := yrange.Translate(0)

	var vx, vy float64
	var x, y int

	if style.ShouldDrawStroke() && style.ShouldDrawFill() {
		style.GetFillOptions().WriteDrawingOptionsToRenderer(r)
		r.MoveTo(x0, y0)
		for i := 1; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)
			r.LineTo(x, y)
		}
		r.LineTo(x, mathutil.MinInt(cb, cb-yv0))
		r.LineTo(x0, mathutil.MinInt(cb, cb-yv0))
		r.LineTo(x0, y0)
		r.Fill()
	}

	if style.ShouldDrawStroke() {
		style.GetStrokeOptions().WriteDrawingOptionsToRenderer(r)

		r.MoveTo(x0, y0)
		for i := 1; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)
			r.LineTo(x, y)
		}
		r.Stroke()
	}

	if style.ShouldDrawDot() {
		defaultDotWidth := style.GetDotWidth()

		style.GetDotOptions().WriteDrawingOptionsToRenderer(r)
		for i := 0; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)

			dotWidth := defaultDotWidth
			if style.DotWidthProvider != nil {
				dotWidth = style.DotWidthProvider(xrange, yrange, i, vx, vy)
			}

			if style.DotColorProvider != nil {
				dotColor := style.DotColorProvider(xrange, yrange, i, vx, vy)

				r.SetFillColor(dotColor)
				r.SetStrokeColor(dotColor)
			}

			r.Circle(dotWidth, x, y)
			r.FillStroke()
		}
	}
}
