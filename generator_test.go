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

func TestGeneratorControl(t *testing.T) {
	routes, err := route.RoutesForFrance()
	require.NoError(t, err)

	var isClosed bool
	gen := New(WithInterval(100 * time.Millisecond))
	gen.OnClose = func() {
		isClosed = true
	}
	gen.Run()

	t.Run("generator", func(t *testing.T) {
		t.Parallel()
		dev, err := Drone("myDrone", nil, routes...)
		require.NoError(t, err)

		// attach
		gen.Attach(dev)
		dev, err = gen.Lookup(dev.ID())
		require.NoError(t, err)
		require.NotNil(t, dev)

		// detach
		gen.Detach(dev.ID())

		// lookup
		dev, err = gen.Lookup(dev.ID())
		require.Error(t, err)
		require.Nil(t, dev)

		// close
		gen.Close()
		require.True(t, isClosed)
	})
}
