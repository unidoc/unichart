package unichart

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// Chart represents a line, curve or histogram chart.
type Chart struct {
	Title      string
	TitleStyle render.Style

	Font         render.Font
	Background   render.Style
	Canvas       render.Style
	ColorPalette render.ColorPalette

	XAxis          XAxis
	YAxis          YAxis
	YAxisSecondary YAxis

	Series   []dataset.Series
	Elements []render.Renderable

	width  int
	height int
	dpi    float64
}

// DPI returns the DPI for the chart.
func (c *Chart) DPI(defaults ...float64) float64 {
	if c.dpi == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return defaultDPI
	}
	return c.dpi
}

// SetDPI sets the DPI for the chart.
func (c *Chart) SetDPI(dpi float64) {
	c.dpi = dpi
}

// GetFont returns the text font.
func (c *Chart) GetFont() render.Font {
	return c.Font
}

// Width returns the chart width.
func (c *Chart) Width() int {
	if c.width == 0 {
		return defaultChartWidth
	}
	return c.width
}

// SetWidth sets the chart width.
func (c *Chart) SetWidth(width int) {
	c.width = width
}

// Height returns the chart height.
func (c *Chart) Height() int {
	if c.height == 0 {
		return defaultChartHeight
	}
	return c.height
}

// SetHeight sets the chart height.
func (c *Chart) SetHeight(height int) {
	c.height = height
}

// Render renders the chart with the given renderer to the given io.Writer.
func (c *Chart) Render(rp render.RendererProvider, w io.Writer) error {
	if len(c.Series) == 0 {
		return errors.New("please provide at least one series")
	}
	if err := c.checkHasVisibleSeries(); err != nil {
		return err
	}

	c.YAxisSecondary.AxisType = dataset.YAxisSecondary

	r, err := rp(c.Width(), c.Height())
	if err != nil {
		return err
	}
	r.SetDPI(c.DPI(defaultDPI))

	c.drawBackground(r)

	var xt, yt, yta []Tick
	xr, yr, yra := c.getRanges()
	canvasBox := c.getDefaultCanvasBox()
	xf, yf, yfa := c.getValueFormatters()
	xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)

	err = c.checkRanges(xr, yr, yra)
	if err != nil {
		r.Save(w)
		return err
	}

	if c.hasAxes() {
		xt, yt, yta = c.getAxesTicks(r, xr, yr, yra, xf, yf, yfa)
		canvasBox = c.getAxesAdjustedCanvasBox(r, canvasBox, xr, yr, yra, xt, yt, yta)
		xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)

		// do a second pass in case things haven't settled yet.
		xt, yt, yta = c.getAxesTicks(r, xr, yr, yra, xf, yf, yfa)
		canvasBox = c.getAxesAdjustedCanvasBox(r, canvasBox, xr, yr, yra, xt, yt, yta)
		xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)
	}

	if c.hasAnnotationSeries() {
		canvasBox = c.getAnnotationAdjustedCanvasBox(r, canvasBox, xr, yr, yra, xf, yf, yfa)
		xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)
		xt, yt, yta = c.getAxesTicks(r, xr, yr, yra, xf, yf, yfa)
	}

	c.drawCanvas(r, canvasBox)
	c.drawAxes(r, canvasBox, xr, yr, yra, xt, yt, yta)
	for index, series := range c.Series {
		c.drawSeries(r, canvasBox, xr, yr, yra, series, index)
	}

	c.drawTitle(r)

	for _, a := range c.Elements {
		a(r, canvasBox, c.styleDefaultsElements())
	}

	return r.Save(w)
}

func (c *Chart) checkHasVisibleSeries() error {
	var style render.Style
	for _, s := range c.Series {
		style = s.GetStyle()
		if !style.Hidden {
			return nil
		}
	}
	return fmt.Errorf("chart render; must have (1) visible series")
}

func (c *Chart) validateSeries() error {
	var err error
	for _, s := range c.Series {
		err = s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Chart) getRanges() (xrange, yrange, yrangeAlt sequence.Range) {
	var minx, maxx float64 = math.MaxFloat64, -math.MaxFloat64
	var miny, maxy float64 = math.MaxFloat64, -math.MaxFloat64
	var minya, maxya float64 = math.MaxFloat64, -math.MaxFloat64

	seriesMappedToSecondaryAxis := false

	// Note: a possible future optimization is to not scan the series values
	// if all axis are represented by either custom ticks or custom ranges.
	for _, s := range c.Series {
		if !s.GetStyle().Hidden {
			seriesAxis := s.GetYAxis()
			if bvp, isBoundedValuesProvider := s.(dataset.BoundedValuesProvider); isBoundedValuesProvider {
				seriesLength := bvp.Len()
				for index := 0; index < seriesLength; index++ {
					vx, vy1, vy2 := bvp.GetBoundedValues(index)

					minx = math.Min(minx, vx)
					maxx = math.Max(maxx, vx)

					if seriesAxis == dataset.YAxisPrimary {
						miny = math.Min(miny, vy1)
						miny = math.Min(miny, vy2)
						maxy = math.Max(maxy, vy1)
						maxy = math.Max(maxy, vy2)
					} else if seriesAxis == dataset.YAxisSecondary {
						minya = math.Min(minya, vy1)
						minya = math.Min(minya, vy2)
						maxya = math.Max(maxya, vy1)
						maxya = math.Max(maxya, vy2)
						seriesMappedToSecondaryAxis = true
					}
				}
			} else if vp, isValuesProvider := s.(dataset.ValuesProvider); isValuesProvider {
				seriesLength := vp.Len()
				for index := 0; index < seriesLength; index++ {
					vx, vy := vp.GetValues(index)

					minx = math.Min(minx, vx)
					maxx = math.Max(maxx, vx)

					if seriesAxis == dataset.YAxisPrimary {
						miny = math.Min(miny, vy)
						maxy = math.Max(maxy, vy)
					} else if seriesAxis == dataset.YAxisSecondary {
						minya = math.Min(minya, vy)
						maxya = math.Max(maxya, vy)
						seriesMappedToSecondaryAxis = true
					}
				}
			}
		}
	}

	if c.XAxis.Range == nil {
		xrange = &sequence.ContinuousRange{}
	} else {
		xrange = c.XAxis.Range
	}

	if c.YAxis.Range == nil {
		yrange = &sequence.ContinuousRange{}
	} else {
		yrange = c.YAxis.Range
	}

	if c.YAxisSecondary.Range == nil {
		yrangeAlt = &sequence.ContinuousRange{}
	} else {
		yrangeAlt = c.YAxisSecondary.Range
	}

	if len(c.XAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.XAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		xrange.SetMin(tickMin)
		xrange.SetMax(tickMax)
	} else if xrange.IsZero() {
		xrange.SetMin(minx)
		xrange.SetMax(maxx)
	}

	if len(c.YAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrange.SetMin(tickMin)
		yrange.SetMax(tickMax)
	} else if yrange.IsZero() {
		yrange.SetMin(miny)
		yrange.SetMax(maxy)

		if !c.YAxis.Style.Hidden {
			delta := yrange.GetDelta()
			roundTo := mathutil.RoundTo(delta)
			rmin, rmax := mathutil.RoundDown(yrange.GetMin(), roundTo), mathutil.RoundUp(yrange.GetMax(), roundTo)

			yrange.SetMin(rmin)
			yrange.SetMax(rmax)
		}
	}

	if len(c.YAxisSecondary.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrangeAlt.SetMin(tickMin)
		yrangeAlt.SetMax(tickMax)
	} else if seriesMappedToSecondaryAxis && yrangeAlt.IsZero() {
		yrangeAlt.SetMin(minya)
		yrangeAlt.SetMax(maxya)

		if !c.YAxisSecondary.Style.Hidden {
			delta := yrangeAlt.GetDelta()
			roundTo := mathutil.RoundTo(delta)
			rmin, rmax := mathutil.RoundDown(yrangeAlt.GetMin(), roundTo), mathutil.RoundUp(yrangeAlt.GetMax(), roundTo)
			yrangeAlt.SetMin(rmin)
			yrangeAlt.SetMax(rmax)
		}
	}

	return
}

func (c *Chart) checkRanges(xr, yr, yra sequence.Range) error {
	xDelta := xr.GetDelta()
	if math.IsInf(xDelta, 0) {
		return errors.New("infinite x-range delta")
	}
	if math.IsNaN(xDelta) {
		return errors.New("nan x-range delta")
	}
	if xDelta == 0 {
		return errors.New("zero x-range delta; there needs to be at least (2) values")
	}

	yDelta := yr.GetDelta()
	if math.IsInf(yDelta, 0) {
		return errors.New("infinite y-range delta")
	}
	if math.IsNaN(yDelta) {
		return errors.New("nan y-range delta")
	}

	if c.hasSecondarySeries() {
		yraDelta := yra.GetDelta()
		if math.IsInf(yraDelta, 0) {
			return errors.New("infinite secondary y-range delta")
		}
		if math.IsNaN(yraDelta) {
			return errors.New("nan secondary y-range delta")
		}
	}

	return nil
}

func (c *Chart) getDefaultCanvasBox() render.Box {
	return c.Box()
}

func (c *Chart) getValueFormatters() (x, y, ya dataset.ValueFormatter) {
	for _, s := range c.Series {
		if vfp, isVfp := s.(dataset.ValueFormatterProvider); isVfp {
			sx, sy := vfp.GetValueFormatters()
			if s.GetYAxis() == dataset.YAxisPrimary {
				x = sx
				y = sy
			} else if s.GetYAxis() == dataset.YAxisSecondary {
				x = sx
				ya = sy
			}
		}
	}
	if c.XAxis.ValueFormatter != nil {
		x = c.XAxis.GetValueFormatter()
	}
	if c.YAxis.ValueFormatter != nil {
		y = c.YAxis.GetValueFormatter()
	}
	if c.YAxisSecondary.ValueFormatter != nil {
		ya = c.YAxisSecondary.GetValueFormatter()
	}
	return
}

func (c *Chart) hasAxes() bool {
	return !c.XAxis.Style.Hidden || !c.YAxis.Style.Hidden || !c.YAxisSecondary.Style.Hidden
}

func (c *Chart) getAxesTicks(r render.Renderer, xr, yr, yar sequence.Range, xf, yf, yfa dataset.ValueFormatter) (xticks, yticks, yticksAlt []Tick) {
	if !c.XAxis.Style.Hidden {
		xticks = c.XAxis.GetTicks(r, xr, c.styleDefaultsAxes(), xf)
	}
	if !c.YAxis.Style.Hidden {
		yticks = c.YAxis.GetTicks(r, yr, c.styleDefaultsAxes(), yf)
	}
	if !c.YAxisSecondary.Style.Hidden {
		yticksAlt = c.YAxisSecondary.GetTicks(r, yar, c.styleDefaultsAxes(), yfa)
	}
	return
}

func (c *Chart) getAxesAdjustedCanvasBox(r render.Renderer, canvasBox render.Box, xr, yr, yra sequence.Range, xticks, yticks, yticksAlt []Tick) render.Box {
	axesOuterBox := canvasBox.Clone()
	if !c.XAxis.Style.Hidden {
		axesBounds := c.XAxis.Measure(r, canvasBox, xr, c.styleDefaultsAxes(), xticks)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}
	if !c.YAxis.Style.Hidden {
		axesBounds := c.YAxis.Measure(r, canvasBox, yr, c.styleDefaultsAxes(), yticks)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}
	if !c.YAxisSecondary.Style.Hidden && c.hasSecondarySeries() {
		axesBounds := c.YAxisSecondary.Measure(r, canvasBox, yra, c.styleDefaultsAxes(), yticksAlt)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}

	return canvasBox.OuterConstrain(c.Box(), axesOuterBox)
}

func (c *Chart) setRangeDomains(canvasBox render.Box, xr, yr, yra sequence.Range) (sequence.Range, sequence.Range, sequence.Range) {
	xr.SetDomain(canvasBox.Width())
	yr.SetDomain(canvasBox.Height())
	yra.SetDomain(canvasBox.Height())
	return xr, yr, yra
}

func (c *Chart) hasAnnotationSeries() bool {
	for _, s := range c.Series {
		if as, isAnnotationSeries := s.(dataset.AnnotationSeries); isAnnotationSeries {
			if !as.GetStyle().Hidden {
				return true
			}
		}
	}
	return false
}

func (c *Chart) hasSecondarySeries() bool {
	for _, s := range c.Series {
		if s.GetYAxis() == dataset.YAxisSecondary {
			return true
		}
	}
	return false
}

func (c *Chart) getAnnotationAdjustedCanvasBox(r render.Renderer, canvasBox render.Box, xr, yr, yra sequence.Range, xf, yf, yfa dataset.ValueFormatter) render.Box {
	annotationSeriesBox := canvasBox.Clone()
	for seriesIndex, s := range c.Series {
		if as, isAnnotationSeries := s.(dataset.AnnotationSeries); isAnnotationSeries {
			if !as.GetStyle().Hidden {
				style := c.styleDefaultsSeries(seriesIndex)
				var annotationBounds render.Box
				if as.YAxis == dataset.YAxisPrimary {
					annotationBounds = as.Measure(r, canvasBox, xr, yr, style)
				} else if as.YAxis == dataset.YAxisSecondary {
					annotationBounds = as.Measure(r, canvasBox, xr, yra, style)
				}

				annotationSeriesBox = annotationSeriesBox.Grow(annotationBounds)
			}
		}
	}

	return canvasBox.OuterConstrain(c.Box(), annotationSeriesBox)
}

func (c *Chart) getBackgroundStyle() render.Style {
	return c.Background.InheritFrom(c.styleDefaultsBackground())
}

func (c *Chart) drawBackground(r render.Renderer) {
	render.Box{
		Right:  c.Width(),
		Bottom: c.Height(),
	}.Draw(r, c.getBackgroundStyle())
}

func (c *Chart) getCanvasStyle() render.Style {
	return c.Canvas.InheritFrom(c.styleDefaultsCanvas())
}

func (c *Chart) drawCanvas(r render.Renderer, canvasBox render.Box) {
	canvasBox.Draw(r, c.getCanvasStyle())
}

func (c *Chart) drawAxes(r render.Renderer, canvasBox render.Box, xrange, yrange, yrangeAlt sequence.Range, xticks, yticks, yticksAlt []Tick) {
	if !c.XAxis.Style.Hidden {
		c.XAxis.Render(r, canvasBox, xrange, c.styleDefaultsAxes(), xticks)
	}
	if !c.YAxis.Style.Hidden {
		c.YAxis.Render(r, canvasBox, yrange, c.styleDefaultsAxes(), yticks)
	}
	if !c.YAxisSecondary.Style.Hidden {
		c.YAxisSecondary.Render(r, canvasBox, yrangeAlt, c.styleDefaultsAxes(), yticksAlt)
	}
}

func (c *Chart) drawSeries(r render.Renderer, canvasBox render.Box, xrange, yrange, yrangeAlt sequence.Range, s dataset.Series, seriesIndex int) {
	if !s.GetStyle().Hidden {
		if s.GetYAxis() == dataset.YAxisPrimary {
			s.Render(r, canvasBox, xrange, yrange, c.styleDefaultsSeries(seriesIndex))
		} else if s.GetYAxis() == dataset.YAxisSecondary {
			s.Render(r, canvasBox, xrange, yrangeAlt, c.styleDefaultsSeries(seriesIndex))
		}
	}
}

func (c *Chart) drawTitle(r render.Renderer) {
	if len(c.Title) > 0 && !c.TitleStyle.Hidden {
		r.SetFont(c.TitleStyle.GetFont(c.GetFont()))
		r.SetFontColor(c.TitleStyle.GetFontColor(c.GetColorPalette().TextColor()))
		titleFontSize := c.TitleStyle.GetFontSize(defaultTitleFontSize)
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(c.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (c.Width() >> 1) - (textWidth >> 1)
		titleY := c.TitleStyle.Padding.GetTop(defaultTitleTop) + textHeight

		r.Text(c.Title, titleX, titleY)
	}
}

func (c *Chart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   c.GetColorPalette().BackgroundColor(),
		StrokeColor: c.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: defaultBackgroundStrokeWidth,
	}
}

func (c *Chart) styleDefaultsCanvas() render.Style {
	return render.Style{
		FillColor:   c.GetColorPalette().CanvasColor(),
		StrokeColor: c.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: defaultCanvasStrokeWidth,
	}
}

func (c *Chart) styleDefaultsSeries(seriesIndex int) render.Style {
	return render.Style{
		DotColor:    c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeColor: c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeWidth: defaultSeriesLineWidth,
		Font:        c.GetFont(),
		FontSize:    render.DefaultFontSize,
	}
}

func (c *Chart) styleDefaultsAxes() render.Style {
	return render.Style{
		Font:        c.GetFont(),
		FontColor:   c.GetColorPalette().TextColor(),
		FontSize:    defaultAxisFontSize,
		StrokeColor: c.GetColorPalette().AxisStrokeColor(),
		StrokeWidth: defaultAxisLineWidth,
	}
}

func (c *Chart) styleDefaultsElements() render.Style {
	return render.Style{
		Font: c.GetFont(),
	}
}

// GetColorPalette returns the color palette for the chart.
func (c *Chart) GetColorPalette() render.ColorPalette {
	if c.ColorPalette != nil {
		return c.ColorPalette
	}
	return render.DefaultColorPalette
}

// Box returns the chart bounds as a box.
func (c *Chart) Box() render.Box {
	dpr := c.Background.Padding.GetRight(defaultBackgroundPadding.Right)
	dpb := c.Background.Padding.GetBottom(defaultBackgroundPadding.Bottom)

	return render.Box{
		Top:    c.Background.Padding.GetTop(defaultBackgroundPadding.Top),
		Left:   c.Background.Padding.GetLeft(defaultBackgroundPadding.Left),
		Right:  c.Width() - dpr,
		Bottom: c.Height() - dpb,
	}
}
