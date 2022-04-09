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

// StackedBar is a bar within a StackedBarChart.
type StackedBar struct {
	Name   string
	Width  int
	Values []dataset.Value
}

// GetWidth returns the width of the bar.
func (sb StackedBar) GetWidth() int {
	if sb.Width == 0 {
		return 50
	}
	return sb.Width
}

// StackedBarChart is a chart that draws sections of a bar based on percentages.
type StackedBarChart struct {
	Title      string
	TitleStyle render.Style

	Font         render.Font
	Background   render.Style
	Canvas       render.Style
	ColorPalette render.ColorPalette

	XAxis render.Style
	YAxis render.Style

	BarSpacing   int
	IsHorizontal bool

	Bars     []StackedBar
	Elements []render.Renderable

	width  int
	height int
	dpi    float64
}

// DPI returns the DPI for the chart.
func (sbc StackedBarChart) DPI(defaults ...float64) float64 {
	if sbc.dpi == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return defaultDPI
	}
	return sbc.dpi
}

// SetDPI sets the DPI for the chart.
func (sbc *StackedBarChart) SetDPI(dpi float64) {
	sbc.dpi = dpi
}

// GetFont returns the text font.
func (sbc StackedBarChart) GetFont() render.Font {
	return sbc.Font
}

// Width returns the chart width or the default value.
func (sbc StackedBarChart) Width() int {
	if sbc.width == 0 {
		return defaultChartWidth
	}
	return sbc.width
}

// SetWidth sets the chart width.
func (sbc *StackedBarChart) SetWidth(width int) {
	sbc.width = width
}

// Height returns the chart height or the default value.
func (sbc StackedBarChart) Height() int {
	if sbc.height == 0 {
		return defaultChartWidth
	}
	return sbc.height
}

// SetHeight sets the chart height.
func (sbc *StackedBarChart) SetHeight(height int) {
	sbc.height = height
}

// GetBarSpacing returns the spacing between bars.
func (sbc StackedBarChart) GetBarSpacing() int {
	if sbc.BarSpacing == 0 {
		return 100
	}
	return sbc.BarSpacing
}

// Render renders the chart with the given renderer to the given io.Writer.
func (sbc StackedBarChart) Render(rp render.RendererProvider, w io.Writer) error {
	if len(sbc.Bars) == 0 {
		return errors.New("please provide at least one bar")
	}

	r, err := rp(sbc.Width(), sbc.Height())
	if err != nil {
		return err
	}
	r.SetDPI(sbc.DPI(defaultDPI))

	var canvasBox render.Box
	if sbc.IsHorizontal {
		canvasBox = sbc.getHorizontalAdjustedCanvasBox(r, sbc.getDefaultCanvasBox())
		sbc.drawCanvas(r, canvasBox)
		sbc.drawHorizontalBars(r, canvasBox)
		sbc.drawHorizontalXAxis(r, canvasBox)
		sbc.drawHorizontalYAxis(r, canvasBox)
	} else {
		canvasBox = sbc.getAdjustedCanvasBox(r, sbc.getDefaultCanvasBox())
		sbc.drawCanvas(r, canvasBox)
		sbc.drawBars(r, canvasBox)
		sbc.drawXAxis(r, canvasBox)
		sbc.drawYAxis(r, canvasBox)
	}

	sbc.drawTitle(r)
	for _, a := range sbc.Elements {
		a(r, canvasBox, sbc.styleDefaultsElements())
	}

	return r.Save(w)
}

func (sbc StackedBarChart) drawCanvas(r render.Renderer, canvasBox render.Box) {
	canvasBox.Draw(r, sbc.getCanvasStyle())
}

func (sbc StackedBarChart) drawBars(r render.Renderer, canvasBox render.Box) {
	xoffset := canvasBox.Left
	for _, bar := range sbc.Bars {
		sbc.drawBar(r, canvasBox, xoffset, bar)
		xoffset += (sbc.GetBarSpacing() + bar.GetWidth())
	}
}

func (sbc StackedBarChart) drawHorizontalBars(r render.Renderer, canvasBox render.Box) {
	yOffset := canvasBox.Top
	for _, bar := range sbc.Bars {
		sbc.drawHorizontalBar(r, canvasBox, yOffset, bar)
		yOffset += sbc.GetBarSpacing() + bar.GetWidth()
	}
}

func (sbc StackedBarChart) drawBar(r render.Renderer, canvasBox render.Box, xoffset int, bar StackedBar) int {
	barSpacing2 := sbc.GetBarSpacing() >> 1
	bxl := xoffset + barSpacing2
	bxr := bxl + bar.GetWidth()

	normalizedBarComponents := dataset.Values(bar.Values).Normalize()
	yoffset := canvasBox.Top
	for index, bv := range normalizedBarComponents {
		barHeight := int(math.Ceil(bv.Value * float64(canvasBox.Height())))
		barBox := render.Box{
			Top:    yoffset,
			Left:   bxl,
			Right:  bxr,
			Bottom: mathutil.MinInt(yoffset+barHeight, canvasBox.Bottom-render.DefaultStrokeWidth),
		}

		barBox.Draw(r, bv.Style.InheritFrom(sbc.styleDefaultsStackedBarValue(index)))
		yoffset += barHeight
	}

	// draw the labels
	yoffset = canvasBox.Top
	var lx, ly int
	for index, bv := range normalizedBarComponents {
		barHeight := int(math.Ceil(bv.Value * float64(canvasBox.Height())))

		if len(bv.Label) > 0 {
			lx = bxl + ((bxr - bxl) / 2)
			ly = yoffset + (barHeight / 2)

			bv.Style.InheritFrom(sbc.styleDefaultsStackedBarValue(index)).WriteToRenderer(r)
			tb := r.MeasureText(bv.Label)
			lx = lx - (tb.Width() >> 1)
			ly = ly + (tb.Height() >> 1)

			if lx < 0 {
				lx = 0
			}
			if ly < 0 {
				lx = 0
			}

			r.Text(bv.Label, lx, ly)
		}
		yoffset += barHeight
	}

	return bxr
}

func (sbc StackedBarChart) drawHorizontalBar(r render.Renderer, canvasBox render.Box, yoffset int, bar StackedBar) {
	halfBarSpacing := sbc.GetBarSpacing() >> 1

	boxTop := yoffset + halfBarSpacing
	boxBottom := boxTop + bar.GetWidth()

	normalizedBarComponents := dataset.Values(bar.Values).Normalize()

	xOffset := canvasBox.Right
	for index, bv := range normalizedBarComponents {
		barHeight := int(math.Floor(bv.Value * float64(canvasBox.Width())))
		barBox := render.Box{
			Top:    boxTop,
			Left:   mathutil.MinInt(xOffset-barHeight, canvasBox.Left+render.DefaultStrokeWidth),
			Right:  xOffset,
			Bottom: boxBottom,
		}

		barBox.Draw(r, bv.Style.InheritFrom(sbc.styleDefaultsStackedBarValue(index)))
		xOffset -= barHeight
	}

	// draw the labels
	xOffset = canvasBox.Right
	var lx, ly int
	for index, bv := range normalizedBarComponents {
		barHeight := int(math.Ceil(bv.Value * float64(canvasBox.Width())))

		if len(bv.Label) > 0 {
			lx = xOffset - (barHeight / 2)
			ly = boxTop + ((boxBottom - boxTop) / 2)

			bv.Style.InheritFrom(sbc.styleDefaultsStackedBarValue(index)).WriteToRenderer(r)
			tb := r.MeasureText(bv.Label)
			lx = lx - (tb.Width() >> 1)
			ly = ly + (tb.Height() >> 1)

			if lx < 0 {
				lx = 0
			}
			if ly < 0 {
				lx = 0
			}

			r.Text(bv.Label, lx, ly)
		}
		xOffset -= barHeight
	}
}

func (sbc StackedBarChart) drawXAxis(r render.Renderer, canvasBox render.Box) {
	if !sbc.XAxis.Hidden {
		axisStyle := sbc.XAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Bottom+defaultVerticalTickHeight)
		r.Stroke()

		cursor := canvasBox.Left
		for _, bar := range sbc.Bars {

			barLabelBox := render.Box{
				Top:    canvasBox.Bottom + defaultXAxisMargin,
				Left:   cursor,
				Right:  cursor + bar.GetWidth() + sbc.GetBarSpacing(),
				Bottom: sbc.Height(),
			}
			if len(bar.Name) > 0 {
				render.Text.DrawWithin(r, bar.Name, barLabelBox, axisStyle)
			}
			axisStyle.WriteToRenderer(r)
			r.MoveTo(barLabelBox.Right, canvasBox.Bottom)
			r.LineTo(barLabelBox.Right, canvasBox.Bottom+defaultVerticalTickHeight)
			r.Stroke()
			cursor += bar.GetWidth() + sbc.GetBarSpacing()
		}
	}
}

func (sbc StackedBarChart) drawHorizontalXAxis(r render.Renderer, canvasBox render.Box) {
	if !sbc.XAxis.Hidden {
		axisStyle := sbc.XAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)
		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Bottom+defaultVerticalTickHeight)
		r.Stroke()

		ticks := sequence.LinearRangeWithStep(0.0, 1.0, 0.2)
		for _, t := range ticks {
			axisStyle.GetStrokeOptions().WriteToRenderer(r)
			tx := canvasBox.Left + int(t*float64(canvasBox.Width()))
			r.MoveTo(tx, canvasBox.Bottom)
			r.LineTo(tx, canvasBox.Bottom+defaultVerticalTickHeight)
			r.Stroke()

			axisStyle.GetTextOptions().WriteToRenderer(r)
			text := fmt.Sprintf("%0.0f%%", t*100)

			textBox := r.MeasureText(text)
			textX := tx - (textBox.Width() >> 1)
			textY := canvasBox.Bottom + defaultXAxisMargin + 10

			if t == 1 {
				textX = canvasBox.Right - textBox.Width()
			}

			render.Text.Draw(r, text, textX, textY, axisStyle)
		}
	}
}

func (sbc StackedBarChart) drawYAxis(r render.Renderer, canvasBox render.Box) {
	if !sbc.YAxis.Hidden {
		axisStyle := sbc.YAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)
		r.MoveTo(canvasBox.Right, canvasBox.Top)
		r.LineTo(canvasBox.Right, canvasBox.Bottom)
		r.Stroke()

		r.MoveTo(canvasBox.Right, canvasBox.Bottom)
		r.LineTo(canvasBox.Right+defaultHorizontalTickWidth, canvasBox.Bottom)
		r.Stroke()

		ticks := sequence.LinearRangeWithStep(0.0, 1.0, 0.2)
		for _, t := range ticks {
			axisStyle.GetStrokeOptions().WriteToRenderer(r)
			ty := canvasBox.Bottom - int(t*float64(canvasBox.Height()))
			r.MoveTo(canvasBox.Right, ty)
			r.LineTo(canvasBox.Right+defaultHorizontalTickWidth, ty)
			r.Stroke()

			axisStyle.GetTextOptions().WriteToRenderer(r)
			text := fmt.Sprintf("%0.0f%%", t*100)

			tb := r.MeasureText(text)
			render.Text.Draw(r, text, canvasBox.Right+defaultYAxisMargin+5, ty+(tb.Height()>>1), axisStyle)
		}
	}
}

func (sbc StackedBarChart) drawHorizontalYAxis(r render.Renderer, canvasBox render.Box) {
	if !sbc.YAxis.Hidden {
		axisStyle := sbc.YAxis.InheritFrom(sbc.styleDefaultsHorizontalAxes())
		axisStyle.WriteToRenderer(r)

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left, canvasBox.Top)
		r.Stroke()

		r.MoveTo(canvasBox.Left, canvasBox.Bottom)
		r.LineTo(canvasBox.Left-defaultHorizontalTickWidth, canvasBox.Bottom)
		r.Stroke()

		cursor := canvasBox.Top
		for _, bar := range sbc.Bars {
			barLabelBox := render.Box{
				Top:    cursor,
				Left:   0,
				Right:  canvasBox.Left - defaultYAxisMargin,
				Bottom: cursor + bar.GetWidth() + sbc.GetBarSpacing(),
			}
			if len(bar.Name) > 0 {
				render.Text.DrawWithin(r, bar.Name, barLabelBox, axisStyle)
			}
			axisStyle.WriteToRenderer(r)
			r.MoveTo(canvasBox.Left, barLabelBox.Bottom)
			r.LineTo(canvasBox.Left-defaultHorizontalTickWidth, barLabelBox.Bottom)
			r.Stroke()
			cursor += bar.GetWidth() + sbc.GetBarSpacing()
		}
	}
}

func (sbc StackedBarChart) drawTitle(r render.Renderer) {
	if len(sbc.Title) > 0 && !sbc.TitleStyle.Hidden {
		r.SetFont(sbc.TitleStyle.GetFont(sbc.GetFont()))
		r.SetFontColor(sbc.TitleStyle.GetFontColor(sbc.GetColorPalette().TextColor()))
		titleFontSize := sbc.TitleStyle.GetFontSize(defaultTitleFontSize)
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(sbc.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (sbc.Width() >> 1) - (textWidth >> 1)
		titleY := sbc.TitleStyle.Padding.GetTop(defaultTitleTop) + textHeight

		r.Text(sbc.Title, titleX, titleY)
	}
}

func (sbc StackedBarChart) getCanvasStyle() render.Style {
	return sbc.Canvas.InheritFrom(sbc.styleDefaultsCanvas())
}

func (sbc StackedBarChart) styleDefaultsCanvas() render.Style {
	return render.Style{
		FillColor:   sbc.GetColorPalette().CanvasColor(),
		StrokeColor: sbc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: defaultCanvasStrokeWidth,
	}
}

// GetColorPalette returns the color palette for the chart.
func (sbc StackedBarChart) GetColorPalette() render.ColorPalette {
	if sbc.ColorPalette != nil {
		return sbc.ColorPalette
	}
	return render.AlternateColorPalette
}

func (sbc StackedBarChart) getDefaultCanvasBox() render.Box {
	return sbc.Box()
}

func (sbc StackedBarChart) getAdjustedCanvasBox(r render.Renderer, canvasBox render.Box) render.Box {
	var totalWidth int
	for _, bar := range sbc.Bars {
		totalWidth += bar.GetWidth() + sbc.GetBarSpacing()
	}

	if !sbc.XAxis.Hidden {
		xaxisHeight := defaultVerticalTickHeight

		axisStyle := sbc.XAxis.InheritFrom(sbc.styleDefaultsAxes())
		axisStyle.WriteToRenderer(r)

		cursor := canvasBox.Left
		for _, bar := range sbc.Bars {
			if len(bar.Name) > 0 {
				barLabelBox := render.Box{
					Top:    canvasBox.Bottom + defaultXAxisMargin,
					Left:   cursor,
					Right:  cursor + bar.GetWidth() + sbc.GetBarSpacing(),
					Bottom: sbc.Height(),
				}
				lines := render.Text.WrapFit(r, bar.Name, barLabelBox.Width(), axisStyle)
				linesBox := render.Text.MeasureLines(r, lines, axisStyle)

				xaxisHeight = mathutil.MaxInt(linesBox.Height()+(2*defaultXAxisMargin), xaxisHeight)
			}
		}
		return render.Box{
			Top:    canvasBox.Top,
			Left:   canvasBox.Left,
			Right:  canvasBox.Left + totalWidth,
			Bottom: sbc.Height() - xaxisHeight,
		}
	}
	return render.Box{
		Top:    canvasBox.Top,
		Left:   canvasBox.Left,
		Right:  canvasBox.Left + totalWidth,
		Bottom: canvasBox.Bottom,
	}

}

func (sbc StackedBarChart) getHorizontalAdjustedCanvasBox(r render.Renderer, canvasBox render.Box) render.Box {
	var totalHeight int
	for _, bar := range sbc.Bars {
		totalHeight += bar.GetWidth() + sbc.GetBarSpacing()
	}

	if !sbc.YAxis.Hidden {
		yAxisWidth := defaultHorizontalTickWidth

		axisStyle := sbc.YAxis.InheritFrom(sbc.styleDefaultsHorizontalAxes())
		axisStyle.WriteToRenderer(r)

		cursor := canvasBox.Top
		for _, bar := range sbc.Bars {
			if len(bar.Name) > 0 {
				barLabelBox := render.Box{
					Top:    cursor,
					Left:   0,
					Right:  canvasBox.Left + defaultYAxisMargin,
					Bottom: cursor + bar.GetWidth() + sbc.GetBarSpacing(),
				}
				lines := render.Text.WrapFit(r, bar.Name, barLabelBox.Width(), axisStyle)
				linesBox := render.Text.MeasureLines(r, lines, axisStyle)

				yAxisWidth = mathutil.MaxInt(linesBox.Height()+(2*defaultXAxisMargin), yAxisWidth)
			}
		}
		return render.Box{
			Top:    canvasBox.Top,
			Left:   canvasBox.Left + yAxisWidth,
			Right:  canvasBox.Right,
			Bottom: canvasBox.Top + totalHeight,
		}
	}
	return render.Box{
		Top:    canvasBox.Top,
		Left:   canvasBox.Left,
		Right:  canvasBox.Right,
		Bottom: canvasBox.Top + totalHeight,
	}
}

// Box returns the chart bounds as a box.
func (sbc StackedBarChart) Box() render.Box {
	dpr := sbc.Background.Padding.GetRight(10)
	dpb := sbc.Background.Padding.GetBottom(50)

	return render.Box{
		Top:    sbc.Background.Padding.GetTop(20),
		Left:   sbc.Background.Padding.GetLeft(20),
		Right:  sbc.Width() - dpr,
		Bottom: sbc.Height() - dpb,
	}
}

func (sbc StackedBarChart) styleDefaultsStackedBarValue(index int) render.Style {
	return render.Style{
		StrokeColor: sbc.GetColorPalette().GetSeriesColor(index),
		StrokeWidth: 3.0,
		FillColor:   sbc.GetColorPalette().GetSeriesColor(index),
		FontSize:    sbc.getScaledFontSize(),
		FontColor:   sbc.GetColorPalette().TextColor(),
		Font:        sbc.GetFont(),
	}
}

func (sbc StackedBarChart) styleDefaultsTitle() render.Style {
	return sbc.TitleStyle.InheritFrom(render.Style{
		FontColor:           render.DefaultTextColor,
		Font:                sbc.GetFont(),
		FontSize:            sbc.getTitleFontSize(),
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	})
}

func (sbc StackedBarChart) getScaledFontSize() float64 {
	effectiveDimension := mathutil.MinInt(sbc.Width(), sbc.Height())
	if effectiveDimension >= 2048 {
		return 48.0
	} else if effectiveDimension >= 1024 {
		return 24.0
	} else if effectiveDimension > 512 {
		return 18.0
	} else if effectiveDimension > 256 {
		return 12.0
	}
	return 10.0
}

func (sbc StackedBarChart) getTitleFontSize() float64 {
	effectiveDimension := mathutil.MinInt(sbc.Width(), sbc.Height())
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

func (sbc StackedBarChart) styleDefaultsAxes() render.Style {
	return render.Style{
		StrokeColor:         render.DefaultLineColor,
		Font:                sbc.GetFont(),
		FontSize:            defaultAxisFontSize,
		FontColor:           render.DefaultLineColor,
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	}
}

func (sbc StackedBarChart) styleDefaultsHorizontalAxes() render.Style {
	return render.Style{
		StrokeColor:         render.DefaultLineColor,
		Font:                sbc.GetFont(),
		FontSize:            defaultAxisFontSize,
		FontColor:           render.DefaultLineColor,
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignMiddle,
		TextWrap:            render.TextWrapWord,
	}
}

func (sbc StackedBarChart) styleDefaultsElements() render.Style {
	return render.Style{
		Font: sbc.GetFont(),
	}
}
