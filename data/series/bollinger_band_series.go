package series

import (
	"fmt"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/data/sequence"
	"github.com/unidoc/unichart/render"
)

// Interface Assertions.
var (
	_ Series = (*BollingerBandsSeries)(nil)
)

// BollingerBandsSeries draws bollinger bands for an inner series.
// Bollinger bands are defined by two lines, one at SMA+k*stddev, one at SMA-k*stdev.
type BollingerBandsSeries struct {
	Name  string
	Style render.Style
	YAxis YAxisType

	Period      int
	K           float64
	InnerSeries data.ValuesProvider

	valueBuffer *data.ValueBuffer
}

// GetName returns the name of the time series.
func (bbs BollingerBandsSeries) GetName() string {
	return bbs.Name
}

// GetStyle returns the line style.
func (bbs BollingerBandsSeries) GetStyle() render.Style {
	return bbs.Style
}

// GetYAxis returns which YAxis the series draws on.
func (bbs BollingerBandsSeries) GetYAxis() YAxisType {
	return bbs.YAxis
}

// GetPeriod returns the window size.
func (bbs BollingerBandsSeries) GetPeriod() int {
	if bbs.Period == 0 {
		return DefaultSimpleMovingAveragePeriod
	}
	return bbs.Period
}

// GetK returns the K value, or the number of standard deviations above and below
// to band the simple moving average with.
// Typical K value is 2.0.
func (bbs BollingerBandsSeries) GetK(defaults ...float64) float64 {
	if bbs.K == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 2.0
	}
	return bbs.K
}

// Len returns the number of elements in the series.
func (bbs BollingerBandsSeries) Len() int {
	return bbs.InnerSeries.Len()
}

// GetBoundedValues gets the bounded value for the series.
func (bbs *BollingerBandsSeries) GetBoundedValues(index int) (x, y1, y2 float64) {
	if bbs.InnerSeries == nil {
		return
	}
	if bbs.valueBuffer == nil || index == 0 {
		bbs.valueBuffer = data.NewValueBufferWithCapacity(bbs.GetPeriod())
	}
	if bbs.valueBuffer.Len() >= bbs.GetPeriod() {
		bbs.valueBuffer.Dequeue()
	}
	px, py := bbs.InnerSeries.GetValues(index)
	bbs.valueBuffer.Enqueue(py)
	x = px

	ay := sequence.NewWrapper(bbs.valueBuffer).Average()
	std := sequence.NewWrapper(bbs.valueBuffer).StdDev()

	y1 = ay + (bbs.GetK() * std)
	y2 = ay - (bbs.GetK() * std)
	return
}

// GetBoundedLastValues returns the last bounded value for the series.
func (bbs *BollingerBandsSeries) GetBoundedLastValues() (x, y1, y2 float64) {
	if bbs.InnerSeries == nil {
		return
	}
	period := bbs.GetPeriod()
	seriesLength := bbs.InnerSeries.Len()
	startAt := seriesLength - period
	if startAt < 0 {
		startAt = 0
	}

	vb := data.NewValueBufferWithCapacity(period)
	for index := startAt; index < seriesLength; index++ {
		xn, yn := bbs.InnerSeries.GetValues(index)
		vb.Enqueue(yn)
		x = xn
	}

	ay := sequence.NewWrapper(vb).Average()
	std := sequence.NewWrapper(vb).StdDev()

	y1 = ay + (bbs.GetK() * std)
	y2 = ay - (bbs.GetK() * std)

	return
}

// Render renders the series.
func (bbs *BollingerBandsSeries) Render(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, defaults render.Style) {
	s := bbs.Style.InheritFrom(defaults.InheritFrom(render.Style{
		StrokeWidth: 1.0,
		StrokeColor: render.ColorWithAlpha(render.DefaultLineColor, 64),
		FillColor:   render.ColorWithAlpha(render.DefaultLineColor, 32),
	}))

	drawBoundedSeries(r, canvasBox, xrange, yrange, s, bbs, bbs.GetPeriod())
}

// Validate validates the series.
func (bbs BollingerBandsSeries) Validate() error {
	if bbs.InnerSeries == nil {
		return fmt.Errorf("bollinger bands series requires InnerSeries to be set")
	}
	return nil
}

// drawBoundedSeries draws a series that implements BoundedValuesProvider.
func drawBoundedSeries(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, style render.Style, bbs data.BoundedValuesProvider, drawOffsetIndexes ...int) {
	drawOffsetIndex := 0
	if len(drawOffsetIndexes) > 0 {
		drawOffsetIndex = drawOffsetIndexes[0]
	}

	cb := canvasBox.Bottom
	cl := canvasBox.Left

	v0x, v0y1, v0y2 := bbs.GetBoundedValues(0)
	x0 := cl + xrange.Translate(v0x)
	y0 := cb - yrange.Translate(v0y1)

	var vx, vy1, vy2 float64
	var x, y int

	xvalues := make([]float64, bbs.Len())
	xvalues[0] = v0x
	y2values := make([]float64, bbs.Len())
	y2values[0] = v0y2

	style.GetFillAndStrokeOptions().WriteToRenderer(r)
	r.MoveTo(x0, y0)
	for i := 1; i < bbs.Len(); i++ {
		vx, vy1, vy2 = bbs.GetBoundedValues(i)

		xvalues[i] = vx
		y2values[i] = vy2

		x = cl + xrange.Translate(vx)
		y = cb - yrange.Translate(vy1)
		if i > drawOffsetIndex {
			r.LineTo(x, y)
		} else {
			r.MoveTo(x, y)
		}
	}
	y = cb - yrange.Translate(vy2)
	r.LineTo(x, y)
	for i := bbs.Len() - 1; i >= drawOffsetIndex; i-- {
		vx, vy2 = xvalues[i], y2values[i]
		x = cl + xrange.Translate(vx)
		y = cb - yrange.Translate(vy2)
		r.LineTo(x, y)
	}
	r.Close()
	r.FillStroke()
}
