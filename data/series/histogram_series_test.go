package series

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/data/sequence"
)

func TestHistogramSeries(t *testing.T) {
	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: sequence.LinearRange(1.0, 20.0),
		YValues: sequence.LinearRange(10.0, -10.0),
	}

	hs := HistogramSeries{
		InnerSeries: cs,
	}

	for x := 0; x < hs.Len(); x++ {
		csx, csy := cs.GetValues(0)
		hsx, hsy1, hsy2 := hs.GetBoundedValues(0)
		require.Equal(t, csx, hsx)
		require.True(t, hsy1 > 0)
		require.True(t, hsy2 <= 0)
		require.True(t, csy < 0 || (csy > 0 && csy == hsy1))
		require.True(t, csy > 0 || (csy < 0 && csy == hsy2))
	}
}
