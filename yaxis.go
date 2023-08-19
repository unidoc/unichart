package unichart

import (
	"math"

	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// YAxis is a veritcal rule of the range.
// There can be (2) y-axes; a primary and secondary.
type YAxis struct {
	Name      string
	NameStyle render.Style
	Style     render.Style

	Zero      GridLine
	AxisType  dataset.YAxisType
	Ascending bool

	ValueFormatter dataset.ValueFormatter
	Range          sequence.Range

	TickStyle render.Style
	Ticks     []Tick

	GridLines      []GridLine
	GridMajorStyle render.Style
	GridMinorStyle render.Style
}

// GetName returns the name.
func (ya YAxis) GetName() string {
	return ya.Name
}

// GetNameStyle returns the name style.
func (ya YAxis) GetNameStyle() render.Style {
	return ya.NameStyle
}

// GetStyle returns the style.
func (ya YAxis) GetStyle() render.Style {
	return ya.Style
}

// GetValueFormatter returns the value formatter for the axis.
func (ya YAxis) GetValueFormatter() dataset.ValueFormatter {
	if ya.ValueFormatter != nil {
		return ya.ValueFormatter
	}
	return dataset.FloatValueFormatter
}

// GetTickStyle returns the tick style.
func (ya YAxis) GetTickStyle() render.Style {
	return ya.TickStyle
}

// GetTicks returns the ticks for a series.
// The coalesce priority is:
// 	- User Supplied Ticks (i.e. Ticks array on the axis itself).
// 	- Range ticks (i.e. if the range provides ticks).
//	- Generating continuous ticks based on minimum spacing and canvas width.
func (ya YAxis) GetTicks(r render.Renderer, ra sequence.Range, defaults render.Style, vf dataset.ValueFormatter) []Tick {
	if len(ya.Ticks) > 0 {
		return ya.Ticks
	}
	if tp, isTickProvider := ra.(TicksProvider); isTickProvider {
		return tp.GetTicks(r, defaults, vf)
	}

	tickStyle := ya.Style.InheritFrom(defaults)
	return generateContinuousTicks(r, ra, true, tickStyle, vf)
}

// GetGridLines returns the gridlines for the axis.
func (ya YAxis) GetGridLines(ticks []Tick) []GridLine {
	if len(ya.GridLines) > 0 {
		return ya.GridLines
	}
	return GenerateGridLines(ticks, ya.GridMajorStyle, ya.GridMinorStyle)
}

// Measure returns the bounds of the axis.
func (ya YAxis) Measure(r render.Renderer, canvasBox render.Box, ra sequence.Range, defaults render.Style, ticks []Tick) render.Box {
	var tx int
	if ya.AxisType == dataset.YAxisPrimary {
		tx = canvasBox.Right + defaultYAxisMargin
	} else if ya.AxisType == dataset.YAxisSecondary {
		tx = canvasBox.Left - defaultYAxisMargin
	}

	ya.TickStyle.InheritFrom(ya.Style.InheritFrom(defaults)).WriteToRenderer(r)
	var minx, maxx, miny, maxy = math.MaxInt32, 0, math.MaxInt32, 0
	var maxTextHeight int
	for _, t := range ticks {
		v := t.Value
		ly := canvasBox.Bottom - ra.Translate(v)

		tb := r.MeasureText(t.Label)
		tbh2 := tb.Height() >> 1
		finalTextX := tx
		if ya.AxisType == dataset.YAxisSecondary {
			finalTextX = tx - tb.Width()
		}

		maxTextHeight = mathutil.MaxInt(tb.Height(), maxTextHeight)

		if ya.AxisType == dataset.YAxisPrimary {
			minx = canvasBox.Right
			maxx = mathutil.MaxInt(maxx, tx+tb.Width())
		} else if ya.AxisType == dataset.YAxisSecondary {
			minx = mathutil.MinInt(minx, finalTextX)
			maxx = mathutil.MaxInt(maxx, tx)
		}

		miny = mathutil.MinInt(miny, ly-tbh2)
		maxy = mathutil.MaxInt(maxy, ly+tbh2)
	}

	if !ya.NameStyle.Hidden && len(ya.Name) > 0 {
		maxx += (defaultYAxisMargin + maxTextHeight)
	}

	return render.Box{
		Top:    miny,
		Left:   minx,
		Right:  maxx,
		Bottom: maxy,
	}
}

// Render renders the axis.
func (ya YAxis) Render(r render.Renderer, canvasBox render.Box, ra sequence.Range, defaults render.Style, ticks []Tick) {
	tickStyle := ya.TickStyle.InheritFrom(ya.Style.InheritFrom(defaults))
	tickStyle.WriteToRenderer(r)

	sw := tickStyle.GetStrokeWidth(defaults.StrokeWidth)

	var lx int
	var tx int
	if ya.AxisType == dataset.YAxisPrimary {
		lx = canvasBox.Right
		tx = lx + defaultYAxisMargin
	} else if ya.AxisType == dataset.YAxisSecondary {
		lx = canvasBox.Left
		tx = lx - defaultYAxisMargin
	}

	var maxTextWidth int
	var finalTextX, finalTextY int
	for _, t := range ticks {
		v := t.Value
		ly := canvasBox.Bottom - ra.Translate(v)

		tb := render.Text.Measure(r, t.Label, tickStyle)

		if tb.Width() > maxTextWidth {
			maxTextWidth = tb.Width()
		}

		if ya.AxisType == dataset.YAxisSecondary {
			finalTextX = tx - tb.Width()
		} else {
			finalTextX = tx
		}

		if tickStyle.TextRotationDegrees == 0 {
			finalTextY = ly + tb.Height()>>1
		} else {
			finalTextY = ly
		}

		tickStyle.WriteToRenderer(r)

		if lx < 0 || ly < 0 {
			continue
		}

		r.MoveTo(lx, ly)
		if ya.AxisType == dataset.YAxisPrimary {
			r.LineTo(lx+defaultHorizontalTickWidth, ly)
		} else if ya.AxisType == dataset.YAxisSecondary {
			r.LineTo(lx-defaultHorizontalTickWidth, ly)
		}
		r.Stroke()

		render.Text.Draw(r, t.Label, finalTextX, finalTextY, tickStyle)
	}

	nameStyle := ya.NameStyle.InheritFrom(defaults.InheritFrom(render.Style{TextRotationDegrees: 90}))
	if !ya.NameStyle.Hidden && len(ya.Name) > 0 {
		nameStyle.GetTextOptions().WriteToRenderer(r)
		tb := render.Text.Measure(r, ya.Name, nameStyle)

		var tx int
		if ya.AxisType == dataset.YAxisPrimary {
			tx = canvasBox.Right + int(sw) + defaultYAxisMargin + maxTextWidth + defaultYAxisMargin
		} else if ya.AxisType == dataset.YAxisSecondary {
			tx = canvasBox.Left - (defaultYAxisMargin + int(sw) + maxTextWidth + defaultYAxisMargin)
		}

		var ty int
		if nameStyle.TextRotationDegrees == 0 {
			ty = canvasBox.Top + (canvasBox.Height()>>1 - tb.Width()>>1)
		} else {
			ty = canvasBox.Top + (canvasBox.Height()>>1 - tb.Height()>>1)
		}

		render.Text.Draw(r, ya.Name, tx, ty, nameStyle)
	}

	if !ya.Zero.Style.Hidden {
		ya.Zero.Render(r, canvasBox, ra, false, render.Style{})
	}

	if !ya.GridMajorStyle.Hidden || !ya.GridMinorStyle.Hidden {
		for _, gl := range ya.GetGridLines(ticks) {
			if (gl.IsMinor && !ya.GridMinorStyle.Hidden) || (!gl.IsMinor && !ya.GridMajorStyle.Hidden) {
				defaults := ya.GridMajorStyle
				if gl.IsMinor {
					defaults = ya.GridMinorStyle
				}
				gl.Render(r, canvasBox, ra, false, gl.Style.InheritFrom(defaults))
			}
		}
	}
}

func (ya YAxis) RenderAxisLine(r render.Renderer, canvasBox render.Box, ra sequence.Range, defaults render.Style, ticks []Tick) {
	tickStyle := ya.TickStyle.InheritFrom(ya.Style.InheritFrom(defaults))
	tickStyle.WriteToRenderer(r)

	sw := tickStyle.GetStrokeWidth(defaults.StrokeWidth)

	var lx int
	if ya.AxisType == dataset.YAxisPrimary {
		lx = canvasBox.Right
	} else if ya.AxisType == dataset.YAxisSecondary {
		lx = canvasBox.Left
	}

	r.SetStrokeWidth(sw)
	r.MoveTo(lx, canvasBox.Bottom)
	r.LineTo(lx, canvasBox.Top)
	r.Stroke()
}
