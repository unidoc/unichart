package chart

import (
	"math"

	"github.com/unidoc/unichart/data"
	"github.com/unidoc/unichart/mathutil"
	"github.com/unidoc/unichart/render"
)

// TickPosition is an enumeration of possible tick drawing positions.
type TickPosition int

const (
	// TickPositionUnset means to use the default tick position.
	TickPositionUnset TickPosition = 0

	// TickPositionBetweenTicks draws the labels for a tick between the previous and current tick.
	TickPositionBetweenTicks TickPosition = 1

	// TickPositionUnderTick draws the tick below the tick.
	TickPositionUnderTick TickPosition = 2
)

// TicksProvider is a type that provides ticks.
type TicksProvider interface {
	GetTicks(r render.Renderer, defaults render.Style, vf data.ValueFormatter) []Tick
}

// Tick represents a label on an axis.
type Tick struct {
	Value float64
	Label string
}

// generateContinuousTicks generates a set of ticks.
func generateContinuousTicks(r render.Renderer, ra data.Range, isVertical bool, style render.Style, vf data.ValueFormatter) []Tick {
	if vf == nil {
		vf = data.FloatValueFormatter
	}

	var ticks []Tick
	min, max := ra.GetMin(), ra.GetMax()

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	}

	minLabel := vf(min)
	style.GetTextOptions().WriteToRenderer(r)
	labelBox := r.MeasureText(minLabel)

	var tickSize float64
	if isVertical {
		tickSize = float64(labelBox.Height() + defaultMinimumTickVerticalSpacing)
	} else {
		tickSize = float64(labelBox.Width() + defaultMinimumTickHorizontalSpacing)
	}

	domain := float64(ra.GetDomain())
	domainRemainder := domain - (tickSize * 2)
	intermediateTickCount := int(math.Floor(float64(domainRemainder) / float64(tickSize)))

	rangeDelta := math.Abs(max - min)
	tickStep := rangeDelta / float64(intermediateTickCount)

	roundTo := mathutil.RoundTo(rangeDelta) / 10
	intermediateTickCount = mathutil.MinInt(intermediateTickCount, defaultTickCountSanityCheck)

	for x := 1; x < intermediateTickCount; x++ {
		var tickValue float64
		if ra.IsDescending() {
			tickValue = max - mathutil.RoundUp(tickStep*float64(x), roundTo)
		} else {
			tickValue = min + mathutil.RoundUp(tickStep*float64(x), roundTo)
		}
		ticks = append(ticks, Tick{
			Value: tickValue,
			Label: vf(tickValue),
		})
	}

	if ra.IsDescending() {
		ticks = append(ticks, Tick{
			Value: min,
			Label: vf(min),
		})
	} else {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	}

	return ticks
}
