package navigator

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
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
		dist := DistanceTo(
			routepath[0].X,
			routepath[0].Y,
			nav.point.X,
			nav.point.Y,
		)

		lat, lon := DestinationPoint(routepath[0].X, routepath[0].Y, float64(i), nav.Segment().Bearing())
		distDiff := DistanceTo(lat, lon, nav.point.X, nav.point.Y)
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
		curDist = nav.CurrentDistance()
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
		curDist = nav.CurrentDistance()
	}
	require.True(t, (totalDist-curDist) < maxDistance)
	require.Equal(t, 1, nav.RouteIndex())
}

func TestNavigatorToProto(t *testing.T) {
	r1, err := NewRoute([][]Point{
		{
			{X: 55.74966429698134, Y: 37.624339525581576},
			{X: 55.748482140161286, Y: 37.62444198526788},
		},
	})
	require.NoError(t, err)
	nav, err := New()
	require.NoError(t, err)
	nav.AddRoutes(r1)
	// walk
	for i := 0; i < 5; i++ {
		nav.Next(1, 1)
	}
	protoNav := nav.ToProto()
	require.NotNil(t, protoNav)

	require.Equal(t, nav.TotalRoutes(), len(protoNav.Routes))
	routes := nav.AllRoutes()
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		protoRoute := protoNav.Routes[i]
		require.Equal(t, route.Distance(), protoRoute.Distance)
		require.Equal(t, route.NumTracks(), len(protoRoute.Tracks))
		for j := 0; j < route.NumTracks(); j++ {
			segmentIndex := 0
			route.EachSegment(j, func(s *Segment) {
				protoSegment := protoRoute.Tracks[j].Segmenets[segmentIndex]
				require.NotNil(t, protoSegment)
				segmentIndex++
				require.Equal(t, s.Bearing(), protoSegment.Bearing)
				require.Equal(t, s.Dist(), protoSegment.Distance)
				require.Equal(t, s.Rel(), int(protoSegment.Rel))
				require.Equal(t, s.PointA().X, protoSegment.PointA.Lat)
				require.Equal(t, s.PointA().Y, protoSegment.PointA.Lon)
				require.Equal(t, s.PointB().X, protoSegment.PointB.Lat)
				require.Equal(t, s.PointB().Y, protoSegment.PointB.Lon)
			})
		}
	}

	require.Equal(t, nav.CurrentDistance(), protoNav.CurrentDistance)
	require.Equal(t, nav.RouteIndex(), int(protoNav.RouteIndex))
	require.Equal(t, nav.TrackIndex(), int(protoNav.TrackIndex))
	require.Equal(t, nav.SegmentIndex(), int(protoNav.SegmentIndex))
	require.Equal(t, nav.offlineIndex, int(protoNav.OfflineIndex))
	require.Equal(t, nav.CurrentDistance(), protoNav.CurrentDistance)
	require.Equal(t, nav.point.X, protoNav.Point.Lat)
	require.Equal(t, nav.point.Y, protoNav.Point.Lon)
	require.Equal(t, nav.offline.Min(), int(protoNav.OfflineMin))
	require.Equal(t, nav.offline.Max(), int(protoNav.OfflineMax))
	require.Equal(t, nav.TotalDistance(), protoNav.TotalDistance)
	require.NotNil(t, protoNav.Elevation)
}

func TestNavigatorFromProto(t *testing.T) {
	protoNav := &proto.NavigatorState{
		RouteIndex:      1,
		TrackIndex:      1,
		SegmentIndex:    1,
		SegmentDistance: 100,
		CurrentDistance: 10,
		OfflineIndex:    3,
		Point:           &proto.Point{Lat: 3, Lon: 4},
		Elevation: &proto.SensorState{
			Min:  1,
			Max:  4,
			ValX: 2,
			Gen: &proto.Curve{
				Points: []*proto.Curve_ControlPoint{
					{
						Vp: &proto.Curve_Point{X: 1, Y: 2},
						Cp: &proto.Curve_Point{X: 3, Y: 4},
					},
				},
			},
		},
		OfflineMin:    1,
		OfflineMax:    3,
		TotalDistance: 300,
		SkipOffline:   false,
		Routes: []*proto.Route{
			{
				Distance: 1000,
				Tracks: []*proto.Route_Track{
					{
						Segmenets: []*proto.Route_Track_Segment{
							{
								PointA:   &proto.Point{Lat: 1, Lon: 2},
								PointB:   &proto.Point{Lat: 3, Lon: 4},
								Bearing:  5,
								Distance: 33,
								Rel:      -1,
							},
						},
					},
				},
			},
		},
	}

	nav := new(Navigator)
	nav.FromProto(protoNav)
	require.Equal(t, len(protoNav.Routes), nav.TotalRoutes())
	routes := nav.AllRoutes()
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		require.Equal(t, protoNav.Routes[i].Distance, route.Distance())
		require.Equal(t, len(protoNav.Routes[i].Tracks), route.NumTracks())
		for s := 0; s < len(protoNav.Routes[i].Tracks); s++ {
			for si := 0; si < len(protoNav.Routes[i].Tracks[s].Segmenets); si++ {
				protoSeg := protoNav.Routes[i].Tracks[s].Segmenets[si]
				seg := route.SegmentAt(s, si)
				require.Equal(t, protoSeg.Bearing, seg.Bearing())
				require.Equal(t, protoSeg.Distance, seg.Dist())
				require.Equal(t, protoSeg.PointA.Lat, seg.PointA().X)
				require.Equal(t, protoSeg.PointA.Lon, seg.PointA().Y)
				require.Equal(t, protoSeg.PointB.Lat, seg.PointB().X)
				require.Equal(t, protoSeg.PointB.Lon, seg.PointB().Y)
				require.Equal(t, int(protoSeg.Rel), seg.Rel())
			}
		}
	}
	require.Equal(t, int(protoNav.RouteIndex), nav.routeIndex)
	require.Equal(t, int(protoNav.TrackIndex), nav.trackIndex)
	require.Equal(t, int(protoNav.SegmentIndex), nav.segmentIndex)
	require.Equal(t, protoNav.SegmentDistance, nav.segmentDistance)
	require.Equal(t, protoNav.CurrentDistance, nav.currentDistance)
	require.Equal(t, int(protoNav.OfflineIndex), nav.offlineIndex)
	require.Equal(t, protoNav.Point.Lat, nav.point.X)
	require.Equal(t, protoNav.Point.Lon, nav.point.Y)
	require.NotNil(t, nav.elevation)
	require.Equal(t, int(protoNav.OfflineMin), nav.offline.Min())
	require.Equal(t, int(protoNav.OfflineMax), nav.offline.Max())
	require.Equal(t, protoNav.TotalDistance, nav.totalDist)
	require.Equal(t, protoNav.SkipOffline, nav.skipOffline)
}
