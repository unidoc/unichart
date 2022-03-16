package mathutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix(10, 5)
	rows, cols := m.Size()
	require.Equal(t, 10, rows)
	require.Equal(t, 5, cols)
	require.Zero(t, m.Get(0, 0))
	require.Zero(t, m.Get(9, 4))
}

func TestNewMatrixWithValues(t *testing.T) {
	m := NewMatrix(5, 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	rows, cols := m.Size()
	require.Equal(t, 5, rows)
	require.Equal(t, 2, cols)
	require.Equal(t, 1.0, m.Get(0, 0))
	require.Equal(t, 10.0, m.Get(4, 1))
}

func TestIdentitiyMatrix(t *testing.T) {
	id := IdentityMatrix(5)
	rows, cols := id.Size()
	require.Equal(t, 5, rows)
	require.Equal(t, 5, cols)
	require.Equal(t, 1.0, id.Get(0, 0))
	require.Equal(t, 1.0, id.Get(1, 1))
	require.Equal(t, 1.0, id.Get(2, 2))
	require.Equal(t, 1.0, id.Get(3, 3))
	require.Equal(t, 1.0, id.Get(4, 4))
	require.Equal(t, 0.0, id.Get(0, 1))
	require.Equal(t, 0.0, id.Get(1, 0))
	require.Equal(t, 0.0, id.Get(4, 0))
	require.Equal(t, 0.0, id.Get(0, 4))
}

func TestNewMatrixFromArrays(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	})
	require.NotNil(t, m)

	rows, cols := m.Size()
	require.Equal(t, 2, rows)
	require.Equal(t, 4, cols)
}

func TestOnes(t *testing.T) {
	ones := OnesMatrix(5, 10)
	rows, cols := ones.Size()
	require.Equal(t, 5, rows)
	require.Equal(t, 10, cols)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			require.Equal(t, 1.0, ones.Get(row, col))
		}
	}
}

func TestMatrixEpsilon(t *testing.T) {
	ones := OnesMatrix(2, 2)
	ones = ones.WithEpsilon(0.001)
	require.Equal(t, 0.001, ones.Epsilon())
}

func TestMatrixArrays(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})

	require.NotNil(t, m)

	arrays := m.Arrays()

	require.Equal(t, arrays, [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})
}

func TestMatrixIsSquare(t *testing.T) {
	require.False(t, NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}).IsSquare())

	require.False(t, NewMatrixFromArrays([][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}).IsSquare())

	require.True(t, NewMatrixFromArrays([][]float64{
		{1, 2},
		{3, 4},
	}).IsSquare())
}

func TestMatrixIsSymmetric(t *testing.T) {
	require.False(t, NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{2, 1, 2},
	}).IsSymmetric())

	require.False(t, NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}).IsSymmetric())

	require.True(t, NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{2, 1, 2},
		{3, 2, 1},
	}).IsSymmetric())

}

func TestMatrixGet(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	require.Equal(t, 1.0, m.Get(0, 0))
	require.Equal(t, 2.0, m.Get(0, 1))
	require.Equal(t, 3.0, m.Get(0, 2))
	require.Equal(t, 4.0, m.Get(1, 0))
	require.Equal(t, 5.0, m.Get(1, 1))
	require.Equal(t, 6.0, m.Get(1, 2))
	require.Equal(t, 7.0, m.Get(2, 0))
	require.Equal(t, 8.0, m.Get(2, 1))
	require.Equal(t, 9.0, m.Get(2, 2))
}

func TestMatrixSet(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m.Set(1, 1, 99)
	require.Equal(t, 99.0, m.Get(1, 1))
}

func TestMatrixCol(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	require.EqualValues(t, []float64{1, 4, 7}, m.Col(0))
	require.EqualValues(t, []float64{2, 5, 8}, m.Col(1))
	require.EqualValues(t, []float64{3, 6, 9}, m.Col(2))
}

func TestMatrixRow(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	require.EqualValues(t, []float64{1, 2, 3}, m.Row(0))
	require.EqualValues(t, []float64{4, 5, 6}, m.Row(1))
	require.EqualValues(t, []float64{7, 8, 9}, m.Row(2))
}

func TestMatrixSwapRows(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m.SwapRows(0, 1)

	require.EqualValues(t, []float64{4, 5, 6}, m.Row(0))
	require.EqualValues(t, []float64{1, 2, 3}, m.Row(1))
	require.EqualValues(t, []float64{7, 8, 9}, m.Row(2))
}

func TestMatrixCopy(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m2 := m.Copy()
	require.False(t, m == m2)
	require.True(t, m.Equals(m2))
}

func TestMatrixDiagonalVector(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	diag := m.DiagonalVector()
	require.EqualValues(t, []float64{1, 2, 3}, diag)
}

func TestMatrixDiagonalVectorLandscape(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 4, 7, 99},
		{4, 2, 8, 99},
	})

	diag := m.DiagonalVector()
	require.EqualValues(t, []float64{1, 2}, diag)
}

func TestMatrixDiagonalVectorPortrait(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 4},
		{4, 2},
		{99, 99},
	})

	diag := m.DiagonalVector()
	require.EqualValues(t, []float64{1, 2}, diag)
}

func TestMatrixDiagonal(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	m2 := NewMatrixFromArrays([][]float64{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 3},
	})

	require.True(t, m.Diagonal().Equals(m2))
}

func TestMatrixEquals(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	require.False(t, m.Equals(nil))
	var nilMatrix *Matrix
	require.True(t, nilMatrix.Equals(nil))
	require.False(t, m.Equals(NewMatrix(1, 1)))
	require.False(t, m.Equals(NewMatrix(3, 3)))
	require.True(t, m.Equals(NewMatrix(3, 3, 1, 4, 7, 4, 2, 8, 7, 8, 3)))
}

func TestMatrixL(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	l := m.L()
	require.True(t, l.Equals(NewMatrix(3, 3, 1, 2, 3, 0, 5, 6, 0, 0, 9)))
}

func TestMatrixU(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	u := m.U()
	require.True(t, u.Equals(NewMatrix(3, 3, 0, 0, 0, 4, 0, 0, 7, 8, 0)))
}

func TestMatrixString(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	require.Equal(t, "1 2 3 \n4 5 6 \n7 8 9 \n", m.String())
}

func TestMatrixLU(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 3, 5},
		{2, 4, 7},
		{1, 1, 0},
	})

	l, u, p := m.LU()
	require.NotNil(t, l)
	require.NotNil(t, u)
	require.NotNil(t, p)
}

func TestMatrixQR(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	})

	q, r := m.QR()
	require.NotNil(t, q)
	require.NotNil(t, r)
}

func TestMatrixTranspose(t *testing.T) {
	m := NewMatrixFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})

	m2 := m.Transpose()

	rows, cols := m2.Size()
	require.Equal(t, 3, rows)
	require.Equal(t, 4, cols)

	require.Equal(t, 1.0, m2.Get(0, 0))
	require.Equal(t, 10.0, m2.Get(0, 3))
	require.Equal(t, 3.0, m2.Get(2, 0))
}
