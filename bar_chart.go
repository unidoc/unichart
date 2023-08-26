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

// BarChart is a chart that draws bars on a range.
type BarChart struct {
	Title      string
	TitleStyle render.Style

	Font         render.Font
	Background   render.Style
	Canvas       render.Style
	ColorPalette render.ColorPalette

	XAxis render.Style
	YAxis YAxis

	BarWidth     int
	BarSpacing   int
	IsHorizontal bool

	UseBaseValue bool
	BaseValue    float64

	Bars     []dataset.Value
	Elements []render.Renderable

	width  int
	height int
	dpi    float64
}

// DPI returns the DPI for the chart.
func (bc *BarChart) DPI() float64 {
	if bc.dpi == 0 {
		return defaultDPI
	}
	return bc.dpi
}

// SetDPI sets the DPI for the chart.
func (bc *BarChart) SetDPI(dpi float64) {
	bc.dpi = dpi
}

// GetFont returns the text font.
func (bc *BarChart) GetFont() render.Font {
	return bc.Font
}

// Width returns the chart width or the default value.
func (bc *BarChart) Width() int {
	if bc.width == 0 {
		return defaultChartWidth
	}
	return bc.width
}

// SetWidth sets the chart width.
func (bc *BarChart) SetWidth(width int) {
	bc.width = width
}

// Height returns the chart height or the default value.
func (bc *BarChart) Height() int {
	if bc.height == 0 {
		return defaultChartHeight
	}
	return bc.height
}

// SetHeight sets the chart height.
func (bc *BarChart) SetHeight(height int) {
	bc.height = height
}

// GetBarSpacing returns the spacing between bars.
func (bc *BarChart) GetBarSpacing() int {
	if bc.BarSpacing == 0 {
		return defaultBarSpacing
	}
	return bc.BarSpacing
}

// GetBarWidth returns the default bar width.
func (bc *BarChart) GetBarWidth() int {
	if bc.BarWidth == 0 {
		return defaultBarWidth
	}
	return bc.BarWidth
}

// Render renders the chart with the given renderer to the given io.Writer.
func (bc *BarChart) Render(rp render.RendererProvider, w io.Writer) error {
	if len(bc.Bars) == 0 {
		return errors.New("please provide at least one bar")
	}

	r, err := rp(bc.Width(), bc.Height())
	if err != nil {
		return err
	}
	r.SetDPI(bc.DPI())

	bc.drawBackground(r)

	var canvasBox render.Box
	var yt []Tick
	var yr sequence.Range
	var yf dataset.ValueFormatter

	canvasBox = bc.getDefaultCanvasBox()
	yr = bc.getRanges()
	if yr.GetMax()-yr.GetMin() == 0 {
		return fmt.Errorf("invalid data range; cannot be zero")
	}
	yr = bc.setRangeDomains(canvasBox, yr)
	yf = bc.getValueFormatters()

	if bc.hasAxes() {
		yt = bc.getAxesTicks(r, yr, yf)

		// Adjust domain range before adjusting the canvas box
		// if the generated max tick value is exceeding the original max range.
		if len(yt) > 0 {
			if yr.IsDescending() {
				yr.SetMax(yt[0].Value)
			} else {
				yr.SetMax(yt[len(yt)-1].Value)
			}
		}

		yr = bc.setRangeDomains(canvasBox, yr)

		if bc.IsHorizontal {
			canvasBox = bc.getAdjustedHorizontalCanvasBox(r, canvasBox, yr, yt)
		} else {
			canvasBox = bc.getAdjustedCanvasBox(r, canvasBox, yr, yt)
		}
		yr = bc.setRangeDomains(canvasBox, yr)
	}

	bc.drawCanvas(r, canvasBox)

	if bc.IsHorizontal {
		bc.drawHorizontalBars(r, canvasBox, yr)
		bc.drawHorizontalXAxis(r, canvasBox, yr, yt)
		bc.drawHorizontalYAxis(r, canvasBox)
	} else {
		bc.drawBars(r, canvasBox, yr)
		bc.drawXAxis(r, canvasBox)
		bc.drawYAxis(r, canvasBox, yr, yt)
	}

	bc.drawTitle(r)

	for _, a := range bc.Elements {
		a(r, canvasBox, bc.styleDefaultsElements())
	}

	return r.Save(w)
}

func (bc *BarChart) drawCanvas(r render.Renderer, canvasBox render.Box) {
	canvasBox.Draw(r, bc.getCanvasStyle())
}

func (bc *BarChart) getRanges() sequence.Range {
	var yrange sequence.Range
	if bc.YAxis.Range != nil && !bc.YAxis.Range.IsZero() {
		yrange = bc.YAxis.Range
	} else {
		yrange = &sequence.ContinuousRange{}
	}

	if !yrange.IsZero() {
		return yrange
	}

	if len(bc.YAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range bc.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrange.SetMin(tickMin)
		yrange.SetMax(tickMax)
		return yrange
	}

	min, max := math.MaxFloat64, -math.MaxFloat64
	for _, b := range bc.Bars {
		min = math.Min(b.Value, min)
		max = math.Max(b.Value, max)
	}

	yrange.SetMin(min)
	yrange.SetMax(max)

	return yrange
}

func (bc *BarChart) drawBackground(r render.Renderer) {
	render.Box{
		Right:  bc.Width(),
		Bottom: bc.Height(),
	}.Draw(r, bc.getBackgroundStyle())
}

func (bc *BarChart) drawBars(r render.Renderer, canvasBox render.Box, yr sequence.Range) {
	xoffset := canvasBox.Left

	width, spacing, _ := bc.calculateScaledTotalSize(canvasBox)
	bs2 := spacing >> 1

	var barBox render.Box
	var bxl, bxr, by int
	for index, bar := range bc.Bars {
		bxl = xoffset + bs2
		bxr = bxl + width

		barStyle := bar.Style.InheritFrom(bc.styleDefaultsBar(index))
		strokeWidth := barStyle.GetStrokeWidth()
		strokeOffset := int(strokeWidth / 2)

		height := yr.Translate(bar.Value)
		if height == 0 {
			height = int(strokeWidth)
		}
		by = canvasBox.Bottom - height

		if bc.UseBaseValue {
			barBox = render.Box{
				Top:    by,
				Left:   bxl,
				Right:  bxr,
				Bottom: canvasBox.Bottom - yr.Translate(bc.BaseValue) - strokeOffset,
			}
		} else {
			barBox = render.Box{
				Top:    by,
				Left:   bxl,
				Right:  bxr,
				Bottom: canvasBox.Bottom - strokeOffset,
			}
		}

		barBox.Draw(r, barStyle)
		xoffset += width + spacing
	}
}

func (bc *BarChart) drawHorizontalBars(r render.Renderer, canvasBox render.Box, yr sequence.Range) {
	height, spacing, _ := bc.calculateScaledTotalSize(canvasBox)
	bs2 := spacing >> 1

	axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
	axisStyle.WriteToRenderer(r)

	yoffset := canvasBox.Bottom - len(bc.Bars)*(height+spacing)

	maxTextWidth := 0
	for _, bar := range bc.Bars {
		tb := r.MeasureText(bar.Label)

		maxTextWidth = mathutil.MaxInt(maxTextWidth, tb.Width())
	}

	var barBox render.Box
	var byt, byb, bx int
	for index, bar := range bc.Bars {
		byt = yoffset + bs2
		byb = byt + height

		barStyle := bar.Style.InheritFrom(bc.styleDefaultsBar(index))
		strokeWidth := barStyle.GetStrokeWidth()
		strokeOffset := int(strokeWidth / 2)

		width := yr.Translate(bar.Value)
		if width == 0 {
			width = int(strokeWidth)
		}
		bx = canvasBox.Left + defaultYAxisMargin + width

		if bc.UseBaseValue {
			barBox = render.Box{
				Top:    byt,
				Left:   canvasBox.Left + yr.Translate(bc.BaseValue) - strokeOffset,
				Right:  bx,
				Bottom: byb,
			}
		} else {
			barBox = render.Box{
				Top:    byt,
				Left:   canvasBox.Left + strokeOffset,
				Right:  bx,
				Bottom: byb,
			}
		}

		barBox.Draw(r, barStyle)
		yoffset += height + spacing
	}
}

func (bc *BarChart) drawXAxis(r render.Renderer, canvasBox render.Box) {
	if !bc.XAxis.Hidden {
		axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		width, spacing, _ := bc.calculateScaledTotalSize(canvasBox)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Bottom+defaultVerticalTickHeight)
		r.Stroke()

		cursor := canvasBox.Left
		for index, bar := range bc.Bars {
			barLabelBox := render.Box{
				Top:    canvasBox.Bottom + defaultXAxisMargin,
				Left:   cursor,
				Right:  cursor + width + spacing,
				Bottom: bc.Height(),
			}

			if len(bar.Label) > 0 {
				render.Text.DrawWithin(r, bar.Label, barLabelBox, axisStyle)
			}

			axisStyle.WriteToRenderer(r)
			if index < len(bc.Bars)-1 {
				r.MoveTo(barLabelBox.Right, canvasBox.Bottom)
				r.LineTo(barLabelBox.Right, canvasBox.Bottom+defaultVerticalTickHeight)
				r.Stroke()
			}
			cursor += width + spacing
		}
	}
}

func (bc *BarChart) drawHorizontalXAxis(r render.Renderer, canvasBox render.Box, yr sequence.Range, ticks []Tick) {
	if !bc.YAxis.Style.Hidden {
		axisStyle := bc.YAxis.Style.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left-defaultHorizontalTickWidth, canvasBox.Bottom)
		r.Stroke()

		var tx int
		var tb render.Box
		for _, t := range ticks {
			tx = canvasBox.Left + yr.Translate(t.Value)

			axisStyle.GetStrokeOptions().WriteToRenderer(r)
			r.MoveTo(tx, canvasBox.Bottom)
			r.LineTo(tx, canvasBox.Bottom+defaultHorizontalTickWidth)
			r.Stroke()

			axisStyle.GetTextOptions().WriteToRenderer(r)
			tb = r.MeasureText(t.Label)
			render.Text.Draw(r, t.Label, tx-(tb.Width()>>1), canvasBox.Bottom+defaultXAxisMargin+5, axisStyle)
		}
	}
}

func (bc *BarChart) drawYAxis(r render.Renderer, canvasBox render.Box, yr sequence.Range, ticks []Tick) {
	if !bc.YAxis.Style.Hidden {
		axisStyle := bc.YAxis.Style.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		r.MoveTo(canvasBox.Right, canvasBox.Top)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Right, canvasBox.Bottom)
		r.LineTo(canvasBox.Right+defaultHorizontalTickWidth, canvasBox.Bottom)
		r.Stroke()

		var ty int
		var tb render.Box
		for _, t := range ticks {
			ty = canvasBox.Bottom - yr.Translate(t.Value)

			axisStyle.GetStrokeOptions().WriteToRenderer(r)
			r.MoveTo(canvasBox.Right, ty)
			r.LineTo(canvasBox.Right+defaultHorizontalTickWidth, ty)
			r.Stroke()

			axisStyle.GetTextOptions().WriteToRenderer(r)
			tb = r.MeasureText(t.Label)
			render.Text.Draw(r, t.Label, canvasBox.Right+defaultYAxisMargin+5, ty+(tb.Height()>>1), axisStyle)
		}

	}
}

func (bc *BarChart) drawHorizontalYAxis(r render.Renderer, canvasBox render.Box) {
	if !bc.XAxis.Hidden {
		defaultStyle := bc.styleDefaultsAxes()
		defaultStyle.TextHorizontalAlign = render.TextHorizontalAlignRight
		defaultStyle.TextVerticalAlign = render.TextVerticalAlignMiddle

		axisStyle := bc.XAxis.InheritFrom(defaultStyle)

		axisStyle.WriteToRenderer(r)

		width, spacing, _ := bc.calculateScaledTotalSize(canvasBox)

		cursor := canvasBox.Bottom - len(bc.Bars)*(width+spacing)

		r.MoveTo(canvasBox.Left, cursor)
		r.LineTo(canvasBox.Left, canvasBox.Bottom)
		r.Stroke()

		for index, bar := range bc.Bars {
			tb := r.MeasureText(bar.Label)

			barLabelBox := render.Box{
				Top:    cursor + spacing,
				Left:   canvasBox.Left - tb.Width() - (2 * defaultYAxisMargin),
				Right:  canvasBox.Left - defaultYAxisMargin,
				Bottom: cursor + width,
			}

			if len(bar.Label) > 0 {
				render.Text.DrawWithin(r, bar.Label, barLabelBox, axisStyle)
			}

			axisStyle.WriteToRenderer(r)
			if index < len(bc.Bars) {
				r.MoveTo(canvasBox.Left, cursor)
				r.LineTo(canvasBox.Left-defaultHorizontalTickWidth, cursor)
				r.Stroke()
			}
			cursor += width + spacing
		}
	}
}

func (bc *BarChart) drawTitle(r render.Renderer) {
	if len(bc.Title) > 0 && !bc.TitleStyle.Hidden {
		r.SetFont(bc.TitleStyle.GetFont(bc.GetFont()))
		r.SetFontColor(bc.TitleStyle.GetFontColor(bc.GetColorPalette().TextColor()))
		titleFontSize := bc.TitleStyle.GetFontSize(bc.getTitleFontSize())
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(bc.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (bc.Width() >> 1) - (textWidth >> 1)
		titleY := bc.TitleStyle.Padding.GetTop(defaultTitleTop) + textHeight

		r.Text(bc.Title, titleX, titleY)
	}
}

func (bc *BarChart) getCanvasStyle() render.Style {
	return bc.Canvas.InheritFrom(bc.styleDefaultsCanvas())
}

func (bc *BarChart) styleDefaultsCanvas() render.Style {
	return render.Style{
		FillColor:   bc.GetColorPalette().CanvasColor(),
		StrokeColor: bc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: defaultCanvasStrokeWidth,
	}
}

func (bc *BarChart) hasAxes() bool {
	return !bc.YAxis.Style.Hidden
}

func (bc *BarChart) setRangeDomains(canvasBox render.Box, yr sequence.Range) sequence.Range {
	if bc.IsHorizontal {
		yr.SetDomain(canvasBox.Width())
	} else {
		yr.SetDomain(canvasBox.Height())
	}
	return yr
}

func (bc *BarChart) getDefaultCanvasBox() render.Box {
	return bc.box()
}

func (bc *BarChart) getValueFormatters() dataset.ValueFormatter {
	if bc.YAxis.ValueFormatter != nil {
		return bc.YAxis.ValueFormatter
	}
	return dataset.FloatValueFormatter
}

func (bc *BarChart) getAxesTicks(r render.Renderer, yr sequence.Range, yf dataset.ValueFormatter) (yticks []Tick) {
	if !bc.YAxis.Style.Hidden {
		yticks = bc.YAxis.GetTicks(r, yr, bc.styleDefaultsAxes(), yf)
	}
	return
}

func (bc *BarChart) calculateEffectiveBarSpacing(canvasBox render.Box) int {
	canvasLength := canvasBox.Width()
	if bc.IsHorizontal {
		canvasLength = canvasBox.Height()
	}

	totalWithBaseSpacing := bc.calculateTotalBarSize(bc.GetBarWidth(), bc.GetBarSpacing())
	if totalWithBaseSpacing > canvasLength {
		lessBarWidths := canvasLength - (len(bc.Bars) * bc.GetBarWidth()) - defaultHorizontalTickWidth
		if lessBarWidths > 0 {
			return int(math.Ceil(float64(lessBarWidths) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarSpacing()
}

func (bc *BarChart) calculateEffectiveBarSize(canvasBox render.Box, spacing int) int {
	canvasLength := canvasBox.Width()
	if bc.IsHorizontal {
		canvasLength = canvasBox.Height()
	}

	totalWithBaseWidth := bc.calculateTotalBarSize(bc.GetBarWidth(), spacing)
	if totalWithBaseWidth > canvasLength {
		totalLessBarSpacings := canvasLength - (len(bc.Bars) * spacing) - defaultHorizontalTickWidth
		if totalLessBarSpacings > 0 {
			return int(math.Ceil(float64(totalLessBarSpacings) / float64(len(bc.Bars))))
		}
		return 0
	}
	return bc.GetBarWidth()
}

func (bc *BarChart) calculateTotalBarSize(barWidth, spacing int) int {
	return len(bc.Bars) * (barWidth + spacing)
}

func (bc *BarChart) calculateScaledTotalSize(canvasBox render.Box) (size, spacing, total int) {
	spacing = bc.calculateEffectiveBarSpacing(canvasBox)
	size = bc.calculateEffectiveBarSize(canvasBox, spacing)
	total = bc.calculateTotalBarSize(size, spacing)
	return
}

func (bc *BarChart) getAdjustedCanvasBox(r render.Renderer, canvasBox render.Box, yrange sequence.Range, yticks []Tick) render.Box {
	axesOuterBox := canvasBox.Clone()

	_, _, totalWidth := bc.calculateScaledTotalSize(canvasBox)

	if !bc.XAxis.Hidden {
		xaxisHeight := defaultVerticalTickHeight
		axisStyle := bc.XAxis.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		cursor := canvasBox.Left
		for _, bar := range bc.Bars {
			if len(bar.Label) > 0 {
				barLabelBox := render.Box{
					Top:    canvasBox.Bottom + defaultXAxisMargin,
					Left:   cursor,
					Right:  cursor + bc.GetBarWidth() + bc.GetBarSpacing(),
					Bottom: bc.Height(),
				}
				lines := render.Text.WrapFit(r, bar.Label, barLabelBox.Width(), axisStyle)
				linesBox := render.Text.MeasureLines(r, lines, axisStyle)

				xaxisHeight = mathutil.MinInt(linesBox.Height()+(2*defaultXAxisMargin), xaxisHeight)
			}
		}

		xbox := render.Box{
			Top:    canvasBox.Top,
			Left:   canvasBox.Left,
			Right:  canvasBox.Left + totalWidth,
			Bottom: canvasBox.Bottom + defaultXAxisMargin + xaxisHeight,
		}

		axesOuterBox = axesOuterBox.Grow(xbox)
	}

	if !bc.YAxis.Style.Hidden {
		axesBounds := bc.YAxis.Measure(r, canvasBox, yrange, bc.styleDefaultsAxes(), yticks)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}

	return canvasBox.OuterConstrain(bc.box(), axesOuterBox)
}

func (bc *BarChart) getAdjustedHorizontalCanvasBox(r render.Renderer, canvasBox render.Box, yrange sequence.Range, yticks []Tick) render.Box {
	axesOuterBox := canvasBox.Clone()
	size, spacing, totalHeight := bc.calculateScaledTotalSize(canvasBox)

	if len(bc.Title) > 0 && !bc.TitleStyle.Hidden {
		r.SetFont(bc.TitleStyle.GetFont(bc.GetFont()))
		r.SetFontColor(bc.TitleStyle.GetFontColor(bc.GetColorPalette().TextColor()))
		titleFontSize := bc.TitleStyle.GetFontSize(bc.getTitleFontSize())
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(bc.Title)

		tbox := render.Box{
			Top:    canvasBox.Top - textBox.Height() - bc.TitleStyle.Padding.Height(),
			Left:   canvasBox.Left,
			Right:  canvasBox.Right,
			Bottom: canvasBox.Top,
		}

		axesOuterBox = axesOuterBox.Grow(tbox)

		r.ResetStyle()
	}

	if !bc.YAxis.Style.Hidden {
		yAxisWidth := defaultHorizontalTickWidth
		axisStyle := bc.YAxis.Style.InheritFrom(bc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		cursor := canvasBox.Bottom - len(bc.Bars)*(size+spacing)
		for _, bar := range bc.Bars {
			if len(bar.Label) > 0 {
				barLabelBox := render.Box{
					Top:    cursor,
					Left:   0,
					Right:  canvasBox.Left + defaultYAxisMargin,
					Bottom: cursor,
				}

				lines := render.Text.WrapFit(r, bar.Label, barLabelBox.Width(), axisStyle)
				linesBox := render.Text.MeasureLines(r, lines, axisStyle)

				yAxisWidth = mathutil.MaxInt(linesBox.Width()+(2*defaultYAxisMargin), yAxisWidth)
			}

			cursor += size + spacing
		}

		ybox := render.Box{
			Top:    canvasBox.Top,
			Left:   canvasBox.Left - yAxisWidth,
			Right:  canvasBox.Right,
			Bottom: canvasBox.Top + totalHeight,
		}

		axesOuterBox = axesOuterBox.Grow(ybox)
	}

	if !bc.XAxis.Hidden {
		var ltx, rtx int
		var tx int
		var left, right, bottom = math.MaxInt32, 0, 0
		for _, t := range yticks {
			v := t.Value

			tx = canvasBox.Left + yrange.Translate(v)
			tb := render.Text.Measure(r, t.Label, bc.XAxis.GetTextOptions())
			ltx = tx - tb.Width()>>1
			rtx = tx + tb.Width()>>1
			bottom = mathutil.MaxInt(bottom, tb.Height())

			left = mathutil.MinInt(left, ltx)
			right = mathutil.MaxInt(right, rtx)
		}

		xbox := render.Box{
			Top:    canvasBox.Bottom,
			Left:   left,
			Right:  right,
			Bottom: canvasBox.Bottom + defaultXAxisMargin + bottom,
		}

		axesOuterBox = axesOuterBox.Grow(xbox)
	}

	return canvasBox.OuterConstrain(bc.box(), axesOuterBox)
}

// box returns the chart bounds as a box.
func (bc *BarChart) box() render.Box {
	dpr := bc.Background.Padding.GetRight(defaultBackgroundPadding.Right)
	dpb := bc.Background.Padding.GetBottom(defaultBackgroundPadding.Bottom)

	return render.Box{
		Top:    bc.Background.Padding.GetTop(defaultBackgroundPadding.Top),
		Left:   bc.Background.Padding.GetLeft(defaultBackgroundPadding.Left),
		Right:  bc.Width() - dpr,
		Bottom: bc.Height() - dpb,
	}
}

func (bc *BarChart) getBackgroundStyle() render.Style {
	return bc.Background.InheritFrom(bc.styleDefaultsBackground())
}

func (bc *BarChart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   bc.GetColorPalette().BackgroundColor(),
		StrokeColor: bc.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (bc *BarChart) styleDefaultsBar(index int) render.Style {
	return render.Style{
		StrokeColor: bc.GetColorPalette().GetSeriesColor(index),
		FillColor:   bc.GetColorPalette().GetSeriesColor(index),
	}
}

func (bc *BarChart) styleDefaultsTitle() render.Style {
	return bc.TitleStyle.InheritFrom(render.Style{
		FontColor:           bc.GetColorPalette().TextColor(),
		Font:                bc.GetFont(),
		FontSize:            bc.getTitleFontSize(),
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	})
}

func (bc *BarChart) getTitleFontSize() float64 {
	effectiveDimension := mathutil.MinInt(bc.Width(), bc.Height())
	if effectiveDimension >= 2048 {
		return 48
	} else if effectiveDimension >= 1024 {
		return 24
	} else if effectiveDimension >= 512 {
		return 18
	} else if effectiveDimension >= 256 {
		return 12
	}
	return 10
}

func (bc *BarChart) styleDefaultsAxes() render.Style {
	return render.Style{
		StrokeColor:         bc.GetColorPalette().AxisStrokeColor(),
		Font:                bc.GetFont(),
		FontSize:            defaultAxisFontSize,
		FontColor:           bc.GetColorPalette().TextColor(),
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	}
}

func (bc *BarChart) styleDefaultsElements() render.Style {
	return render.Style{
		Font: bc.GetFont(),
	}
}

// GetColorPalette returns the color palette for the chart.
func (bc *BarChart) GetColorPalette() render.ColorPalette {
	if bc.ColorPalette != nil {
		return bc.ColorPalette
	}
	return render.AlternateColorPalette
}
