package chart

import (
	"math"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// XAxis represents the horizontal axis.
type XAxis struct {
	Name      string
	NameStyle render.Style

	Style          render.Style
	ValueFormatter data.ValueFormatter
	Range          data.Range

	TickStyle    render.Style
	Ticks        []Tick
	TickPosition TickPosition

	GridLines      []GridLine
	GridMajorStyle render.Style
	GridMinorStyle render.Style
}

// GetName returns the name.
func (xa XAxis) GetName() string {
	return xa.Name
}

// GetStyle returns the style.
func (xa XAxis) GetStyle() render.Style {
	return xa.Style
}

// GetValueFormatter returns the value formatter for the axis.
func (xa XAxis) GetValueFormatter() data.ValueFormatter {
	if xa.ValueFormatter != nil {
		return xa.ValueFormatter
	}
	return data.FloatValueFormatter
}

// GetTickPosition returns the tick position option for the axis.
func (xa XAxis) GetTickPosition(defaults ...TickPosition) TickPosition {
	if xa.TickPosition == TickPositionUnset {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return TickPositionUnderTick
	}
	return xa.TickPosition
}

// GetTicks returns the ticks for a series.
// The coalesce priority is:
// 	- User Supplied Ticks (i.e. Ticks array on the axis itself).
// 	- Range ticks (i.e. if the range provides ticks).
//	- Generating continuous ticks based on minimum spacing and canvas width.
func (xa XAxis) GetTicks(r render.Renderer, ra data.Range, defaults render.Style, vf data.ValueFormatter) []Tick {
	if len(xa.Ticks) > 0 {
		return xa.Ticks
	}
	if tp, isTickProvider := ra.(TicksProvider); isTickProvider {
		return tp.GetTicks(r, defaults, vf)
	}
	tickStyle := xa.Style.InheritFrom(defaults)
	return GenerateContinuousTicks(r, ra, false, tickStyle, vf)
}

// GetGridLines returns the gridlines for the axis.
func (xa XAxis) GetGridLines(ticks []Tick) []GridLine {
	if len(xa.GridLines) > 0 {
		return xa.GridLines
	}
	return GenerateGridLines(ticks, xa.GridMajorStyle, xa.GridMinorStyle)
}

// Measure returns the bounds of the axis.
func (xa XAxis) Measure(r render.Renderer, canvasBox render.Box, ra data.Range, defaults render.Style, ticks []Tick) render.Box {
	tickStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))

	tp := xa.GetTickPosition()

	var ltx, rtx int
	var tx, ty int
	var left, right, bottom = math.MaxInt32, 0, 0
	for index, t := range ticks {
		v := t.Value
		tb := render.Text.Measure(r, t.Label, tickStyle.GetTextOptions())

		tx = canvasBox.Left + ra.Translate(v)
		ty = canvasBox.Bottom + DefaultXAxisMargin + tb.Height()
		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			ltx = tx - tb.Width()>>1
			rtx = tx + tb.Width()>>1
			break
		case TickPositionBetweenTicks:
			if index > 0 {
				ltx = ra.Translate(ticks[index-1].Value)
				rtx = tx
			}
			break
		}

		left = mathutil.MinInt(left, ltx)
		right = mathutil.MaxInt(right, rtx)
		bottom = mathutil.MaxInt(bottom, ty)
	}

	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := render.Text.Measure(r, xa.Name, xa.NameStyle.InheritFrom(defaults))
		bottom += DefaultXAxisMargin + tb.Height()
	}

	return render.Box{
		Top:    canvasBox.Bottom,
		Left:   left,
		Right:  right,
		Bottom: bottom,
	}
}

// Render renders the axis
func (xa XAxis) Render(r render.Renderer, canvasBox render.Box, ra data.Range, defaults render.Style, ticks []Tick) {
	tickStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))

	tickStyle.GetStrokeOptions().WriteToRenderer(r)
	r.MoveTo(canvasBox.Left, canvasBox.Bottom)
	r.LineTo(canvasBox.Right, canvasBox.Bottom)
	r.Stroke()

	tp := xa.GetTickPosition()

	var tx, ty int
	var maxTextHeight int
	for index, t := range ticks {
		v := t.Value
		lx := ra.Translate(v)

		tx = canvasBox.Left + lx

		tickStyle.GetStrokeOptions().WriteToRenderer(r)
		r.MoveTo(tx, canvasBox.Bottom)
		r.LineTo(tx, canvasBox.Bottom+DefaultVerticalTickHeight)
		r.Stroke()

		tickWithAxisStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))
		tb := render.Text.Measure(r, t.Label, tickWithAxisStyle)

		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			if tickStyle.TextRotationDegrees == 0 {
				tx = tx - tb.Width()>>1
				ty = canvasBox.Bottom + DefaultXAxisMargin + tb.Height()
			} else {
				ty = canvasBox.Bottom + (2 * DefaultXAxisMargin)
			}
			render.Text.Draw(r, t.Label, tx, ty, tickWithAxisStyle)
			maxTextHeight = mathutil.MaxInt(maxTextHeight, tb.Height())
			break
		case TickPositionBetweenTicks:
			if index > 0 {
				llx := ra.Translate(ticks[index-1].Value)
				ltx := canvasBox.Left + llx
				finalTickStyle := tickWithAxisStyle.InheritFrom(render.Style{TextHorizontalAlign: render.TextHorizontalAlignCenter})

				render.Text.DrawWithin(r, t.Label, render.Box{
					Left:   ltx,
					Right:  tx,
					Top:    canvasBox.Bottom + DefaultXAxisMargin,
					Bottom: canvasBox.Bottom + DefaultXAxisMargin,
				}, finalTickStyle)

				ftb := render.Text.MeasureLines(r, render.Text.WrapFit(r, t.Label, tx-ltx, finalTickStyle), finalTickStyle)
				maxTextHeight = mathutil.MaxInt(maxTextHeight, ftb.Height())
			}
			break
		}
	}

	nameStyle := xa.NameStyle.InheritFrom(defaults)
	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := render.Text.Measure(r, xa.Name, nameStyle)
		tx := canvasBox.Right - (canvasBox.Width()>>1 + tb.Width()>>1)
		ty := canvasBox.Bottom + DefaultXAxisMargin + maxTextHeight + DefaultXAxisMargin + tb.Height()
		render.Text.Draw(r, xa.Name, tx, ty, nameStyle)
	}

	if !xa.GridMajorStyle.Hidden || !xa.GridMinorStyle.Hidden {
		for _, gl := range xa.GetGridLines(ticks) {
			if (gl.IsMinor && !xa.GridMinorStyle.Hidden) || (!gl.IsMinor && !xa.GridMajorStyle.Hidden) {
				defaults := xa.GridMajorStyle
				if gl.IsMinor {
					defaults = xa.GridMinorStyle
				}
				gl.Render(r, canvasBox, ra, true, gl.Style.InheritFrom(defaults))
			}
		}
	}
}
