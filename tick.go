package unichart

import (
	"math"
	"strconv"

	"github.com/unidoc/unichart/dataset"
	"github.com/unidoc/unichart/dataset/sequence"
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
	GetTicks(r render.Renderer, defaults render.Style, vf dataset.ValueFormatter) []Tick
}

// Tick represents a label on an axis.
type Tick struct {
	Value float64
	Label string
}

// generateContinuousTicks generates a set of ticks.
func generateContinuousTicks(r render.Renderer, ra sequence.Range, isVertical bool, style render.Style, vf dataset.ValueFormatter) []Tick {
	if vf == nil {
		vf = dataset.FloatValueFormatter
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

	// Decide to use nice ticks by checking if the first tick label is a float number.
	_, err := strconv.ParseFloat(ticks[0].Label, 64)
	if err == nil {
		intermediateTickCount = mathutil.MinInt(intermediateTickCount, defaultTickCountSanityCheck)

		// Make sure to have at least two ticks generated.
		numTicks := mathutil.MaxInt(intermediateTickCount, 2)
		nTicks := niceTicks(min, max, numTicks)
		if ra.IsDescending() {
			// Reverse generated ticks.
			for i := len(nTicks)/2 - 1; i >= 0; i-- {
				opp := len(nTicks) - 1 - i
				nTicks[i], nTicks[opp] = nTicks[opp], nTicks[i]
			}
		}

		ticks = []Tick{}
		for _, t := range nTicks {
			ticks = append(ticks, Tick{
				Value: t,
				Label: vf(t),
			})
		}
	} else {
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
	}

	return ticks
}

// niceNum would round tick value to the nearest nice number.
func niceNum(value float64, round bool) float64 {
	exponent := math.Floor(math.Log10(value))
	fraction := value / math.Pow(10, exponent)

	var niceFraction float64
	if round {
		if fraction < 1.5 {
			niceFraction = 1
		} else if fraction < 3 {
			niceFraction = 2
		} else if fraction < 7 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	} else {
		if fraction <= 1 {
			niceFraction = 1
		} else if fraction <= 2 {
			niceFraction = 2
		} else if fraction <= 5 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	}

	return niceFraction * math.Pow(10, exponent)
}

// niceTicks generates ticks value with a rounded up values.
func niceTicks(min, max float64, numTicks int) (tickValues []float64) {
	rangeValue := niceNum(max-min, false)
	tickSpacing := niceNum(rangeValue/(float64(numTicks)-1), true)
	niceMin := math.Floor(min/tickSpacing) * tickSpacing
	niceMax := math.Ceil(max/tickSpacing) * tickSpacing

	for value := niceMin; value <= niceMax; value += tickSpacing {
		tickValues = append(tickValues, value)
	}

	return tickValues
}
