package dataset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValuesValues(t *testing.T) {
	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Values()
	require.Len(t, values, 7)
	require.Equal(t, 10.0, values[0])
	require.Equal(t, 9.0, values[1])
	require.Equal(t, 8.0, values[2])
	require.Equal(t, 7.0, values[3])
	require.Equal(t, 6.0, values[4])
	require.Equal(t, 5.0, values[5])
	require.Equal(t, 2.0, values[6])
}

func TestValuesValuesNormalized(t *testing.T) {
	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).ValuesNormalized()
	require.Len(t, values, 7)
	require.Equal(t, 0.2127, values[0])
	require.Equal(t, 0.0425, values[6])
}

func TestValuesNormalize(t *testing.T) {
	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Normalize()
	require.Len(t, values, 7)
	require.Equal(t, 0.2127, values[0].Value)
	require.Equal(t, 0.0425, values[6].Value)
}
