package unichart

import (
	"math"

	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// XAxis represents the horizontal axis.
type XAxis struct {
	Name      string
	NameStyle render.Style

	Style          render.Style
	ValueFormatter dataset.ValueFormatter
	Range          sequence.Range

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
func (xa XAxis) GetValueFormatter() dataset.ValueFormatter {
	if xa.ValueFormatter != nil {
		return xa.ValueFormatter
	}
	return dataset.FloatValueFormatter
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
func (xa XAxis) GetTicks(r render.Renderer, ra sequence.Range, defaults render.Style, vf dataset.ValueFormatter) []Tick {
	if len(xa.Ticks) > 0 {
		return xa.Ticks
	}
	if tp, isTickProvider := ra.(TicksProvider); isTickProvider {
		return tp.GetTicks(r, defaults, vf)
	}

	tickStyle := xa.Style.InheritFrom(defaults)
	return generateContinuousTicks(r, ra, false, tickStyle, vf)
}

// GetGridLines returns the gridlines for the axis.
func (xa XAxis) GetGridLines(ticks []Tick) []GridLine {
	if len(xa.GridLines) > 0 {
		return xa.GridLines
	}
	return GenerateGridLines(ticks, xa.GridMajorStyle, xa.GridMinorStyle)
}

// Measure returns the bounds of the axis.
func (xa XAxis) Measure(r render.Renderer, canvasBox render.Box, ra sequence.Range, defaults render.Style, ticks []Tick) render.Box {
	tickStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))

	tp := xa.GetTickPosition()

	var ltx, rtx int
	var tx int
	var left, right, bottom = math.MaxInt32, 0, 0
	for index, t := range ticks {
		v := t.Value

		tx = canvasBox.Left + ra.Translate(v)
		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			tb := render.Text.Measure(r, t.Label, tickStyle.GetTextOptions())
			ltx = tx - tb.Width()>>1
			rtx = tx + tb.Width()>>1
			bottom = mathutil.MaxInt(bottom, tb.Height())
		case TickPositionBetweenTicks:
			if index > 0 {
				ltx = ra.Translate(ticks[index-1].Value)
				rtx = tx

				finalTickStyle := tickStyle.InheritFrom(render.Style{TextHorizontalAlign: render.TextHorizontalAlignCenter})
				ftb := render.Text.MeasureLines(r, render.Text.WrapFit(r, t.Label, tx-ltx, finalTickStyle), finalTickStyle)
				bottom = mathutil.MaxInt(bottom, ftb.Height())
			}
		}

		left = mathutil.MinInt(left, ltx)
		right = mathutil.MaxInt(right, rtx)
	}

	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := render.Text.Measure(r, xa.Name, xa.NameStyle.InheritFrom(defaults))
		bottom += defaultXAxisMargin + tb.Height()
	}

	return render.Box{
		Top:    canvasBox.Bottom,
		Left:   left,
		Right:  right,
		Bottom: canvasBox.Bottom + defaultXAxisMargin + bottom,
	}
}

// Render renders the axis
func (xa XAxis) Render(r render.Renderer, canvasBox render.Box, ra sequence.Range, defaults render.Style, ticks []Tick) {
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
		r.LineTo(tx, canvasBox.Bottom+defaultVerticalTickHeight)
		r.Stroke()

		tickWithAxisStyle := xa.TickStyle.InheritFrom(xa.Style.InheritFrom(defaults))
		tb := render.Text.Measure(r, t.Label, tickWithAxisStyle)

		switch tp {
		case TickPositionUnderTick, TickPositionUnset:
			if tickStyle.TextRotationDegrees == 0 {
				tx = tx - tb.Width()>>1
				ty = canvasBox.Bottom + defaultXAxisMargin + tb.Height()
			} else {
				ty = canvasBox.Bottom + (2 * defaultXAxisMargin)
			}
			render.Text.Draw(r, t.Label, tx, ty, tickWithAxisStyle)
			maxTextHeight = mathutil.MaxInt(maxTextHeight, tb.Height())
		case TickPositionBetweenTicks:
			if index > 0 {
				llx := ra.Translate(ticks[index-1].Value)
				ltx := canvasBox.Left + llx
				finalTickStyle := tickWithAxisStyle.InheritFrom(render.Style{TextHorizontalAlign: render.TextHorizontalAlignCenter})

				render.Text.DrawWithin(r, t.Label, render.Box{
					Left:   ltx,
					Right:  tx,
					Top:    canvasBox.Bottom + defaultXAxisMargin,
					Bottom: canvasBox.Bottom + defaultXAxisMargin,
				}, finalTickStyle)

				ftb := render.Text.MeasureLines(r, render.Text.WrapFit(r, t.Label, tx-ltx, finalTickStyle), finalTickStyle)
				maxTextHeight = mathutil.MaxInt(maxTextHeight, ftb.Height())
			}
		}
	}

	nameStyle := xa.NameStyle.InheritFrom(defaults)
	if !xa.NameStyle.Hidden && len(xa.Name) > 0 {
		tb := render.Text.Measure(r, xa.Name, nameStyle)
		tx := canvasBox.Right - (canvasBox.Width()>>1 + tb.Width()>>1)
		ty := canvasBox.Bottom + defaultXAxisMargin + maxTextHeight + defaultXAxisMargin + tb.Height()
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
