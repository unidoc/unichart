package series

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeSeriesGetValue(t *testing.T) {
	ts := TimeSeries{
		Name: "Test",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}

	x0, y0 := ts.GetValues(0)
	require.NotZero(t, x0)
	require.Equal(t, 1.0, y0)
}

func TestTimeSeriesValidate(t *testing.T) {
	cs := TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	require.Nil(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			time.Now().AddDate(0, 0, -5),
			time.Now().AddDate(0, 0, -4),
			time.Now().AddDate(0, 0, -3),
			time.Now().AddDate(0, 0, -2),
			time.Now().AddDate(0, 0, -1),
		},
	}
	require.NotNil(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	require.NotNil(t, cs.Validate())
}
