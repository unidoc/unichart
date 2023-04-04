package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/model"
)

// CircularProgressBar is a component that will render circular progress bar component.
type CircularProgressBar struct {
	// BackgroundStyle is the style for the background bar.
	BackgroundStyle render.Style

	// ForegroundStyle is the style for the foreground bar.
	ForegroundStyle render.Style

	// LabelStyle is the style for the label that will displayed in the center of the progress bar.
	LabelStyle render.Style

	// ColorPalette is the color pallete that could be used to add colors in this progress bar
	ColorPalette render.ColorPalette

	// Reversed is a flag where if the value is true then the progress bar would rendered counter clockwise.
	Reversed bool

	// size is the size of the progress bar in width and height.
	size int
	dpi  float64

	// progress is the progress bar values which should be between 0.0 - 1.0.
	progress float64

	// label is the text to be displayed in the center of the chart.
	label string
}

// SetProgress set the progress that will represented by this chart.
// Expected value should be float number between 0.0 - 1.0.
func (cp *CircularProgressBar) SetProgress(progress float64) {
	cp.progress = math.Max(0.0, math.Min(progress, 1.0))
}

// GetProgress returns the progress represented by this chart.
// Returned value should be a float number between 0.0 - 1.0.
func (cp *CircularProgressBar) GetProgress() float64 {
	return cp.progress
}

// DPI returns the DPI of the progress bar.
func (cp *CircularProgressBar) DPI() float64 {
	if cp.dpi == 0 {
		return defaultDPI
	}
	return cp.dpi
}

// SetDPI sets the DPI for the progress bar.
func (cp *CircularProgressBar) SetDPI(dpi float64) {
	cp.dpi = dpi
}

// Size returns the chart size or the default value.
func (cp *CircularProgressBar) Size() int {
	if cp.size == 0 {
		return defaultChartWidth
	}
	return cp.size
}

// SetSize sets the chart size.
func (cp *CircularProgressBar) SetSize(size int) {
	cp.size = size
}

// SetLabel sets the label that would be displayed in the center of the progress bar.
func (cp *CircularProgressBar) SetLabel(label string) {
	cp.label = label
}

// GetLabel returns the label displayed in the center of the progress bar.
func (cp *CircularProgressBar) GetLabel() string {
	return cp.label
}

// Width returns the chart width.
func (cp *CircularProgressBar) Width() int {
	return cp.Size()
}

// SetWidth method is exists to fuifill the requirements of render.ChartRenderable interface.
// To set width or height of this circular progress bar, use SetSize instead.
func (cp *CircularProgressBar) SetWidth(width int) {
}

// Height returns the chart height.
func (cp *CircularProgressBar) Height() int {
	return cp.Size()
}

// SetHeight method is exists to fuifill the requirements of render.ChartRenderable interface.
// To set width or height of this circular progress bar, use SetSize instead.
func (cp *CircularProgressBar) SetHeight(height int) {
}

func (cp *CircularProgressBar) getBackgroundStyle() render.Style {
	return cp.BackgroundStyle.InheritFrom(cp.styleDefaultsBackground())
}

func (cp *CircularProgressBar) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   cp.GetColorPalette().BackgroundColor(),
		StrokeColor: cp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (cp *CircularProgressBar) getForegroundStyle() render.Style {
	return cp.ForegroundStyle.InheritFrom(cp.styleDefaultsForeground())
}

func (cp *CircularProgressBar) styleDefaultsForeground() render.Style {
	return render.Style{
		FillColor:   cp.GetColorPalette().BackgroundColor(),
		StrokeColor: cp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (cp *CircularProgressBar) getLabelStyle() render.Style {
	return cp.LabelStyle.InheritFrom(cp.styleDefaultsLabel())
}

func (cp *CircularProgressBar) styleDefaultsLabel() render.Style {
	return render.Style{
		Font:                model.DefaultFont(),
		FontSize:            render.DefaultFontSize,
		FontColor:           cp.getForegroundStyle().StrokeColor,
		TextHorizontalAlign: render.TextHorizontalAlignCenter,
		TextVerticalAlign:   render.TextVerticalAlignMiddle,
	}
}

// GetColorPalette returns the color palette for the chart.
func (cp *CircularProgressBar) GetColorPalette() render.ColorPalette {
	if cp.ColorPalette != nil {
		return cp.ColorPalette
	}
	return render.AlternateColorPalette
}

func (cp *CircularProgressBar) drawBackground(r render.Renderer) {
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

func (cp *CircularProgressBar) drawForeground(r render.Renderer) {
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

func (cp *CircularProgressBar) drawLabel(r render.Renderer) {
	labelStyle := cp.getLabelStyle()

	if labelStyle.Hidden || cp.label == "" {
		return
	}

	fgStrokeWidth := int(labelStyle.StrokeWidth)
	render.Text.DrawWithin(r, cp.label, render.NewBox(fgStrokeWidth, fgStrokeWidth, cp.Size()-fgStrokeWidth, cp.Size()-fgStrokeWidth), labelStyle)
}

// Render renders the progress bar with the given renderer to the given io.Writer.
func (cp *CircularProgressBar) Render(rp render.RendererProvider, w io.Writer) error {
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
