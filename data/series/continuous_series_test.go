package series

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unidoc/unichart/data/sequence"
)

func TestContinuousSeries(t *testing.T) {
	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: sequence.LinearRange(1.0, 10.0),
		YValues: sequence.LinearRange(1.0, 10.0),
	}

	require.Equal(t, "Test Series", cs.GetName())
	require.Equal(t, 10, cs.Len())
	x0, y0 := cs.GetValues(0)
	require.Equal(t, 1.0, x0)
	require.Equal(t, 1.0, y0)

	xn, yn := cs.GetValues(9)
	require.Equal(t, 10.0, xn)
	require.Equal(t, 10.0, yn)

	xn, yn = cs.GetLastValues()
	require.Equal(t, 10.0, xn)
	require.Equal(t, 10.0, yn)
}

func TestContinuousSeriesValueFormatter(t *testing.T) {
	cs := ContinuousSeries{
		XValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f foo", v)
		},
		YValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f bar", v)
		},
	}

	xf, yf := cs.GetValueFormatters()
	require.Equal(t, "0.100000 foo", xf(0.1))
	require.Equal(t, "0.100000 bar", yf(0.1))
}

func TestContinuousSeriesValidate(t *testing.T) {
	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: sequence.LinearRange(1.0, 10.0),
		YValues: sequence.LinearRange(1.0, 10.0),
	}
	require.Nil(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		XValues: sequence.LinearRange(1.0, 10.0),
	}
	require.NotNil(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		YValues: sequence.LinearRange(1.0, 10.0),
	}
	require.NotNil(t, cs.Validate())
}
