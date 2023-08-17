package gpsgen

import (
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestGenerator_Attach(t *testing.T) {
	tests := []struct {
		name   string
		assert func(*Generator, *Device)
	}{
		{
			name: "should not return error when device is nil",
			assert: func(g *Generator, d *Device) {
				require.NoError(t, g.Attach(nil))
			},
		},
		{
			name: "should return error when attach device more than one times",
			assert: func(g *Generator, d *Device) {
				require.NoError(t, g.Attach(d))
				require.Error(t, g.Attach(d))
			},
		},
		{
			name: "should not return error when successfully attached",
			assert: func(g *Generator, d *Device) {
				require.NoError(t, g.Attach(d))
				require.Equal(t, 1, g.NumDevices())
				require.Equal(t, Running, d.Status())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(nil)

			dev, err := NewDevice(nil)
			require.NoError(t, err)

			tt.assert(gen, dev)
		})
	}
}

func TestGenerator_Detach(t *testing.T) {
	tests := []struct {
		name   string
		assert func(*Generator, *Device)
	}{
		{
			name: "should not return error when deviceID is empty",
			assert: func(g *Generator, d *Device) {
				require.NoError(t, g.Detach(""))
			},
		},
		{
			name: "should not return error when device not found",
			assert: func(g *Generator, d *Device) {
				require.NoError(t, g.Detach("someid"))
			},
		},
		{
			name: "should not return error when successfully detached",
			assert: func(g *Generator, d *Device) {
				require.Equal(t, Running, d.Status())
				require.Equal(t, 1, g.NumDevices())
				require.NoError(t, g.Detach(d.ID()))
				require.Equal(t, Stopped, d.Status())
				require.Equal(t, 0, g.NumDevices())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(nil)

			dev, err := NewDevice(nil)
			require.NoError(t, err)
			require.NoError(t, gen.Attach(dev))

			tt.assert(gen, dev)
		})
	}
}

func TestGenerator_Lookup(t *testing.T) {
	tests := []struct {
		name   string
		assert func(*Generator, string)
	}{
		{
			name: "should not return error when deviceID is empty",
			assert: func(g *Generator, id string) {
				dev, ok := g.Lookup("")
				require.Nil(t, dev)
				require.False(t, ok)
			},
		},
		{
			name: "should not return error when device not found",
			assert: func(g *Generator, id string) {
				dev, ok := g.Lookup("someid")
				require.Nil(t, dev)
				require.False(t, ok)
			},
		},
		{
			name: "should not return error when successfully found",
			assert: func(g *Generator, id string) {
				dev, ok := g.Lookup(id)
				require.NotNil(t, dev)
				require.True(t, ok)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := New(nil)

			dev, err := NewDevice(nil)
			require.NoError(t, err)
			require.NoError(t, gen.Attach(dev))

			tt.assert(gen, dev.ID())
		})
	}
}

func TestGenerator_Run(t *testing.T) {
	gen := New(&Options{
		Interval:   10 * time.Millisecond,
		NumWorkers: 64,
	})

	var next uint32
	var packet uint32
	var e uint32

	gen.OnNext(func() {
		atomic.AddUint32(&next, 1)
	})

	gen.OnPacket(func(pck []byte) {
		atomic.AddUint32(&packet, uint32(len(pck)))
		p := new(pb.Packet)
		require.NoError(t, proto.Unmarshal(pck, p))
		require.NotEmpty(t, p.Devices)
	})

	gen.OnError(func(err error) {
		if err != nil {
			atomic.AddUint32(&e, 1)
		}
	})

	devices := make([]*Device, 0)
	for i := 0; i < 1000; i++ {
		// make new device
		devOpts := NewDeviceOptions()
		devOpts.Navigator.SkipOffline = true
		dev, err := NewDevice(nil)
		require.NoError(t, err)
		require.NotNil(t, dev)
		// add routes
		routes := testRoutes()
		dev.AddRoute(routes...)
		devices = append(devices, dev)
		// attach
		gen.Attach(dev)
	}

	go func() {
		<-time.After(50 * time.Millisecond)
		gen.Close()
	}()

	gen.Run()

	require.NotZero(t, atomic.LoadUint32(&packet))
	require.NotZero(t, atomic.LoadUint32(&next))
	require.Zero(t, atomic.LoadUint32(&e))

	for _, dev := range devices {
		state := dev.State()
		require.NotZero(t, state.Distance.CurrentDistance)
	}
}

func TestGenerator_Close(t *testing.T) {
	gen := New(&Options{
		Interval:   10 * time.Millisecond,
		PacketSize: -1, // default 32
	})
	for i := 0; i < 1000; i++ {
		devOpts := NewDeviceOptions()
		devOpts.Navigator.SkipOffline = true
		dev, err := NewDevice(nil)
		require.NoError(t, err)
		require.NotNil(t, dev)
		routes := testRoutes()
		dev.AddRoute(routes...)
		gen.Attach(dev)
	}
	gen.OnNext(func() {
		time.Sleep(50 * time.Millisecond)
	})
	gen.OnPacket(func(_ []byte) {
		time.Sleep(50 * time.Millisecond)
	})
	go func() {
		<-time.After(50 * time.Millisecond)
		gen.Close()
	}()
	gen.Run()
}
