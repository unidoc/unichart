package dataset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/dataset/sequence"
)

func TestConcatSeries(t *testing.T) {
	s1 := ContinuousSeries{
		XValues: sequence.LinearRange(1.0, 10.0),
		YValues: sequence.LinearRange(1.0, 10.0),
	}

	s2 := ContinuousSeries{
		XValues: sequence.LinearRange(11, 20.0),
		YValues: sequence.LinearRange(10.0, 1.0),
	}

	s3 := ContinuousSeries{
		XValues: sequence.LinearRange(21, 30.0),
		YValues: sequence.LinearRange(1.0, 10.0),
	}

	cs := ConcatSeries([]Series{s1, s2, s3})
	require.Equal(t, 30, cs.Len())

	x0, y0 := cs.GetValue(0)
	require.Equal(t, 1.0, x0)
	require.Equal(t, 1.0, y0)

	xm, ym := cs.GetValue(19)
	require.Equal(t, 20.0, xm)
	require.Equal(t, 1.0, ym)

	xn, yn := cs.GetValue(29)
	require.Equal(t, 30.0, xn)
	require.Equal(t, 10.0, yn)
}
