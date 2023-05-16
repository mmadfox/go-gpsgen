package gpsgen

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/mmadfox/go-gpsgen/route"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	routes, err := route.RoutesForChina()
	require.NoError(t, err)

	gen := New(WithInterval(100 * time.Millisecond))

	d1, err := Drone("Tx", nil, routes...)
	require.NoError(t, err)

	var tick uint32
	d1.OnStateChange = func(s *State, snapshot []byte) {
		atomic.AddUint32(&tick, 1)
		require.NotNil(t, s)
		require.NotZero(t, snapshot)
	}

	gen.Attach(d1)
	gen.Run()

	<-time.After(500 * time.Millisecond)

	require.GreaterOrEqual(t, uint32(5), atomic.LoadUint32(&tick))
}
