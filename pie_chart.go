package unichart

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// PieChart is a chart that draws sections of a circle based on percentages.
type PieChart struct {
	Title      string
	TitleStyle render.Style

	Font         render.Font
	Background   render.Style
	Canvas       render.Style
	SliceStyle   render.Style
	ColorPalette render.ColorPalette

	Values   []dataset.Value
	Elements []render.Renderable

	width  int
	height int
	dpi    float64
}

// DPI returns the DPI for the chart.
func (pc *PieChart) DPI(defaults ...float64) float64 {
	if pc.dpi == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return defaultDPI
	}
	return pc.dpi
}

// SetDPI sets the DPI for the chart.
func (pc *PieChart) SetDPI(dpi float64) {
	pc.dpi = dpi
}

// GetFont returns the text font.
func (pc *PieChart) GetFont() render.Font {
	return pc.Font
}

// Width returns the chart width.
func (pc *PieChart) Width() int {
	if pc.width == 0 {
		return defaultChartWidth
	}
	return pc.width
}

// SetWidth sets the chart width.
func (pc *PieChart) SetWidth(width int) {
	pc.width = width
}

// Height returns the chart height.
func (pc *PieChart) Height() int {
	if pc.height == 0 {
		return defaultChartHeight
	}
	return pc.height
}

// SetHeight sets the chart height.
func (pc *PieChart) SetHeight(height int) {
	pc.height = height
}

// Render renders the chart with the given renderer to the given io.Writer.
func (pc *PieChart) Render(rp render.RendererProvider, w io.Writer) error {
	if len(pc.Values) == 0 {
		return errors.New("please provide at least one value")
	}

	r, err := rp(pc.Width(), pc.Height())
	if err != nil {
		return err
	}
	r.SetDPI(pc.DPI(defaultDPI))

	canvasBox := pc.getDefaultCanvasBox()
	canvasBox = pc.getCircleAdjustedCanvasBox(canvasBox)

	pc.drawBackground(r)
	pc.drawCanvas(r, canvasBox)

	finalValues, err := pc.finalizeValues(pc.Values)
	if err != nil {
		return err
	}
	pc.drawSlices(r, canvasBox, finalValues)
	pc.drawTitle(r)
	for _, a := range pc.Elements {
		a(r, canvasBox, pc.styleDefaultsElements())
	}

	return r.Save(w)
}

func (pc *PieChart) drawBackground(r render.Renderer) {
	render.Box{
		Right:  pc.Width(),
		Bottom: pc.Height(),
	}.Draw(r, pc.getBackgroundStyle())
}

func (pc *PieChart) drawCanvas(r render.Renderer, canvasBox render.Box) {
	canvasBox.Draw(r, pc.getCanvasStyle())
}

func (pc *PieChart) drawTitle(r render.Renderer) {
	if len(pc.Title) > 0 && !pc.TitleStyle.Hidden {
		render.Text.DrawWithin(r, pc.Title, pc.Box(), pc.styleDefaultsTitle())
	}
}

func (pc *PieChart) drawSlices(r render.Renderer, canvasBox render.Box, values []dataset.Value) {
	cx, cy := canvasBox.Center()
	diameter := mathutil.MinInt(canvasBox.Width(), canvasBox.Height())
	radius := float64(diameter >> 1)
	labelRadius := (radius * 2.0) / 3.0

	// Draw the pie slices.
	var rads, delta, delta2, total float64
	var lx, ly int

	if len(values) == 1 {
		pc.stylePieChartValue(0).WriteToRenderer(r)
		r.MoveTo(cx, cy)
		r.Circle(radius, cx, cy)
		r.FillStroke()
	} else {
		for index, v := range values {
			v.Style.InheritFrom(pc.stylePieChartValue(index)).WriteToRenderer(r)

			r.MoveTo(cx, cy)
			rads = mathutil.PercentToRadians(total)
			delta = mathutil.PercentToRadians(v.Value)

			r.ArcTo(cx, cy, radius, radius, rads, delta)

			r.LineTo(cx, cy)
			r.FillStroke()
			r.Close()
			total = total + v.Value
		}
	}

	// Draw the labels.
	total = 0
	for index, v := range values {
		v.Style.InheritFrom(pc.stylePieChartValue(index)).WriteToRenderer(r)
		if len(v.Label) > 0 {
			delta2 = mathutil.PercentToRadians(total + (v.Value / 2.0))
			delta2 = mathutil.RadiansAdd(delta2, math.Pi/2.0)
			lx, ly = mathutil.CirclePoint(cx, cy, labelRadius, delta2)

			tb := r.MeasureText(v.Label)
			lx = lx - (tb.Width() >> 1)
			ly = ly + (tb.Height() >> 1)

			if lx < 0 {
				lx = 0
			}
			if ly < 0 {
				lx = 0
			}

			r.Text(v.Label, lx, ly)
		}
		total = total + v.Value
	}
}

func (pc *PieChart) finalizeValues(values []dataset.Value) ([]dataset.Value, error) {
	finalValues := dataset.Values(values).Normalize()
	if len(finalValues) == 0 {
		return nil, fmt.Errorf("pie chart must contain at least (1) non-zero value")
	}
	return finalValues, nil
}

func (pc *PieChart) getDefaultCanvasBox() render.Box {
	return pc.Box()
}

func (pc *PieChart) getCircleAdjustedCanvasBox(canvasBox render.Box) render.Box {
	circleDiameter := mathutil.MinInt(canvasBox.Width(), canvasBox.Height())

	square := render.Box{
		Right:  circleDiameter,
		Bottom: circleDiameter,
	}

	return canvasBox.Fit(square)
}

func (pc *PieChart) getBackgroundStyle() render.Style {
	return pc.Background.InheritFrom(pc.styleDefaultsBackground())
}

func (pc *PieChart) getCanvasStyle() render.Style {
	return pc.Canvas.InheritFrom(pc.styleDefaultsCanvas())
}

func (pc *PieChart) styleDefaultsCanvas() render.Style {
	return render.Style{
		FillColor:   pc.GetColorPalette().CanvasColor(),
		StrokeColor: pc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (pc *PieChart) styleDefaultsPieChartValue() render.Style {
	return render.Style{
		StrokeColor: pc.GetColorPalette().TextColor(),
		StrokeWidth: 5.0,
		FillColor:   pc.GetColorPalette().TextColor(),
	}
}

func (pc *PieChart) stylePieChartValue(index int) render.Style {
	return pc.SliceStyle.InheritFrom(render.Style{
		StrokeColor: render.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   pc.GetColorPalette().GetSeriesColor(index),
		FontSize:    pc.getScaledFontSize(),
		FontColor:   pc.GetColorPalette().TextColor(),
		Font:        pc.GetFont(),
	})
}

func (pc *PieChart) getScaledFontSize() float64 {
	effectiveDimension := mathutil.MinInt(pc.Width(), pc.Height())
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

func (pc *PieChart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   pc.GetColorPalette().BackgroundColor(),
		StrokeColor: pc.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (pc *PieChart) styleDefaultsElements() render.Style {
	return render.Style{
		Font: pc.GetFont(),
	}
}

func (pc *PieChart) styleDefaultsTitle() render.Style {
	return pc.TitleStyle.InheritFrom(render.Style{
		FontColor:           pc.GetColorPalette().TextColor(),
		Font:                pc.GetFont(),
		FontSize:            pc.getTitleFontSize(),
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	})
}

func (pc *PieChart) getTitleFontSize() float64 {
	effectiveDimension := mathutil.MinInt(pc.Width(), pc.Height())
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

// GetColorPalette returns the color palette for the chart.
func (pc *PieChart) GetColorPalette() render.ColorPalette {
	if pc.ColorPalette != nil {
		return pc.ColorPalette
	}
	return render.AlternateColorPalette
}

// Box returns the chart bounds as a box.
func (pc *PieChart) Box() render.Box {
	dpr := pc.Background.Padding.GetRight(defaultBackgroundPadding.Right)
	dpb := pc.Background.Padding.GetBottom(defaultBackgroundPadding.Bottom)

	return render.Box{
		Top:    pc.Background.Padding.GetTop(defaultBackgroundPadding.Top),
		Left:   pc.Background.Padding.GetLeft(defaultBackgroundPadding.Left),
		Right:  pc.Width() - dpr,
		Bottom: pc.Height() - dpb,
	}
}
