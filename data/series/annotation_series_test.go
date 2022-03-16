package series

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLastValueAnnotationSeries(t *testing.T) {
	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	lva := LastValueAnnotationSeries(series)
	require.NotEmpty(t, lva.Annotations)
	lvaa := lva.Annotations[0]
	require.Equal(t, 5.0, lvaa.XValue)
	require.Equal(t, 1.0, lvaa.YValue)
}
