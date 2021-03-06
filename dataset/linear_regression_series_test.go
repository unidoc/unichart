package dataset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/dataset/sequence"
)

func TestLinearRegressionSeries(t *testing.T) {
	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: sequence.LinearRange(1.0, 100.0),
		YValues: sequence.LinearRange(1.0, 100.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	require.InDelta(t, 1.0, lrx0, 0.0000001)
	require.InDelta(t, 1.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	require.InDelta(t, 100.0, lrxn, 0.0000001)
	require.InDelta(t, 100.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesDesc(t *testing.T) {
	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: sequence.LinearRange(100.0, 1.0),
		YValues: sequence.LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	require.InDelta(t, 100.0, lrx0, 0.0000001)
	require.InDelta(t, 100.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	require.InDelta(t, 1.0, lrxn, 0.0000001)
	require.InDelta(t, 1.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesWindowAndOffset(t *testing.T) {
	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: sequence.LinearRange(100.0, 1.0),
		YValues: sequence.LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
		Offset:      10,
		Limit:       10,
	}
	require.Equal(t, 10, linRegSeries.Len())

	lrx0, lry0 := linRegSeries.GetValues(0)
	require.InDelta(t, 90.0, lrx0, 0.0000001)
	require.InDelta(t, 90.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	require.InDelta(t, 80.0, lrxn, 0.0000001)
	require.InDelta(t, 80.0, lryn, 0.0000001)
}
