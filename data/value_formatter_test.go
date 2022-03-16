package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeValueFormatterWithFormat(t *testing.T) {
	d := time.Now()
	di := d.UnixNano()
	df := float64(di)

	s := formatTime(d, defaultDateFormat)
	si := formatTime(di, defaultDateFormat)
	sf := formatTime(df, defaultDateFormat)
	require.Equal(t, s, si)
	require.Equal(t, s, sf)

	sd := TimeValueFormatter(d)
	sdi := TimeValueFormatter(di)
	sdf := TimeValueFormatter(df)
	require.Equal(t, s, sd)
	require.Equal(t, s, sdi)
	require.Equal(t, s, sdf)
}

func TestFloatValueFormatter(t *testing.T) {
	require.Equal(t, "1234.00", FloatValueFormatter(1234.00))
}

func TestFloatValueFormatterWithFloat32Input(t *testing.T) {
	require.Equal(t, "1234.00", FloatValueFormatter(float32(1234.00)))
}

func TestFloatValueFormatterWithIntegerInput(t *testing.T) {
	require.Equal(t, "1234.00", FloatValueFormatter(1234))
}

func TestFloatValueFormatterWithInt64Input(t *testing.T) {
	require.Equal(t, "1234.00", FloatValueFormatter(int64(1234)))
}

func TestFloatValueFormatterWithFormat(t *testing.T) {
	v := 123.456
	sv := FloatValueFormatterWithFormat(v, "%.3f")
	require.Equal(t, "123.456", sv)
	require.Equal(t, "123.000", FloatValueFormatterWithFormat(123, "%.3f"))
}
