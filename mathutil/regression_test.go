package mathutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoly(t *testing.T) {
	// replaced new assertions helper
	var xGiven = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var yGiven = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
	var degree = 2

	c, err := PolyRegression(xGiven, yGiven, degree)
	require.Nil(t, err)
	require.Equal(t, 3, len(c))

	require.InDelta(t, c[0], 0.999999999, DefaultEpsilon)
	require.InDelta(t, c[1], 2, DefaultEpsilon)
	require.InDelta(t, c[2], 3, DefaultEpsilon)
}
