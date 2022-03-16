package series

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstValueAnnotation(t *testing.T) {
	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	fva := FirstValueAnnotation(series)
	require.NotEmpty(t, fva.Annotations)

	fvaa := fva.Annotations[0]
	require.Equal(t, 1.0, fvaa.XValue)
	require.Equal(t, 5.0, fvaa.YValue)
}

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
