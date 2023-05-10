package navigator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNavigatorLocation(t *testing.T) {
	routepath := []Point{
		{X: 55.74966429698134, Y: 37.624339525581576},
		{X: 55.748482140161286, Y: 37.62444198526788},
	}

	routedist := DistanceTo(
		routepath[0].X,
		routepath[0].Y,
		routepath[1].X,
		routepath[1].Y,
	)

	r1, err := NewRoute([][]Point{routepath})
	require.NoError(t, err)

	nav, err := New(SkipOfflineMode())
	require.NoError(t, err)
	nav.AddRoute(r1)

	speed := 1.0       // m/s
	tick := 1.0        // second
	maxDistance := 1.5 // meters
	var curDist float64
	for i := 0; i < int(routedist); i++ {
		ok := nav.Next(tick, speed)
		require.True(t, nav.IsOnline())
		require.True(t, ok)
		loc := nav.Location()
		dist := DistanceTo(
			routepath[0].X,
			routepath[0].Y,
			loc.Lat,
			loc.Lon,
		)

		lat, lon := DestinationPoint(routepath[0].X, routepath[0].Y, float64(i), nav.Segment().Bearing())
		distDiff := DistanceTo(lat, lon, loc.Lat, loc.Lon)
		require.True(t, distDiff < maxDistance)

		if i == 0 {
			curDist = dist
			continue
		}
		require.True(t, dist > curDist)
		curDist = dist
	}

	require.True(t, (routedist-curDist) < maxDistance)
}

func TestNavigatorMultiTracks(t *testing.T) {
	routepath := [][]Point{
		{
			{X: 55.74966429698134, Y: 37.624339525581576},
			{X: 55.748482140161286, Y: 37.62444198526788},
		}, // track-1
		{
			{X: 55.74969397655789, Y: 37.62458990661304},
			{X: 55.74847692761699, Y: 37.62472920358056},
		}, // track-2
	}

	r1, err := NewRoute(routepath)
	require.NoError(t, err)

	nav, err := New(SkipOfflineMode())
	require.NoError(t, err)
	nav.AddRoute(r1)

	speed := 1.0 // m/s
	tick := 1.0  // second
	var curDist float64
	maxDistance := 2.0
	for i := 0; i < int(r1.TotalDistance()); i++ {
		ok := nav.Next(tick, speed)
		endOfFirstSegment := i == 131
		if endOfFirstSegment {
			require.False(t, ok)
		} else {
			require.True(t, ok)
		}
		loc := nav.Location()
		curDist = loc.CurrentDistance
	}
	require.True(t, (r1.TotalDistance()-curDist) < maxDistance)
}

func TestNavigatorMultiRoutes(t *testing.T) {
	r1, err := NewRoute([][]Point{
		{
			{X: 55.74966429698134, Y: 37.624339525581576},
			{X: 55.748482140161286, Y: 37.62444198526788},
		},
	})
	require.NoError(t, err)

	r2, err := NewRoute([][]Point{
		{
			{X: 55.74969397655789, Y: 37.62458990661304},
			{X: 55.74847692761699, Y: 37.62472920358056},
		},
	})
	require.NoError(t, err)

	nav, err := New(SkipOfflineMode())
	require.NoError(t, err)
	nav.AddRoutes(r1, r2)

	totalDist := r1.TotalDistance() + r2.TotalDistance()
	var curDist float64
	tick := 1.0
	speed := 1.0
	maxDistance := 2.0
	for i := 0; i < int(totalDist); i++ {
		nav.Next(tick, speed)
		loc := nav.Location()
		curDist = loc.CurrentDistance
	}
	require.True(t, (totalDist-curDist) < maxDistance)
	require.Equal(t, 1, nav.RouteIndex())
}
