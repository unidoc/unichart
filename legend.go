package unichart

import (
	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// Legend returns a legend renderable function.
func Legend(c *Chart, userDefaults ...render.Style) render.Renderable {
	return func(r render.Renderer, cb render.Box, chartDefaults render.Style) {
		legendDefaults := render.Style{
			FillColor:   render.ColorWhite,
			FontColor:   render.DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: render.DefaultLineColor,
			StrokeWidth: defaultAxisLineWidth,
		}

		var legendStyle render.Style
		if len(userDefaults) > 0 {
			legendStyle = userDefaults[0].InheritFrom(chartDefaults.InheritFrom(legendDefaults))
		} else {
			legendStyle = chartDefaults.InheritFrom(legendDefaults)
		}

		// DEFAULTS
		legendPadding := render.Box{
			Top:    5,
			Left:   5,
			Right:  5,
			Bottom: 5,
		}
		lineTextGap := 5
		lineLengthMinimum := 25

		var labels []string
		var lines []render.Style
		for index, s := range c.Series {
			if !s.GetStyle().Hidden {
				if _, isAnnotationSeries := s.(dataset.AnnotationSeries); !isAnnotationSeries {
					labels = append(labels, s.GetName())
					lines = append(lines, s.GetStyle().InheritFrom(c.styleDefaultsSeries(index)))
				}
			}
		}

		legend := render.Box{
			Top:  cb.Top,
			Left: cb.Left,
			// bottom and right will be sized by the legend content + relevant padding.
		}

		legendContent := render.Box{
			Top:    legend.Top + legendPadding.Top,
			Left:   legend.Left + legendPadding.Left,
			Right:  legend.Left + legendPadding.Left,
			Bottom: legend.Top + legendPadding.Top,
		}

		legendStyle.GetTextOptions().WriteToRenderer(r)

		// measure
		labelCount := 0
		for x := 0; x < len(labels); x++ {
			if len(labels[x]) > 0 {
				tb := r.MeasureText(labels[x])
				if labelCount > 0 {
					legendContent.Bottom += defaultMinimumTickVerticalSpacing
				}
				legendContent.Bottom += tb.Height()
				right := legendContent.Left + tb.Width() + lineTextGap + lineLengthMinimum
				legendContent.Right = mathutil.MaxInt(legendContent.Right, right)
				labelCount++
			}
		}

		legend = legend.Grow(legendContent)
		legend.Right = legendContent.Right + legendPadding.Right
		legend.Bottom = legendContent.Bottom + legendPadding.Bottom
		legend.Draw(r, legendStyle)

		legendStyle.GetTextOptions().WriteToRenderer(r)

		ycursor := legendContent.Top
		tx := legendContent.Left
		legendCount := 0
		var label string
		for x := 0; x < len(labels); x++ {
			label = labels[x]
			if len(label) > 0 {
				if legendCount > 0 {
					ycursor += defaultMinimumTickVerticalSpacing
				}

				tb := r.MeasureText(label)

				ty := ycursor + tb.Height()
				r.Text(label, tx, ty)

				th2 := tb.Height() >> 1

				lx := tx + tb.Width() + lineTextGap
				ly := ty - th2
				lx2 := legendContent.Right - legendPadding.Right

				r.SetStrokeColor(lines[x].GetStrokeColor())
				r.SetStrokeWidth(lines[x].GetStrokeWidth())
				r.SetStrokeDashArray(lines[x].GetStrokeDashArray())

				r.MoveTo(lx, ly)
				r.LineTo(lx2, ly)
				r.Stroke()

				ycursor += tb.Height()
				legendCount++
			}
		}
	}
}

// LegendThin is a legend that doesn't obscure the chart area.
func LegendThin(c *Chart, userDefaults ...render.Style) render.Renderable {
	return func(r render.Renderer, cb render.Box, chartDefaults render.Style) {
		legendDefaults := render.Style{
			FillColor:   render.ColorWhite,
			FontColor:   render.DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: render.DefaultLineColor,
			StrokeWidth: defaultAxisLineWidth,
			Padding: render.Box{
				Top:    5,
				Left:   7,
				Right:  7,
				Bottom: 5,
			},
		}

		var legendStyle render.Style
		if len(userDefaults) > 0 {
			legendStyle = userDefaults[0].InheritFrom(chartDefaults.InheritFrom(legendDefaults))
		} else {
			legendStyle = chartDefaults.InheritFrom(legendDefaults)
		}

		r.SetFont(legendStyle.GetFont())
		r.SetFontColor(legendStyle.GetFontColor())
		r.SetFontSize(legendStyle.GetFontSize())

		var labels []string
		var lines []render.Style
		for index, s := range c.Series {
			if !s.GetStyle().Hidden {
				if _, isAnnotationSeries := s.(dataset.AnnotationSeries); !isAnnotationSeries {
					labels = append(labels, s.GetName())
					lines = append(lines, s.GetStyle().InheritFrom(c.styleDefaultsSeries(index)))
				}
			}
		}

		var textHeight int
		var textWidth int
		var textBox render.Box
		for x := 0; x < len(labels); x++ {
			if len(labels[x]) > 0 {
				textBox = r.MeasureText(labels[x])
				textHeight = mathutil.MaxInt(textBox.Height(), textHeight)
				textWidth = mathutil.MaxInt(textBox.Width(), textWidth)
			}
		}

		legendBoxHeight := textHeight + legendStyle.Padding.Top + legendStyle.Padding.Bottom
		chartPadding := cb.Top
		legendYMargin := (chartPadding - legendBoxHeight) >> 1

		legendBox := render.Box{
			Left:   cb.Left,
			Right:  cb.Right,
			Top:    legendYMargin,
			Bottom: legendYMargin + legendBoxHeight,
		}
		legendBox.Draw(r, legendDefaults)

		r.SetFont(legendStyle.GetFont())
		r.SetFontColor(legendStyle.GetFontColor())
		r.SetFontSize(legendStyle.GetFontSize())

		lineTextGap := 5
		lineLengthMinimum := 25

		tx := legendBox.Left + legendStyle.Padding.Left
		ty := legendYMargin + legendStyle.Padding.Top + textHeight
		var label string
		var lx, ly int
		th2 := textHeight >> 1
		for index := range labels {
			label = labels[index]
			if len(label) > 0 {
				textBox = r.MeasureText(label)
				r.Text(label, tx, ty)

				lx = tx + textBox.Width() + lineTextGap
				ly = ty - th2

				r.SetStrokeColor(lines[index].GetStrokeColor())
				r.SetStrokeWidth(lines[index].GetStrokeWidth())
				r.SetStrokeDashArray(lines[index].GetStrokeDashArray())

				r.MoveTo(lx, ly)
				r.LineTo(lx+lineLengthMinimum, ly)
				r.Stroke()

				tx += textBox.Width() + defaultMinimumTickHorizontalSpacing + lineTextGap + lineLengthMinimum
			}
		}
	}
}

// LegendLeft is a legend that is designed for longer series lists.
func LegendLeft(c *Chart, userDefaults ...render.Style) render.Renderable {
	return func(r render.Renderer, cb render.Box, chartDefaults render.Style) {
		legendDefaults := render.Style{
			FillColor:   render.ColorWhite,
			FontColor:   render.DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: render.DefaultLineColor,
			StrokeWidth: defaultAxisLineWidth,
		}

		var legendStyle render.Style
		if len(userDefaults) > 0 {
			legendStyle = userDefaults[0].InheritFrom(chartDefaults.InheritFrom(legendDefaults))
		} else {
			legendStyle = chartDefaults.InheritFrom(legendDefaults)
		}

		// DEFAULTS
		legendPadding := render.Box{
			Top:    5,
			Left:   5,
			Right:  5,
			Bottom: 5,
		}
		lineTextGap := 5
		lineLengthMinimum := 25

		var labels []string
		var lines []render.Style
		for index, s := range c.Series {
			if !s.GetStyle().Hidden {
				if _, isAnnotationSeries := s.(dataset.AnnotationSeries); !isAnnotationSeries {
					labels = append(labels, s.GetName())
					lines = append(lines, s.GetStyle().InheritFrom(c.styleDefaultsSeries(index)))
				}
			}
		}

		legend := render.Box{
			Top:  5,
			Left: 5,
			// bottom and right will be sized by the legend content + relevant padding.
		}

		legendContent := render.Box{
			Top:    legend.Top + legendPadding.Top,
			Left:   legend.Left + legendPadding.Left,
			Right:  legend.Left + legendPadding.Left,
			Bottom: legend.Top + legendPadding.Top,
		}

		legendStyle.GetTextOptions().WriteToRenderer(r)

		// measure
		labelCount := 0
		for x := 0; x < len(labels); x++ {
			if len(labels[x]) > 0 {
				tb := r.MeasureText(labels[x])
				if labelCount > 0 {
					legendContent.Bottom += defaultMinimumTickVerticalSpacing
				}
				legendContent.Bottom += tb.Height()
				right := legendContent.Left + tb.Width() + lineTextGap + lineLengthMinimum
				legendContent.Right = mathutil.MaxInt(legendContent.Right, right)
				labelCount++
			}
		}

		legend = legend.Grow(legendContent)
		legend.Right = legendContent.Right + legendPadding.Right
		legend.Bottom = legendContent.Bottom + legendPadding.Bottom
		legend.Draw(r, legendDefaults)

		legendStyle.GetTextOptions().WriteToRenderer(r)

		ycursor := legendContent.Top
		tx := legendContent.Left
		legendCount := 0
		var label string
		for x := 0; x < len(labels); x++ {
			label = labels[x]
			if len(label) > 0 {
				if legendCount > 0 {
					ycursor += defaultMinimumTickVerticalSpacing
				}

				tb := r.MeasureText(label)

				ty := ycursor + tb.Height()
				r.Text(label, tx, ty)

				th2 := tb.Height() >> 1

				lx := tx + tb.Width() + lineTextGap
				ly := ty - th2
				lx2 := legendContent.Right - legendPadding.Right

				r.SetStrokeColor(lines[x].GetStrokeColor())
				r.SetStrokeWidth(lines[x].GetStrokeWidth())
				r.SetStrokeDashArray(lines[x].GetStrokeDashArray())

				r.MoveTo(lx, ly)
				r.LineTo(lx2, ly)
				r.Stroke()

				ycursor += tb.Height()
				legendCount++
			}
		}
	}
}
