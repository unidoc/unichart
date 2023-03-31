package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// LinearProgressChart is a component that will render progress bar component.
type LinearProgressChart struct {
	Background   render.Style
	Foreground   render.Style
	ColorPalette render.ColorPalette

	RoundedEdgeStart bool
	RoundedEdgeEnd   bool

	height int
	width  int
	dpi    float64

	progress float64
}

// SetProgress set the progress that will represented by this chart.
// Expected value should be float number between 0.0 - 1.0.
func (lp *LinearProgressChart) SetProgress(progress float64) {
	lp.progress = math.Max(0.0, math.Min(progress, 1.0))
}

// GetProgress returns the progress represented by this chart.
// Returned value should be a float number between 0.0 - 1.0.
func (lp *LinearProgressChart) GetProgress() float64 {
	return lp.progress
}

// DPI returns the DPI of the progress bar.
func (lp *LinearProgressChart) DPI() float64 {
	if lp.dpi == 0 {
		return defaultDPI
	}
	return lp.dpi
}

// SetDPI sets the DPI for the progrss bar.
func (lp *LinearProgressChart) SetDPI(dpi float64) {
	lp.dpi = dpi
}

// Width returns the chart width or the default value.
func (lp *LinearProgressChart) Width() int {
	if lp.width == 0 {
		return defaultChartWidth
	}
	return lp.width
}

// SetWidth sets the chart width.
func (lp *LinearProgressChart) SetWidth(width int) {
	lp.width = width
}

// Height returns the chart height or the default value.
func (lp *LinearProgressChart) Height() int {
	if lp.height == 0 {
		return defaultChartHeight
	}
	return lp.height
}

// SetHeight sets the chart height.
func (lp *LinearProgressChart) SetHeight(height int) {
	lp.height = height
}

func (lp *LinearProgressChart) getBackgroundStyle() render.Style {
	return lp.Background.InheritFrom(lp.styleDefaultsBackground())
}

func (lp *LinearProgressChart) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   lp.GetColorPalette().BackgroundColor(),
		StrokeColor: lp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (lp *LinearProgressChart) getForegroundStyle() render.Style {
	return lp.Foreground.InheritFrom(lp.styleDefaultsForeground())
}

func (lp *LinearProgressChart) styleDefaultsForeground() render.Style {
	return render.Style{
		FillColor:   lp.GetColorPalette().BackgroundColor(),
		StrokeColor: lp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

// GetColorPalette returns the color palette for the chart.
func (lp *LinearProgressChart) GetColorPalette() render.ColorPalette {
	if lp.ColorPalette != nil {
		return lp.ColorPalette
	}
	return render.AlternateColorPalette
}

func (lp *LinearProgressChart) roundedEdgeRadius() float64 {
	return float64(lp.Height()) / 2
}

func (lp *LinearProgressChart) drawBar(r render.Renderer, width int, style render.Style) {
	roundStartRadius := 0.0
	if lp.RoundedEdgeStart {
		roundStartRadius = lp.roundedEdgeRadius()
	}

	roundEndRadius := 0.0
	if lp.RoundedEdgeEnd {
		roundEndRadius = lp.roundedEdgeRadius()
	}

	h := lp.Height()

	r.MoveTo(int(roundStartRadius), 0)
	r.LineTo(width-int(roundEndRadius), 0)

	if lp.RoundedEdgeEnd {
		r.ArcTo(width-int(roundEndRadius), h/2, roundEndRadius, roundEndRadius, mathutil.DegreesToRadians(-90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(width, h)
	}

	r.MoveTo(width-int(roundEndRadius), h)
	r.LineTo(int(roundStartRadius), h)

	if lp.RoundedEdgeStart {
		r.ArcTo(int(roundStartRadius), h/2, roundStartRadius, roundStartRadius, mathutil.DegreesToRadians(90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(0, 0)
	}

	r.SetFillColor(style.FillColor)

	strokeColor := style.StrokeColor
	strokeWidth := style.StrokeWidth

	if style.StrokeWidth == 0 {
		strokeColor = style.FillColor
		strokeWidth = lp.getBackgroundStyle().StrokeWidth
	}

	r.SetStrokeColor(strokeColor)
	r.SetStrokeWidth(strokeWidth)
	r.FillStroke()
}

func (lp *LinearProgressChart) drawBackground(r render.Renderer) {
	lp.drawBar(r, lp.Width(), lp.getBackgroundStyle())
}

func (lp *LinearProgressChart) drawForeground(r render.Renderer) {
	w := float64(lp.Width()) * lp.progress

	lp.drawBar(r, int(w), lp.getForegroundStyle())
}

// Render renders the progrss bar with the given renderer to the given io.Writer.
func (lp *LinearProgressChart) Render(rp render.RendererProvider, w io.Writer) error {
	r, err := rp(lp.Width(), lp.Height())
	if err != nil {
		return err
	}
	r.SetDPI(lp.DPI())

	lp.drawBackground(r)
	lp.drawForeground(r)

	return r.Save(w)
}
