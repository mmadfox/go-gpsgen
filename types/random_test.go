package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom_MinMax(t *testing.T) {
	rnd := NewRandom(1, 100)
	require.Equal(t, 1, rnd.Min())
	require.Equal(t, 100, rnd.Max())

	rnd = NewRandom(-1, 100)
	require.Equal(t, 0, rnd.Min())
	require.Equal(t, 100, rnd.Max())

	rnd = NewRandom(-1, -1)
	require.Equal(t, 0, rnd.Min())
	require.Equal(t, 1, rnd.Max())

	rnd = NewRandom(0, 0)
	require.Equal(t, 0, rnd.Min())
	require.Equal(t, 1, rnd.Max())
}

func TestRandom_Value(t *testing.T) {
	values := make([]int, 0, 10)
	rnd := NewRandom(1, 10)
	require.NotNil(t, rnd)
	for i := 0; i < 10; i++ {
		val := rnd.Value()
		values = append(values, val)
	}
	require.GreaterOrEqual(t, values[0], 1)
	require.LessOrEqual(t, values[len(values)-1], 10)
}
