package render

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoxClone(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	b := a.Clone()

	require.True(t, a.Equals(b))
	require.True(t, b.Equals(a))
}

func TestBoxEquals(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	b := Box{Top: 10, Left: 10, Right: 30, Bottom: 30}
	c := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}

	require.True(t, a.Equals(a))
	require.True(t, a.Equals(c))
	require.True(t, c.Equals(a))
	require.False(t, a.Equals(b))
	require.False(t, c.Equals(b))
	require.False(t, b.Equals(a))
	require.False(t, b.Equals(c))
}

func TestBoxIsBiggerThan(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20}
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}

	require.True(t, a.IsBiggerThan(b))
	require.False(t, a.IsBiggerThan(c))
	require.True(t, c.IsBiggerThan(a))
}

func TestBoxIsSmallerThan(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20}
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}

	require.False(t, a.IsSmallerThan(b))
	require.True(t, a.IsSmallerThan(c))
	require.False(t, c.IsSmallerThan(a))
}

func TestBoxGrow(t *testing.T) {
	a := Box{Top: 1, Left: 2, Right: 15, Bottom: 15}
	b := Box{Top: 4, Left: 5, Right: 30, Bottom: 35}
	c := a.Grow(b)

	require.False(t, c.Equals(b))
	require.False(t, c.Equals(a))
	require.Equal(t, 1, c.Top)
	require.Equal(t, 2, c.Left)
	require.Equal(t, 30, c.Right)
	require.Equal(t, 35, c.Bottom)
}

func TestBoxFit(t *testing.T) {
	a := Box{Top: 64, Left: 64, Right: 192, Bottom: 192}
	b := Box{Top: 16, Left: 16, Right: 256, Bottom: 170}
	c := Box{Top: 16, Left: 16, Right: 170, Bottom: 256}

	fab := a.Fit(b)
	require.Equal(t, a.Left, fab.Left)
	require.Equal(t, a.Right, fab.Right)
	require.True(t, fab.Top < fab.Bottom)
	require.True(t, fab.Left < fab.Right)
	require.True(t, math.Abs(b.AspectRatio()-fab.AspectRatio()) < 0.02)

	fac := a.Fit(c)
	require.Equal(t, a.Top, fac.Top)
	require.Equal(t, a.Bottom, fac.Bottom)
	require.True(t, math.Abs(c.AspectRatio()-fac.AspectRatio()) < 0.02)
}

func TestBoxConstrain(t *testing.T) {
	a := Box{Top: 64, Left: 64, Right: 192, Bottom: 192}
	b := Box{Top: 16, Left: 16, Right: 256, Bottom: 170}
	c := Box{Top: 16, Left: 16, Right: 170, Bottom: 256}

	cab := a.Constrain(b)
	require.Equal(t, 64, cab.Top)
	require.Equal(t, 64, cab.Left)
	require.Equal(t, 192, cab.Right)
	require.Equal(t, 170, cab.Bottom)

	cac := a.Constrain(c)
	require.Equal(t, 64, cac.Top)
	require.Equal(t, 64, cac.Left)
	require.Equal(t, 170, cac.Right)
	require.Equal(t, 192, cac.Bottom)
}

func TestBoxOuterConstrain(t *testing.T) {
	box := NewBox(0, 0, 100, 100)
	canvas := NewBox(5, 5, 95, 95)
	taller := NewBox(-10, 5, 50, 50)

	c := canvas.OuterConstrain(box, taller)
	require.Equal(t, 15, c.Top, c.String())
	require.Equal(t, 5, c.Left, c.String())
	require.Equal(t, 95, c.Right, c.String())
	require.Equal(t, 95, c.Bottom, c.String())

	wider := NewBox(5, 5, 110, 50)
	d := canvas.OuterConstrain(box, wider)
	require.Equal(t, 5, d.Top, d.String())
	require.Equal(t, 5, d.Left, d.String())
	require.Equal(t, 85, d.Right, d.String())
	require.Equal(t, 95, d.Bottom, d.String())
}

func TestBoxShift(t *testing.T) {
	b := Box{Top: 5, Left: 5, Right: 10, Bottom: 10}
	shifted := b.Shift(1, 2)

	require.Equal(t, 7, shifted.Top)
	require.Equal(t, 6, shifted.Left)
	require.Equal(t, 11, shifted.Right)
	require.Equal(t, 12, shifted.Bottom)
}

func TestBoxCenter(t *testing.T) {
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 30}
	cx, cy := b.Center()

	require.Equal(t, 15, cx)
	require.Equal(t, 20, cy)
}

func TestBoxCornersCenter(t *testing.T) {
	bc := BoxCorners{
		TopLeft:     Point{5, 5},
		TopRight:    Point{15, 5},
		BottomRight: Point{15, 15},
		BottomLeft:  Point{5, 15},
	}
	cx, cy := bc.Center()

	require.Equal(t, 10, cx)
	require.Equal(t, 10, cy)
}

func TestBoxCornersRotate(t *testing.T) {
	bc := BoxCorners{
		TopLeft:     Point{5, 5},
		TopRight:    Point{15, 5},
		BottomRight: Point{15, 15},
		BottomLeft:  Point{5, 15},
	}
	rotated := bc.Rotate(45)

	require.True(t, rotated.TopLeft.Equals(Point{10, 3}), rotated.String())
}
