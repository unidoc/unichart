package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/model"
)

// CircularProgressChart is a component that will render circular progress bar component.
type CircularProgressChart struct {
	BackgroundStyle render.Style
	ForegroundStyle render.Style
	LabelStyle      render.Style
	ColorPalette    render.ColorPalette

	// If the value is true then the progress bar would rendered counter clockwise.
	Reversed bool

	// Size of the progress bar in width and height.
	size int
	dpi  float64

	// Progress values which should be between 0.0 - 1.0.
	progress float64

	// Label to be displayed in the center of the chart.
	label string
}

// SetProgress set the progress that will represented by this chart.
// Expected value should be float number between 0.0 - 1.0.
func (cp *CircularProgressChart) SetProgress(progress float64) {
	cp.progress = math.Max(0.0, math.Min(progress, 1.0))
}

// GetProgress returns the progress represented by this chart.
// Returned value should be a float number between 0.0 - 1.0.
func (cp *CircularProgressChart) GetProgress() float64 {
	return cp.progress
}

// DPI returns the DPI of the progress bar.
func (cp *CircularProgressChart) DPI() float64 {
	if cp.dpi == 0 {
		return defaultDPI
	}
	return cp.dpi
}

// SetDPI sets the DPI for the progrss bar.
func (cp *CircularProgressChart) SetDPI(dpi float64) {
	cp.dpi = dpi
}

// Width returns the chart size or the default value.
func (cp *CircularProgressChart) Size() int {
	if cp.size == 0 {
		return defaultChartWidth
	}
	return cp.size
}

// SetWidth sets the chart size.
func (cp *CircularProgressChart) SetSize(size int) {
	cp.size = size
}

func (cp *CircularProgressChart) SetLabel(label string) {
	cp.label = label
}

func (cp *CircularProgressChart) GetLabel() string {
	return cp.label
}

// Width returns the chart width.
func (cp *CircularProgressChart) Width() int {
	return cp.Size()
}

// SetWidth sets the chart width.
func (cp *CircularProgressChart) SetWidth(width int) {
}

// Height returns the chart height.
func (cp *CircularProgressChart) Height() int {
	return cp.Size()
}

// SetHeight sets the chart height.
func (cp *CircularProgressChart) SetHeight(height int) {
}

func (cp *CircularProgressChart) getBackgroundStyle() render.Style {
	return cp.BackgroundStyle.InheritFrom(cp.styleDefaultsBackground())
}

func (cp *CircularProgressChart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   cp.GetColorPalette().BackgroundColor(),
		StrokeColor: cp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (cp *CircularProgressChart) getForegroundStyle() render.Style {
	return cp.ForegroundStyle.InheritFrom(cp.styleDefaultsForeground())
}

func (cp *CircularProgressChart) styleDefaultsForeground() render.Style {
	return render.Style{
		FillColor:   cp.GetColorPalette().BackgroundColor(),
		StrokeColor: cp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (cp *CircularProgressChart) getLabelStyle() render.Style {
	return cp.LabelStyle.InheritFrom(cp.styleDefaultsLabel())
}

func (cp *CircularProgressChart) styleDefaultsLabel() render.Style {
	return render.Style{
		Font:                model.DefaultFont(),
		FontSize:            render.DefaultFontSize,
		FontColor:           cp.getForegroundStyle().StrokeColor,
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignMiddle,
	}
}

// GetColorPalette returns the color palette for the chart.
func (cp *CircularProgressChart) GetColorPalette() render.ColorPalette {
	if cp.ColorPalette != nil {
		return cp.ColorPalette
	}
	return render.AlternateColorPalette
}

func (cp *CircularProgressChart) drawBackground(r render.Renderer) {
	bgStyle := cp.getBackgroundStyle()

	if bgStyle.Hidden {
		return
	}

	radius := cp.Size() / 2.0

	r.Circle(float64(radius), radius, radius)
	r.SetFillColor(bgStyle.FillColor)
	r.SetStrokeColor(bgStyle.StrokeColor)
	r.SetStrokeWidth(bgStyle.StrokeWidth)
	r.FillStroke()
}

func (cp *CircularProgressChart) drawForeground(r render.Renderer) {
	fgStyle := cp.getForegroundStyle()

	if fgStyle.Hidden {
		return
	}

	radius := float64(cp.Size()) / 2.0
	progressDeg := cp.progress * 360.0

	if cp.Reversed {
		progressDeg = -progressDeg
	}

	r.MoveTo(int(radius), 0)
	r.ArcTo(int(radius), int(radius), radius, radius, mathutil.DegreesToRadians(-90), mathutil.DegreesToRadians(progressDeg))
	r.SetStrokeColor(fgStyle.StrokeColor)
	r.SetStrokeWidth(fgStyle.StrokeWidth)
	r.Stroke()
}

func (cp *CircularProgressChart) drawLabel(r render.Renderer) {
	labelStyle := cp.getLabelStyle()

	if labelStyle.Hidden || cp.label == "" {
		return
	}

	fgStrokeWidth := int(labelStyle.StrokeWidth)
	render.Text.DrawWithin(r, cp.label, render.NewBox(fgStrokeWidth, fgStrokeWidth, cp.Size()-fgStrokeWidth, cp.Size()-fgStrokeWidth), labelStyle)
}

// Render renders the progrss bar with the given renderer to the given io.Writer.
func (cp *CircularProgressChart) Render(rp render.RendererProvider, w io.Writer) error {
	r, err := rp(cp.Size(), cp.Size())
	if err != nil {
		return err
	}
	r.SetDPI(cp.DPI())

	cp.drawBackground(r)
	cp.drawForeground(r)
	cp.drawLabel(r)

	return r.Save(w)
}
