package series

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/data/sequence"
)

func TestPercentageDifferenceSeries(t *testing.T) {
	cs := ContinuousSeries{
		XValues: sequence.LinearRange(1.0, 10.0),
		YValues: sequence.LinearRange(1.0, 10.0),
	}

	pcs := PercentChangeSeries{
		Name:        "Test Series",
		InnerSeries: cs,
	}

	require.Equal(t, "Test Series", pcs.GetName())
	require.Equal(t, 10, pcs.Len())
	x0, y0 := pcs.GetValues(0)
	require.Equal(t, 1.0, x0)
	require.Equal(t, 0.0, y0)

	xn, yn := pcs.GetValues(9)
	require.Equal(t, 10.0, xn)
	require.Equal(t, 9.0, yn)

	xn, yn = pcs.GetLastValues()
	require.Equal(t, 10.0, xn)
	require.Equal(t, 9.0, yn)
}
