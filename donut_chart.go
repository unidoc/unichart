package chart

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/unidoc/unichart/data/series"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// DonutChart is a chart that draws sections of a circle based on percentages with an hole.
type DonutChart struct {
	Title      string
	TitleStyle render.Style

	ColorPalette render.ColorPalette

	Width  int
	Height int
	DPI    float64

	Font       render.Font
	Background render.Style
	Canvas     render.Style
	SliceStyle render.Style

	Values   []series.Value
	Elements []render.Renderable
}

// GetDPI returns the dpi for the chart.
func (pc DonutChart) GetDPI(defaults ...float64) float64 {
	if pc.DPI == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return defaultDPI
	}
	return pc.DPI
}

// GetFont returns the text font.
func (pc DonutChart) GetFont() render.Font {
	return pc.Font
}

// GetWidth returns the chart width or the default value.
func (pc DonutChart) GetWidth() int {
	if pc.Width == 0 {
		return defaultChartWidth
	}
	return pc.Width
}

// GetHeight returns the chart height or the default value.
func (pc DonutChart) GetHeight() int {
	if pc.Height == 0 {
		return defaultChartWidth
	}
	return pc.Height
}

// Render renders the chart with the given renderer to the given io.Writer.
func (pc DonutChart) Render(rp render.RendererProvider, w io.Writer) error {
	if len(pc.Values) == 0 {
		return errors.New("please provide at least one value")
	}

	r, err := rp(pc.GetWidth(), pc.GetHeight())
	if err != nil {
		return err
	}
	r.SetDPI(pc.GetDPI(defaultDPI))

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

func (pc DonutChart) drawBackground(r render.Renderer) {
	render.Box{
		Right:  pc.GetWidth(),
		Bottom: pc.GetHeight(),
	}.Draw(r, pc.getBackgroundStyle())
}

func (pc DonutChart) drawCanvas(r render.Renderer, canvasBox render.Box) {
	canvasBox.Draw(r, pc.getCanvasStyle())
}

func (pc DonutChart) drawTitle(r render.Renderer) {
	if len(pc.Title) > 0 && !pc.TitleStyle.Hidden {
		render.Text.DrawWithin(r, pc.Title, pc.Box(), pc.styleDefaultsTitle())
	}
}

func (pc DonutChart) drawSlices(r render.Renderer, canvasBox render.Box, values []series.Value) {
	cx, cy := canvasBox.Center()
	diameter := mathutil.MinInt(canvasBox.Width(), canvasBox.Height())
	radius := float64(diameter>>1) / 1.1
	labelRadius := (radius * 2.83) / 3.0

	// draw the donut slices
	var rads, delta, delta2, total float64
	var lx, ly int

	if len(values) == 1 {
		pc.styleDonutChartValue(0).WriteToRenderer(r)
		r.MoveTo(cx, cy)
		r.Circle(radius, cx, cy)
	} else {
		for index, v := range values {
			v.Style.InheritFrom(pc.styleDonutChartValue(index)).WriteToRenderer(r)
			r.MoveTo(cx, cy)
			rads = mathutil.PercentToRadians(total)
			delta = mathutil.PercentToRadians(v.Value)

			r.ArcTo(cx, cy, (radius / 1.25), (radius / 1.25), rads, delta)

			r.LineTo(cx, cy)
			r.Close()
			r.FillStroke()
			total = total + v.Value
		}
	}

	//making the donut hole
	v := series.Value{Value: 100, Label: "center"}
	styletemp := pc.SliceStyle.InheritFrom(render.Style{
		StrokeColor: render.ColorWhite, StrokeWidth: 4.0, FillColor: render.ColorWhite, FontColor: render.ColorWhite, //Font:        pc.GetFont(),//FontSize:    pc.getScaledFontSize(),
	})
	v.Style.InheritFrom(styletemp).WriteToRenderer(r)
	r.MoveTo(cx, cy)
	r.ArcTo(cx, cy, (radius / 3.5), (radius / 3.5), mathutil.DegreesToRadians(0), mathutil.DegreesToRadians(359))
	r.LineTo(cx, cy)
	r.Close()
	r.FillStroke()

	// draw the labels
	total = 0
	for index, v := range values {
		v.Style.InheritFrom(pc.styleDonutChartValue(index)).WriteToRenderer(r)
		if len(v.Label) > 0 {
			delta2 = mathutil.PercentToRadians(total + (v.Value / 2.0))
			delta2 = mathutil.RadiansAdd(delta2, math.Pi/2.0)
			lx, ly = mathutil.CirclePoint(cx, cy, labelRadius, delta2)

			tb := r.MeasureText(v.Label)
			lx = lx - (tb.Width() >> 1)
			ly = ly + (tb.Height() >> 1)

			r.Text(v.Label, lx, ly)
		}
		total = total + v.Value
	}
}

func (pc DonutChart) finalizeValues(values []series.Value) ([]series.Value, error) {
	finalValues := series.Values(values).Normalize()
	if len(finalValues) == 0 {
		return nil, fmt.Errorf("donut chart must contain at least (1) non-zero value")
	}
	return finalValues, nil
}

func (pc DonutChart) getDefaultCanvasBox() render.Box {
	return pc.Box()
}

func (pc DonutChart) getCircleAdjustedCanvasBox(canvasBox render.Box) render.Box {
	circleDiameter := mathutil.MinInt(canvasBox.Width(), canvasBox.Height())

	square := render.Box{
		Right:  circleDiameter,
		Bottom: circleDiameter,
	}

	return canvasBox.Fit(square)
}

func (pc DonutChart) getBackgroundStyle() render.Style {
	return pc.Background.InheritFrom(pc.styleDefaultsBackground())
}

func (pc DonutChart) getCanvasStyle() render.Style {
	return pc.Canvas.InheritFrom(pc.styleDefaultsCanvas())
}

func (pc DonutChart) styleDefaultsCanvas() render.Style {
	return render.Style{
		FillColor:   pc.GetColorPalette().CanvasColor(),
		StrokeColor: pc.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (pc DonutChart) styleDefaultsDonutChartValue() render.Style {
	return render.Style{
		StrokeColor: pc.GetColorPalette().TextColor(),
		StrokeWidth: 4.0,
		FillColor:   pc.GetColorPalette().TextColor(),
	}
}

func (pc DonutChart) styleDonutChartValue(index int) render.Style {
	return pc.SliceStyle.InheritFrom(render.Style{
		StrokeColor: render.ColorWhite,
		StrokeWidth: 4.0,
		FillColor:   pc.GetColorPalette().GetSeriesColor(index),
		FontSize:    pc.getScaledFontSize(),
		FontColor:   pc.GetColorPalette().TextColor(),
		Font:        pc.GetFont(),
	})
}

func (pc DonutChart) getScaledFontSize() float64 {
	effectiveDimension := mathutil.MinInt(pc.GetWidth(), pc.GetHeight())
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

func (pc DonutChart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   pc.GetColorPalette().BackgroundColor(),
		StrokeColor: pc.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (pc DonutChart) styleDefaultsElements() render.Style {
	return render.Style{
		Font: pc.GetFont(),
	}
}

func (pc DonutChart) styleDefaultsTitle() render.Style {
	return pc.TitleStyle.InheritFrom(render.Style{
		FontColor:           pc.GetColorPalette().TextColor(),
		Font:                pc.GetFont(),
		FontSize:            pc.getTitleFontSize(),
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignTop,
		TextWrap:            render.TextWrapWord,
	})
}

func (pc DonutChart) getTitleFontSize() float64 {
	effectiveDimension := mathutil.MinInt(pc.GetWidth(), pc.GetHeight())
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
func (pc DonutChart) GetColorPalette() render.ColorPalette {
	if pc.ColorPalette != nil {
		return pc.ColorPalette
	}
	return render.AlternateColorPalette
}

// Box returns the chart bounds as a box.
func (pc DonutChart) Box() render.Box {
	dpr := pc.Background.Padding.GetRight(defaultBackgroundPadding.Right)
	dpb := pc.Background.Padding.GetBottom(defaultBackgroundPadding.Bottom)

	return render.Box{
		Top:    pc.Background.Padding.GetTop(defaultBackgroundPadding.Top),
		Left:   pc.Background.Padding.GetLeft(defaultBackgroundPadding.Left),
		Right:  pc.GetWidth() - dpr,
		Bottom: pc.GetHeight() - dpb,
	}
}
