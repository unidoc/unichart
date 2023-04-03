package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
	"github.com/unidoc/unipdf/v3/model"
)

// LinearProgressBar is a component that will render progress bar component.
type LinearProgressBar struct {
	// BackgroundStyle is the style for the background bar.
	BackgroundStyle render.Style

	// ForegroundStyle is the style for the foreground bar.
	ForegroundStyle render.Style

	// LabelStyle is the style for the label that will displayed on the progress bar.
	LabelStyle render.Style

	// ColorPalette is the color pallete that could be used to add colors in this progress bar
	ColorPalette render.ColorPalette

	// RoundedEdgeStart is a flag to enable rounded edge at the start of the bar.
	RoundedEdgeStart bool

	// RoundedEdgeEnd us a flag to enable rounded edge at the end of the bar.
	RoundedEdgeEnd bool

	// CustomTopInfo defines a user provided function to draw a custom info above the progress bar.
	// The callback function should return the height occupied by the top info, which would be used
	// to shift down the progress bar position.
	CustomTopInfo func(r render.Renderer, x int) int

	// CustomBottomInfo defines a user provided function to draw a custom info under the progress bar.
	// The callback function should return the height occupied by the bottom info, which would be used
	// to calculate this progress bar total height.
	CustomBottomInfo func(r render.Renderer, x int) int

	height           int
	width            int
	dpi              float64
	barYPos          int
	bottomInfoHeight int

	// progress is the progress bar values which should be between 0.0 - 1.0.
	progress float64

	// label is the text that will be displayed on the progress bar.
	label string
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
	return lp.height + lp.barYPos + lp.bottomInfoHeight
}

// SetHeight sets the chart height.
func (lp *LinearProgressBar) SetHeight(height int) {
	lp.height = height
}

// SetLabel sets the labe that would be rendered on the progress bar.
func (lp *LinearProgressBar) SetLabel(label string) {
	lp.label = label
}

// GetLabel would returns label that rendered on the progress bar.
func (lp *LinearProgressBar) GetLabel() string {
	return lp.label
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

func (lp *LinearProgressBar) getLabelStyle() render.Style {
	return lp.LabelStyle.InheritFrom(lp.styleDefaultsLabel())
}

func (lp *LinearProgressBar) styleDefaultsLabel() render.Style {
	return render.Style{
		Font:                model.DefaultFont(),
		FontSize:            render.DefaultFontSize,
		FontColor:           lp.getForegroundStyle().FillColor,
		TextHorizontalAlign: render.TextHorizontalAlignLeft,
		TextVerticalAlign:   render.TextVerticalAlignMiddle,
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

	r.MoveTo(int(roundStartRadius), lp.barYPos)
	r.LineTo(width-int(roundEndRadius), lp.barYPos)

	if lp.RoundedEdgeEnd {
		r.ArcTo(width-int(roundEndRadius), lp.barYPos+h/2, roundEndRadius, roundEndRadius, mathutil.DegreesToRadians(-90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(width, lp.barYPos+h)
	}

	r.MoveTo(width-int(roundEndRadius), lp.barYPos+h)
	r.LineTo(int(roundStartRadius), lp.barYPos+h)

	if lp.RoundedEdgeStart {
		r.ArcTo(int(roundStartRadius), lp.barYPos+h/2, roundStartRadius, roundStartRadius, mathutil.DegreesToRadians(90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(0, lp.barYPos)
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

func (lp *LinearProgressBar) drawLabel(r render.Renderer) {
	labelStyle := lp.getLabelStyle()

	if labelStyle.Hidden || lp.label == "" {
		return
	}

	x := float64(lp.Width()) * lp.progress
	y := lp.barYPos

	render.Text.DrawWithin(r, lp.label, render.NewBox(y, int(x)+10, lp.Width(), y+lp.height), labelStyle)
}

func (lp *LinearProgressBar) drawTopInfo(r render.Renderer) {
	x := float64(lp.Width()) * lp.progress

	if lp.CustomTopInfo != nil {
		lp.barYPos = lp.CustomTopInfo(r, int(x))

		return
	}
}

func (lp *LinearProgressBar) drawBottomInfo(r render.Renderer) {
	x := float64(lp.Width()) * lp.progress

	if lp.CustomBottomInfo != nil {
		lp.bottomInfoHeight = lp.CustomBottomInfo(r, int(x))

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

	lp.barYPos = 0

	lp.drawTopInfo(r)
	lp.drawBackground(r)
	lp.drawForeground(r)
	lp.drawLabel(r)
	lp.drawBottomInfo(r)

	return r.Save(w)
}
