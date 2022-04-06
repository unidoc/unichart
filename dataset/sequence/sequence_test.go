package sequence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapperEach(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	values.Each(func(i int, v float64) {
		require.Equal(t, float64(i), v-1)
	})
}

func TestWrapperMap(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	mapped := values.Map(func(i int, v float64) float64 {
		require.Equal(t, float64(i), v-1)
		return v * 2
	})
	require.Equal(t, 4, mapped.Len())
}

func TestWrapperFoldLeft(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	ten := values.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	require.Equal(t, 10.0, ten)

	orderTest := Wrapper{NewArraySequence(10, 3, 2, 1)}
	four := orderTest.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	require.Equal(t, 4.0, four)
}

func TestWrapperFoldRight(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	ten := values.FoldRight(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	require.Equal(t, 10.0, ten)

	orderTest := Wrapper{NewArraySequence(10, 3, 2, 1)}
	notFour := orderTest.FoldRight(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	require.Equal(t, -14.0, notFour)
}

func TestWrapperSum(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	require.Equal(t, 10.0, values.Sum())
}

func TestWrapperAverage(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4)}
	require.Equal(t, 2.5, values.Average())

	valuesOdd := Wrapper{NewArraySequence(1, 2, 3, 4, 5)}
	require.Equal(t, 3.0, valuesOdd.Average())
}

func TestWrapperuenceVariance(t *testing.T) {
	values := Wrapper{NewArraySequence(1, 2, 3, 4, 5)}
	require.Equal(t, 2.0, values.Variance())
}

func TestSequenceNormalize(t *testing.T) {
	normalized := NewArrayWrapper(1, 2, 3, 4, 5).Normalize().Values()

	require.NotEmpty(t, normalized)
	require.Len(t, normalized, 5)
	require.Equal(t, 0.0, normalized[0])
	require.Equal(t, 0.25, normalized[1])
	require.Equal(t, 1.0, normalized[4])
}

func TestLinearRange(t *testing.T) {
	values := LinearRange(1, 100)
	require.Len(t, values, 100)
	require.Equal(t, 1.0, values[0])
	require.Equal(t, 100.0, values[99])
}

func TestLinearRangeWithStep(t *testing.T) {
	values := LinearRangeWithStep(0, 100, 5)
	require.Equal(t, 100.0, values[20])
	require.Len(t, values, 21)
}

func TestLinearRangeReversed(t *testing.T) {
	values := LinearRange(10.0, 1.0)
	require.Equal(t, 10, len(values))
	require.Equal(t, 10.0, values[0])
	require.Equal(t, 1.0, values[9])
}

func TestLinearSequenceRegression(t *testing.T) {
	linearProvider := NewLinearSequence().WithStart(1.0).WithEnd(100.0)
	require.Equal(t, 1.0, linearProvider.Start())
	require.Equal(t, 100.0, linearProvider.End())
	require.Equal(t, 100, linearProvider.Len())

	values := Wrapper{linearProvider}.Values()
	require.Len(t, values, 100)
	require.Equal(t, 1.0, values[0])
	require.Equal(t, 100.0, values[99])
}
