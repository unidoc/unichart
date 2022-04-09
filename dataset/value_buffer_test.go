package dataset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	buffer := NewValueBuffer()

	buffer.Enqueue(1)
	require.Equal(t, 1, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 1.0, buffer.PeekBack())

	buffer.Enqueue(2)
	require.Equal(t, 2, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 2.0, buffer.PeekBack())

	buffer.Enqueue(3)
	require.Equal(t, 3, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 3.0, buffer.PeekBack())

	buffer.Enqueue(4)
	require.Equal(t, 4, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 4.0, buffer.PeekBack())

	buffer.Enqueue(5)
	require.Equal(t, 5, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 5.0, buffer.PeekBack())

	buffer.Enqueue(6)
	require.Equal(t, 6, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 6.0, buffer.PeekBack())

	buffer.Enqueue(7)
	require.Equal(t, 7, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 7.0, buffer.PeekBack())

	buffer.Enqueue(8)
	require.Equal(t, 8, buffer.Len())
	require.Equal(t, 1.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value := buffer.Dequeue()
	require.Equal(t, 1.0, value)
	require.Equal(t, 7, buffer.Len())
	require.Equal(t, 2.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 2.0, value)
	require.Equal(t, 6, buffer.Len())
	require.Equal(t, 3.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 3.0, value)
	require.Equal(t, 5, buffer.Len())
	require.Equal(t, 4.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 4.0, value)
	require.Equal(t, 4, buffer.Len())
	require.Equal(t, 5.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 5.0, value)
	require.Equal(t, 3, buffer.Len())
	require.Equal(t, 6.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 6.0, value)
	require.Equal(t, 2, buffer.Len())
	require.Equal(t, 7.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 7.0, value)
	require.Equal(t, 1, buffer.Len())
	require.Equal(t, 8.0, buffer.Peek())
	require.Equal(t, 8.0, buffer.PeekBack())

	value = buffer.Dequeue()
	require.Equal(t, 8.0, value)
	require.Equal(t, 0, buffer.Len())
	require.Zero(t, buffer.Peek())
	require.Zero(t, buffer.PeekBack())
}

func TestBufferClear(t *testing.T) {
	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)

	require.Equal(t, 8, buffer.Len())

	buffer.Clear()
	require.Equal(t, 0, buffer.Len())
	require.Zero(t, buffer.Peek())
	require.Zero(t, buffer.PeekBack())
}

func TestBufferArray(t *testing.T) {
	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(2)
	buffer.Enqueue(3)
	buffer.Enqueue(4)
	buffer.Enqueue(5)

	contents := buffer.Array()
	require.Len(t, contents, 5)
	require.Equal(t, 1.0, contents[0])
	require.Equal(t, 2.0, contents[1])
	require.Equal(t, 3.0, contents[2])
	require.Equal(t, 4.0, contents[3])
	require.Equal(t, 5.0, contents[4])
}

func TestBufferEach(t *testing.T) {
	buffer := NewValueBuffer()

	for x := 1; x < 17; x++ {
		buffer.Enqueue(float64(x))
	}

	called := 0
	buffer.Each(func(_ int, v float64) {
		if v == float64(called+1) {
			called++
		}
	})

	require.Equal(t, 16, called)
}

func TestNewBuffer(t *testing.T) {
	empty := NewValueBuffer()
	require.NotNil(t, empty)
	require.Zero(t, empty.Len())
	require.Equal(t, bufferDefaultCapacity, empty.Capacity())
	require.Zero(t, empty.Peek())
	require.Zero(t, empty.PeekBack())
}

func TestNewBufferWithValues(t *testing.T) {
	values := NewValueBuffer(1, 2, 3, 4, 5)
	require.NotNil(t, values)
	require.Equal(t, 5, values.Len())
	require.Equal(t, 1.0, values.Peek())
	require.Equal(t, 5.0, values.PeekBack())
}

func TestBufferGrowth(t *testing.T) {
	values := NewValueBuffer(1, 2, 3, 4, 5)
	for i := 0; i < 1<<10; i++ {
		values.Enqueue(float64(i))
	}

	require.Equal(t, float64(1<<10-1), values.PeekBack())
}
