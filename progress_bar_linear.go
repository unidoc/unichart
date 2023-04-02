package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// LinearProgressBar is a component that will render progress bar component.
type LinearProgressBar struct {
	// BackgroundStyle is the style for the background bar.
	BackgroundStyle render.Style

	// ForegroundStyle is the style for the foreground bar.
	ForegroundStyle render.Style

	// ColorPalette is the color pallete that could be used to add colors in this progress bar
	ColorPalette render.ColorPalette

	// RoundedEdgeStart is a flag to enable rounded edge at the start of the bar.
	RoundedEdgeStart bool

	// RoundedEdgeEnd us a flag to enable rounded edge at the end of the bar.
	RoundedEdgeEnd bool

	// CustomTopInfo defines a user provided function to draw a custom info above the progress bar.
	// The callback function should return the height used to render the info, which would be used
	// to shift down the progress bar position.
	CustomTopInfo func(r render.Renderer, x int) int

	// CustomBottomInfo defines a user provided function to draw a custom info under the progress bar.
	CustomBottomInfo func(r render.Renderer, x int)

	height int
	width  int
	dpi    float64
	yPos   int

	// progress is the progress bar values which should be between 0.0 - 1.0.
	progress float64
}

// SetProgress set the progress that will represented by this chart.
// Expected value should be float number between 0.0 - 1.0.
func (lp *LinearProgressBar) SetProgress(progress float64) {
	lp.progress = math.Max(0.0, math.Min(progress, 1.0))
}

// GetProgress returns the progress represented by this chart.
// Returned value should be a float number between 0.0 - 1.0.
func (lp *LinearProgressBar) GetProgress() float64 {
	return lp.progress
}

// DPI returns the DPI of the progress bar.
func (lp *LinearProgressBar) DPI() float64 {
	if lp.dpi == 0 {
		return defaultDPI
	}
	return lp.dpi
}

// SetDPI sets the DPI for the progrss bar.
func (lp *LinearProgressBar) SetDPI(dpi float64) {
	lp.dpi = dpi
}

// Width returns the chart width or the default value.
func (lp *LinearProgressBar) Width() int {
	if lp.width == 0 {
		return defaultChartWidth
	}
	return lp.width
}

// SetWidth sets the chart width.
func (lp *LinearProgressBar) SetWidth(width int) {
	lp.width = width
}

// Height returns the chart height or the default value.
func (lp *LinearProgressBar) Height() int {
	if lp.height == 0 {
		return defaultChartHeight
	}
	return lp.height + lp.yPos
}

// SetHeight sets the chart height.
func (lp *LinearProgressBar) SetHeight(height int) {
	lp.height = height
}

func (lp *LinearProgressBar) getBackgroundStyle() render.Style {
	return lp.BackgroundStyle.InheritFrom(lp.styleDefaultsBackground())
}

func (lp *LinearProgressBar) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   lp.GetColorPalette().BackgroundColor(),
		StrokeColor: lp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (lp *LinearProgressBar) getForegroundStyle() render.Style {
	return lp.ForegroundStyle.InheritFrom(lp.styleDefaultsForeground())
}

func (lp *LinearProgressBar) styleDefaultsForeground() render.Style {
	return render.Style{
		FillColor:   lp.GetColorPalette().BackgroundColor(),
		StrokeColor: lp.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

// GetColorPalette returns the color palette for the chart.
func (lp *LinearProgressBar) GetColorPalette() render.ColorPalette {
	if lp.ColorPalette != nil {
		return lp.ColorPalette
	}
	return render.AlternateColorPalette
}

func (lp *LinearProgressBar) roundedEdgeRadius() float64 {
	return float64(lp.height) / 2
}

func (lp *LinearProgressBar) drawBar(r render.Renderer, width int, style render.Style) {
	roundStartRadius := 0.0
	if lp.RoundedEdgeStart {
		roundStartRadius = lp.roundedEdgeRadius()
	}

	roundEndRadius := 0.0
	if lp.RoundedEdgeEnd {
		roundEndRadius = lp.roundedEdgeRadius()
	}

	h := lp.height

	r.MoveTo(int(roundStartRadius), lp.yPos)
	r.LineTo(width-int(roundEndRadius), lp.yPos)

	if lp.RoundedEdgeEnd {
		r.ArcTo(width-int(roundEndRadius), lp.yPos+h/2, roundEndRadius, roundEndRadius, mathutil.DegreesToRadians(-90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(width, lp.yPos+h)
	}

	r.MoveTo(width-int(roundEndRadius), lp.yPos+h)
	r.LineTo(int(roundStartRadius), lp.yPos+h)

	if lp.RoundedEdgeStart {
		r.ArcTo(int(roundStartRadius), lp.yPos+h/2, roundStartRadius, roundStartRadius, mathutil.DegreesToRadians(90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(0, lp.yPos)
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

func (lp *LinearProgressBar) drawBackground(r render.Renderer) {
	bgStyle := lp.getBackgroundStyle()

	if bgStyle.Hidden {
		return
	}

	lp.drawBar(r, lp.Width(), bgStyle)
}

func (lp *LinearProgressBar) drawForeground(r render.Renderer) {
	fgStyle := lp.getForegroundStyle()

	if fgStyle.Hidden {
		return
	}

	w := float64(lp.Width()) * lp.progress
	lp.drawBar(r, int(w), fgStyle)
}

func (lp *LinearProgressBar) drawTopInfo(r render.Renderer) {
	x := float64(lp.Width()) * lp.progress

	if lp.CustomTopInfo != nil {
		infoHeight := lp.CustomTopInfo(r, int(x))

		lp.yPos = infoHeight

		return
	}
}

func (lp *LinearProgressBar) drawBottomInfo(r render.Renderer) {
	x := float64(lp.Width()) * lp.progress

	if lp.CustomBottomInfo != nil {
		lp.CustomBottomInfo(r, int(x))

		return
	}
}

// Render renders the progrss bar with the given renderer to the given io.Writer.
func (lp *LinearProgressBar) Render(rp render.RendererProvider, w io.Writer) error {
	r, err := rp(lp.Width(), lp.height)
	if err != nil {
		return err
	}
	r.SetDPI(lp.DPI())

	lp.yPos = 0

	lp.drawTopInfo(r)
	lp.drawBackground(r)
	lp.drawForeground(r)
	lp.drawBottomInfo(r)

	return r.Save(w)
}
