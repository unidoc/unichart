package unichart

import (
	"io"
	"math"

	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

type ProgressBar struct {
	Background   render.Style
	Foreground   render.Style
	ColorPalette render.ColorPalette
	Canvas       render.Style

	RoundedEdgeStart bool
	RoundedEdgeEnd   bool

	height int
	width  int
	dpi    float64

	progress float64
}

// SetProgress set the progress that will represented by this chart.
// Expected value should be float number between 0.0 - 1.0.
func (pb *ProgressBar) SetProgress(progress float64) {
	pb.progress = math.Max(0.0, math.Min(progress, 1.0))
}

// GetProgress returns the progress represented by this chart.
// Returned value should be a float number between 0.0 - 1.0.
func (pb *ProgressBar) GetProgress() float64 {
	return pb.progress
}

// DPI returns the DPI of the progress bar.
func (pb *ProgressBar) DPI() float64 {
	if pb.dpi == 0 {
		return defaultDPI
	}
	return pb.dpi
}

// SetDPI sets the DPI for the progrss bar.
func (pb *ProgressBar) SetDPI(dpi float64) {
	pb.dpi = dpi
}

// Width returns the chart width or the default value.
func (pb *ProgressBar) Width() int {
	if pb.width == 0 {
		return defaultChartWidth
	}
	return pb.width
}

// SetWidth sets the chart width.
func (pb *ProgressBar) SetWidth(width int) {
	pb.width = width
}

// Height returns the chart height or the default value.
func (pb *ProgressBar) Height() int {
	if pb.height == 0 {
		return defaultChartHeight
	}
	return pb.height
}

// SetHeight sets the chart height.
func (pb *ProgressBar) SetHeight(height int) {
	pb.height = height
}

func (pb *ProgressBar) getBackgroundStyle() render.Style {
	return pb.Background.InheritFrom(pb.styleDefaultsBackground())
}

func (pb *ProgressBar) styleDefaultsBackground() render.Style {
	return render.Style{
		FillColor:   pb.GetColorPalette().BackgroundColor(),
		StrokeColor: pb.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

func (pb *ProgressBar) getForegroundStyle() render.Style {
	return pb.Foreground.InheritFrom(pb.styleDefaultsForeground())
}

func (pb *ProgressBar) styleDefaultsForeground() render.Style {
	return render.Style{
		FillColor:   pb.GetColorPalette().BackgroundColor(),
		StrokeColor: pb.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: render.DefaultStrokeWidth,
	}
}

// GetColorPalette returns the color palette for the chart.
func (pb *ProgressBar) GetColorPalette() render.ColorPalette {
	if pb.ColorPalette != nil {
		return pb.ColorPalette
	}
	return render.AlternateColorPalette
}

func (pb *ProgressBar) roundedEdgeRadius() float64 {
	return float64(pb.Height()) / 2
}

func (pb *ProgressBar) drawBar(r render.Renderer, width int, style render.Style) {
	roundStartRadius := 0.0
	if pb.RoundedEdgeStart {
		roundStartRadius = pb.roundedEdgeRadius()
	}

	roundEndRadius := 0.0
	if pb.RoundedEdgeEnd {
		roundEndRadius = pb.roundedEdgeRadius()
	}

	h := pb.Height()

	r.MoveTo(int(roundStartRadius), 0)
	r.LineTo(width-int(roundEndRadius), 0)

	if pb.RoundedEdgeEnd {
		r.ArcTo(width-int(roundEndRadius), h/2, roundEndRadius, roundEndRadius, mathutil.DegreesToRadians(-90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(width, h)
	}

	r.MoveTo(width-int(roundEndRadius), h)
	r.LineTo(int(roundStartRadius), h)

	if pb.RoundedEdgeStart {
		r.ArcTo(int(roundStartRadius), h/2, roundStartRadius, roundStartRadius, mathutil.DegreesToRadians(90), mathutil.DegreesToRadians(180))
	} else {
		r.LineTo(0, 0)
	}

	r.SetFillColor(style.FillColor)

	strokeColor := style.StrokeColor
	strokeWidth := style.StrokeWidth

	if style.StrokeWidth == 0 {
		strokeColor = style.FillColor
		strokeWidth = pb.getBackgroundStyle().StrokeWidth
	}

	r.SetStrokeColor(strokeColor)
	r.SetStrokeWidth(strokeWidth)
	r.FillStroke()
}

func (pb *ProgressBar) drawBackground(r render.Renderer) {
	pb.drawBar(r, pb.Width(), pb.getBackgroundStyle())
}

func (pb *ProgressBar) drawForeground(r render.Renderer) {
	w := float64(pb.Width()) * pb.progress

	pb.drawBar(r, int(w), pb.getForegroundStyle())
}

// Render renders the progrss bar with the given renderer to the given io.Writer.
func (pb *ProgressBar) Render(rp render.RendererProvider, w io.Writer) error {
	r, err := rp(pb.Width(), pb.Height())
	if err != nil {
		return err
	}
	r.SetDPI(pb.DPI())

	pb.drawBackground(r)
	pb.drawForeground(r)

	return r.Save(w)
}
