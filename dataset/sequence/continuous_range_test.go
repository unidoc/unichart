package sequence

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/mathutil"
)

func TestRangeTranslate(t *testing.T) {
	// replaced new assertions helper
	values := []float64{1.0, 2.0, 2.5, 2.7, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	r := ContinuousRange{Domain: 1000}
	r.Min, r.Max = mathutil.MinMax(values...)

	// delta = ~7.0
	// value = ~5.0
	// domain = ~1000
	// 5/8 * 1000 ~=
	require.Equal(t, 0, r.Translate(1.0))
	require.Equal(t, 1000, r.Translate(8.0))
	require.Equal(t, 572, r.Translate(5.0))
}

func TestRangeTranslateEqaulMinMax(t *testing.T) {
	r := ContinuousRange{
		Domain: 370,
		Min:    0,
		Max:    0,
	}

	result := r.Translate(0)
	require.Equal(t, 0, result)
}
