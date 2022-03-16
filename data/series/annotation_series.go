package series

import (
	"fmt"
	"math"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

const (
	// DefaultAnnotationDeltaWidth is the width of the left triangle out of annotations.
	DefaultAnnotationDeltaWidth = 10

	// DefaultAnnotationFontSize is the font size of annotations.
	DefaultAnnotationFontSize = 10.0
)

var (
	// DefaultAnnotationPadding is the padding around an annotation.
	DefaultAnnotationPadding = render.Box{
		Top:    5,
		Left:   5,
		Right:  5,
		Bottom: 5,
	}

	// DefaultAnnotationFillColor is the default annotation background color.
	DefaultAnnotationFillColor = render.ColorWhite
)

// Interface Assertions.
var (
	_ Series = (*AnnotationSeries)(nil)
)

// FirstValueAnnotation returns an annotation series of just the first value of a value provider as an annotation.
func FirstValueAnnotation(innerSeries data.ValuesProvider, vfs ...data.ValueFormatter) AnnotationSeries {
	var vf data.ValueFormatter
	if len(vfs) > 0 {
		vf = vfs[0]
	} else if typed, isTyped := innerSeries.(data.ValueFormatterProvider); isTyped {
		_, vf = typed.GetValueFormatters()
	} else {
		vf = data.FloatValueFormatter
	}

	var firstValue Value2
	if typed, isTyped := innerSeries.(data.FirstValuesProvider); isTyped {
		firstValue.XValue, firstValue.YValue = typed.GetFirstValues()
		firstValue.Label = vf(firstValue.YValue)
	} else {
		firstValue.XValue, firstValue.YValue = innerSeries.GetValues(0)
		firstValue.Label = vf(firstValue.YValue)
	}

	var seriesName string
	var seriesStyle render.Style
	if typed, isTyped := innerSeries.(Series); isTyped {
		seriesName = fmt.Sprintf("%s - First Value", typed.GetName())
		seriesStyle = typed.GetStyle()
	}

	return AnnotationSeries{
		Name:        seriesName,
		Style:       seriesStyle,
		Annotations: []Value2{firstValue},
	}
}

// LastValueAnnotationSeries returns an annotation series of just the last value of a value provider.
func LastValueAnnotationSeries(innerSeries data.ValuesProvider, vfs ...data.ValueFormatter) AnnotationSeries {
	var vf data.ValueFormatter
	if len(vfs) > 0 {
		vf = vfs[0]
	} else if typed, isTyped := innerSeries.(data.ValueFormatterProvider); isTyped {
		_, vf = typed.GetValueFormatters()
	} else {
		vf = data.FloatValueFormatter
	}

	var lastValue Value2
	if typed, isTyped := innerSeries.(data.LastValuesProvider); isTyped {
		lastValue.XValue, lastValue.YValue = typed.GetLastValues()
		lastValue.Label = vf(lastValue.YValue)
	} else {
		lastValue.XValue, lastValue.YValue = innerSeries.GetValues(innerSeries.Len() - 1)
		lastValue.Label = vf(lastValue.YValue)
	}

	var seriesName string
	var seriesStyle render.Style
	if typed, isTyped := innerSeries.(Series); isTyped {
		seriesName = fmt.Sprintf("%s - Last Value", typed.GetName())
		seriesStyle = typed.GetStyle()
	}

	return AnnotationSeries{
		Name:        seriesName,
		Style:       seriesStyle,
		Annotations: []Value2{lastValue},
	}
}

// BoundedLastValuesAnnotationSeries returns a last value annotation series for a bounded values provider.
func BoundedLastValuesAnnotationSeries(innerSeries data.FullBoundedValuesProvider, vfs ...data.ValueFormatter) AnnotationSeries {
	lvx, lvy1, lvy2 := innerSeries.GetBoundedLastValues()

	var vf data.ValueFormatter
	if len(vfs) > 0 {
		vf = vfs[0]
	} else if typed, isTyped := innerSeries.(data.ValueFormatterProvider); isTyped {
		_, vf = typed.GetValueFormatters()
	} else {
		vf = data.FloatValueFormatter
	}

	label1 := vf(lvy1)
	label2 := vf(lvy2)

	var seriesName string
	var seriesStyle render.Style
	if typed, isTyped := innerSeries.(Series); isTyped {
		seriesName = fmt.Sprintf("%s - Last Values", typed.GetName())
		seriesStyle = typed.GetStyle()
	}

	return AnnotationSeries{
		Name:  seriesName,
		Style: seriesStyle,
		Annotations: []Value2{
			{XValue: lvx, YValue: lvy1, Label: label1},
			{XValue: lvx, YValue: lvy2, Label: label2},
		},
	}
}

// AnnotationSeries is a series of labels on the chart.
type AnnotationSeries struct {
	Name        string
	Style       render.Style
	YAxis       YAxisType
	Annotations []Value2
}

// GetName returns the name of the time series.
func (as AnnotationSeries) GetName() string {
	return as.Name
}

// GetStyle returns the line style.
func (as AnnotationSeries) GetStyle() render.Style {
	return as.Style
}

// GetYAxis returns which YAxis the series draws on.
func (as AnnotationSeries) GetYAxis() YAxisType {
	return as.YAxis
}

func (as AnnotationSeries) annotationStyleDefaults(defaults render.Style) render.Style {
	return render.Style{
		FontColor:   render.DefaultTextColor,
		Font:        defaults.Font,
		FillColor:   DefaultAnnotationFillColor,
		FontSize:    DefaultAnnotationFontSize,
		StrokeColor: defaults.StrokeColor,
		StrokeWidth: defaults.StrokeWidth,
		Padding:     DefaultAnnotationPadding,
	}
}

// Measure returns a bounds box of the series.
func (as AnnotationSeries) Measure(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, defaults render.Style) render.Box {
	box := render.Box{
		Top:    math.MaxInt32,
		Left:   math.MaxInt32,
		Right:  0,
		Bottom: 0,
	}
	if !as.Style.Hidden {
		seriesStyle := as.Style.InheritFrom(as.annotationStyleDefaults(defaults))
		for _, a := range as.Annotations {
			style := a.Style.InheritFrom(seriesStyle)
			lx := canvasBox.Left + xrange.Translate(a.XValue)
			ly := canvasBox.Bottom - yrange.Translate(a.YValue)
			ab := measureAnnotation(r, canvasBox, style, lx, ly, a.Label)
			box.Top = mathutil.MinInt(box.Top, ab.Top)
			box.Left = mathutil.MinInt(box.Left, ab.Left)
			box.Right = mathutil.MaxInt(box.Right, ab.Right)
			box.Bottom = mathutil.MaxInt(box.Bottom, ab.Bottom)
		}
	}
	return box
}

// Render draws the series.
func (as AnnotationSeries) Render(r render.Renderer, canvasBox render.Box, xrange, yrange data.Range, defaults render.Style) {
	if !as.Style.Hidden {
		seriesStyle := as.Style.InheritFrom(as.annotationStyleDefaults(defaults))
		for _, a := range as.Annotations {
			style := a.Style.InheritFrom(seriesStyle)
			lx := canvasBox.Left + xrange.Translate(a.XValue)
			ly := canvasBox.Bottom - yrange.Translate(a.YValue)
			drawAnnotation(r, canvasBox, style, lx, ly, a.Label)
		}
	}
}

// Validate validates the series.
func (as AnnotationSeries) Validate() error {
	if len(as.Annotations) == 0 {
		return fmt.Errorf("annotation series requires annotations to be set and not empty")
	}
	return nil
}

// measureAnnotation measures how big an annotation would be.
func measureAnnotation(r render.Renderer, canvasBox render.Box, style render.Style, lx, ly int, label string) render.Box {
	style.WriteToRenderer(r)
	defer r.ResetStyle()

	textBox := r.MeasureText(label)
	textWidth := textBox.Width()
	textHeight := textBox.Height()
	halfTextHeight := textHeight >> 1

	pt := style.Padding.GetTop(DefaultAnnotationPadding.Top)
	pl := style.Padding.GetLeft(DefaultAnnotationPadding.Left)
	pr := style.Padding.GetRight(DefaultAnnotationPadding.Right)
	pb := style.Padding.GetBottom(DefaultAnnotationPadding.Bottom)

	strokeWidth := style.GetStrokeWidth()

	top := ly - (pt + halfTextHeight)
	right := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth + int(strokeWidth)
	bottom := ly + (pb + halfTextHeight)

	return render.Box{
		Top:    top,
		Left:   lx,
		Right:  right,
		Bottom: bottom,
	}
}

// drawAnnotation draws an anotation with a renderer.
func drawAnnotation(r render.Renderer, canvasBox render.Box, style render.Style, lx, ly int, label string) {
	style.GetTextOptions().WriteToRenderer(r)
	defer r.ResetStyle()

	textBox := r.MeasureText(label)
	textWidth := textBox.Width()
	halfTextHeight := textBox.Height() >> 1

	style.GetFillAndStrokeOptions().WriteToRenderer(r)

	pt := style.Padding.GetTop(DefaultAnnotationPadding.Top)
	pl := style.Padding.GetLeft(DefaultAnnotationPadding.Left)
	pr := style.Padding.GetRight(DefaultAnnotationPadding.Right)
	pb := style.Padding.GetBottom(DefaultAnnotationPadding.Bottom)

	textX := lx + pl + DefaultAnnotationDeltaWidth
	textY := ly + halfTextHeight

	ltx := lx + DefaultAnnotationDeltaWidth
	lty := ly - (pt + halfTextHeight)

	rtx := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth
	rty := ly - (pt + halfTextHeight)

	rbx := lx + pl + pr + textWidth + DefaultAnnotationDeltaWidth
	rby := ly + (pb + halfTextHeight)

	lbx := lx + DefaultAnnotationDeltaWidth
	lby := ly + (pb + halfTextHeight)

	r.MoveTo(lx, ly)
	r.LineTo(ltx, lty)
	r.LineTo(rtx, rty)
	r.LineTo(rbx, rby)
	r.LineTo(lbx, lby)
	r.LineTo(lx, ly)
	r.Close()
	r.FillStroke()

	style.GetTextOptions().WriteToRenderer(r)
	r.Text(label, textX, textY)
}
