package gpsgen

import (
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	var buf []byte
	var err error
	for i := 0; i < 10; i++ {
		buf, err = proto.MarshalOptions{}.MarshalAppend(buf[:0], &pb.Device{
			UserId: "userID",
		})
		if err != nil {
			panic(err)
		}
	}
}

func TestGenerator(t *testing.T) {
	routes, err := route.RoutesForChina()
	require.NoError(t, err)

	gen := New(WithInterval(100 * time.Millisecond))

	d1, err := Drone("Tx", nil, routes...)
	require.NoError(t, err)

	var tick uint32
	d1.OnStateChange = func(dev *pb.Device) {
		atomic.AddUint32(&tick, 1)
		require.NotNil(t, dev)
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
		dev1 := dev
		for i := 0; i < 3; i++ {
			dev, err = gen.Lookup(dev1.ID())
			time.Sleep(50 * time.Millisecond)
		}
		require.NoError(t, err)
		require.NotNil(t, dev)

		// detach
		gen.Detach(dev.ID())

		// lookup
		dev2 := dev
		for i := 0; i < 3; i++ {
			dev, err = gen.Lookup(dev2.ID())
			time.Sleep(50 * time.Millisecond)
		}
		require.Error(t, err)
		require.Nil(t, dev)

		// close
		gen.Close()
		require.True(t, isClosed)
	})
}
