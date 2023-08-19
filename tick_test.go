package unichart

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TickInput struct {
	min            float64
	max            float64
	num            int
	expectedResult []float64
}

// TestNiceTicks tests `niceTicks` function output
func TestNiceTicks(t *testing.T) {
	inputs := []TickInput{
		{
			min:            0,
			max:            10,
			num:            3,
			expectedResult: []float64{0.0, 5.0, 10.0},
		},
		{
			min:            0,
			max:            10,
			num:            2,
			expectedResult: []float64{0.0, 10.0},
		},
		{
			min:            0,
			max:            10,
			num:            5,
			expectedResult: []float64{0.0, 2.0, 4.0, 6.0, 8.0, 10.0},
		},
		{
			min:            -10,
			max:            10,
			num:            8,
			expectedResult: []float64{-10.0, -8.0, -6.0, -4.0, -2.0, 0.0, 2.0, 4.0, 6.0, 8.0, 10.0},
		},
		{
			min:            1,
			max:            5,
			num:            10,
			expectedResult: []float64{1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0},
		},
	}

	for _, i := range inputs {
		result := niceTicks(i.min, i.max, i.num)

		require.ElementsMatch(t, i.expectedResult, result)
	}
}
