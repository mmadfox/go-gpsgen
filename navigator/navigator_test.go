package navigator

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// 1segment: 300m
	track300m1segment, _ = NewTrack([]geo.LatLonPoint{
		{Lon: 106.49331396675268, Lat: 29.5299004724652},
		{Lon: 106.49523863664103, Lat: 29.532016484207674},
	})
	// 1segment: 301m, 2segment: 697m
	track1km2segment, _ = NewTrack([]geo.LatLonPoint{
		{Lon: 106.46029732125612, Lat: 29.532477955504234},
		{Lon: 106.46229029469634, Lat: 29.534568282538686},
		{Lon: 106.4659759305107, Lat: 29.53996017665672},
	})
	track3km7segments, _ = NewTrack([]geo.LatLonPoint{
		{Lon: 106.48818494999364, Lat: 29.526967168711465},
		{Lon: 106.48818494999364, Lat: 29.53306336155346},
		{Lon: 106.49038855444923, Lat: 29.53586552047419},
		{Lon: 106.49055806248373, Lat: 29.537291150475014},
		{Lon: 106.48007681565423, Lat: 29.53522643843239},
		{Lon: 106.47905976744477, Lat: 29.534439870375238},
		{Lon: 106.47617813084906, Lat: 29.53055610091529},
		{Lon: 106.4755283500491, Lat: 29.529843240954946},
	})
)

func routes(n int) []*Route {
	routes := make([]*Route, 0, n)
	for i := 0; i < n; i++ {
		route := RouteFromTracks(track300m1segment, track1km2segment, track3km7segments)
		routes = append(routes, route)
	}
	return routes
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error when params are invalid",
			args: args{
				opts: []Option{
					WithElevation(1, 10000, 10000, 0),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNavigator_MoveToSegment(t *testing.T) {
	type args struct {
		routeIndex   int
		trackIndex   int
		segmentIndex int
	}
	tests := []struct {
		name    string
		args    args
		arrange func(*Navigator)
		assert  func(*Navigator)
		want    bool
	}{
		{
			name: "should return false when no routes",
			want: false,
		},
		{
			name: "should return false when move to the current segment",
			args: args{
				routeIndex:   0,
				trackIndex:   0,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				n.AddRoute(RouteFromTracks(track300m1segment))
			},
			want: false,
		},
		{
			name: "should return false when routeIndex < 0",
			args: args{
				routeIndex: -1,
			},
			arrange: func(n *Navigator) {
				n.AddRoute(RouteFromTracks(track300m1segment))
			},
			want: false,
		},
		{
			name: "should return false when routeIndex > routes",
			args: args{
				routeIndex: 100,
			},
			arrange: func(n *Navigator) {
				n.AddRoute(RouteFromTracks(track300m1segment))
			},
			want: false,
		},
		{
			name: "should return false when no tracks",
			args: args{
				routeIndex:   1,
				trackIndex:   0,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				route1 := NewRoute().AddTrack(track300m1segment)
				route2 := NewRoute()
				n.AddRoute(route1, route2)
			},
			want: false,
		},
		{
			name: "should return false when trackIndex < 0",
			args: args{
				routeIndex:   0,
				trackIndex:   -1,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				route1 := NewRoute().AddTrack(track300m1segment)
				n.AddRoute(route1)
			},
			want: false,
		},
		{
			name: "should return false when trackIndex > tracks",
			args: args{
				routeIndex:   0,
				trackIndex:   100,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				route1 := NewRoute().AddTrack(track300m1segment)
				n.AddRoute(route1)
			},
			want: false,
		},
		{
			name: "should return false when segmentIndex < 0",
			args: args{
				routeIndex:   0,
				trackIndex:   0,
				segmentIndex: -1,
			},
			arrange: func(n *Navigator) {
				route1 := NewRoute().AddTrack(track300m1segment)
				n.AddRoute(route1)
			},
			want: false,
		},
		{
			name: "should return false when segmentIndex > segments",
			args: args{
				routeIndex:   0,
				trackIndex:   0,
				segmentIndex: 100,
			},
			arrange: func(n *Navigator) {
				route1 := NewRoute().AddTrack(track300m1segment)
				n.AddRoute(route1)
			},
			want: false,
		},
		{
			name: "should return true when segment is switched",
			args: args{
				routeIndex:   0,
				trackIndex:   1,
				segmentIndex: 1,
			},
			arrange: func(n *Navigator) {
				newDistInMeters := float64(600)
				route1 := NewRoute().
					AddTrack(track300m1segment).
					AddTrack(track1km2segment)
				n.AddRoute(route1)
				n.DestinationTo(newDistInMeters)
				assert.Equal(t, newDistInMeters, n.CurrentDistance())
				assert.Equal(t, newDistInMeters, n.CurrentRouteDistance())
				assert.NotZero(t, n.CurrentTrackDistance())
				assert.NotZero(t, n.CurrentSegmentDistance())
			},
			assert: func(n *Navigator) {
				require.NotZero(t, n.CurrentDistance())
				require.NotZero(t, n.CurrentTrackDistance())
				require.Zero(t, n.CurrentSegmentDistance())
			},
			want: true,
		},
		{
			name: "should return true when track is switched",
			args: args{
				routeIndex:   0,
				trackIndex:   1,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				newDistInMeters := float64(10)
				route1 := NewRoute().
					AddTrack(track1km2segment).
					AddTrack(track1km2segment)
				n.AddRoute(route1)
				n.DestinationTo(newDistInMeters)
			},
			assert: func(n *Navigator) {
				expectedDist := float64(10)
				require.Equal(t, expectedDist, n.CurrentDistance())
				require.Equal(t, expectedDist, n.CurrentRouteDistance())
				require.Zero(t, n.CurrentTrackDistance())
				require.Zero(t, n.CurrentSegmentDistance())
			},
			want: true,
		},
		{
			name: "should return true when route is switched",
			args: args{
				routeIndex:   1,
				trackIndex:   0,
				segmentIndex: 0,
			},
			arrange: func(n *Navigator) {
				newDistInMeters := float64(10)
				route1 := NewRoute().AddTrack(track1km2segment)
				route2 := NewRoute().AddTrack(track1km2segment)
				n.AddRoute(route1, route2)
				n.DestinationTo(newDistInMeters)
			},
			assert: func(n *Navigator) {
				expectedDist := float64(10)
				require.Equal(t, expectedDist, n.CurrentDistance())
				require.Zero(t, n.CurrentRouteDistance())
				require.Zero(t, n.CurrentTrackDistance())
				require.Zero(t, n.CurrentSegmentDistance())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := New(SkipOfflineMode())
			require.NoError(t, err)
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if got := n.MoveToSegment(tt.args.routeIndex, tt.args.trackIndex, tt.args.segmentIndex); got != tt.want {
				t.Errorf("Navigator.ToSegment() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_ResetRoutes(t *testing.T) {
	n, _ := New()
	expectedNumRoutes := 10
	newDistInMeters := float64(100)
	require.False(t, n.ResetRoutes())

	n.AddRoute(routes(expectedNumRoutes)...)
	require.Equal(t, expectedNumRoutes, n.NumRoutes())
	n.DestinationTo(newDistInMeters)

	require.True(t, n.ResetRoutes())
	require.Zero(t, n.NumRoutes())
	require.Zero(t, n.CurrentDistance())
	require.Zero(t, n.CurrentRouteDistance())
	require.Zero(t, n.CurrentTrackDistance())
	require.Zero(t, n.CurrentSegmentDistance())
	require.Zero(t, n.RouteIndex())
	require.Zero(t, n.TrackIndex())
	require.Zero(t, n.SegmentIndex())
}

func TestNavigator_AddRoute(t *testing.T) {
	type args struct {
		routes func() []*Route
	}
	tests := []struct {
		name    string
		args    args
		assert  func(*Navigator)
		wantErr bool
	}{
		{
			name: "should return error when are no routes",
			args: args{
				routes: func() []*Route {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "should return nil when each route is nil",
			args: args{
				routes: func() []*Route {
					return []*Route{nil, nil, nil}
				},
			},
			assert: func(n *Navigator) {
				assert.Equal(t, 0, n.NumRoutes())
			},
			wantErr: false,
		},
		{
			name: "routes added successfully",
			args: args{
				routes: func() []*Route {
					return routes(10)
				},
			},
			assert: func(n *Navigator) {
				assert.Equal(t, 10, n.NumRoutes())
			},
			wantErr: false,
		},
		{
			name: "should not add route when already exists",
			args: args{
				routes: func() []*Route {
					routes := routes(1)
					return []*Route{routes[0], routes[0]}
				},
			},
			assert: func(n *Navigator) {
				require.Equal(t, 1, n.NumRoutes())
				require.Equal(t, [3]int{1, 3, 0}, n.Sum())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := New()
			require.NoError(t, err)
			if err := n.AddRoute(tt.args.routes()...); (err != nil) != tt.wantErr {
				t.Errorf("Navigator.AddRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_RouteByID(t *testing.T) {
	type args struct {
		routeID string
	}

	expectedRoute := routes(1)[0]

	tests := []struct {
		name    string
		args    args
		arrange func(*Navigator)
		want    *Route
		wantErr bool
	}{
		{
			name:    "should return error when routeID is empty",
			wantErr: true,
		},
		{
			name: "should return error when route not found",
			args: args{
				routeID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "route found successfully",
			args: args{
				routeID: expectedRoute.ID(),
			},
			arrange: func(n *Navigator) {
				n.AddRoute(expectedRoute)
			},
			want:    expectedRoute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			if tt.arrange != nil {
				tt.arrange(n)
			}
			got, err := n.RouteByID(tt.args.routeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Navigator.RouteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Navigator.RouteByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigator_Routes(t *testing.T) {
	n, err := New()
	require.NoError(t, err)

	route := RouteFromTracks(track1km2segment, track3km7segments)
	n.AddRoute(route)

	require.Len(t, n.Routes(), 1)
}

func TestNavigator_EachRoute(t *testing.T) {
	routes := routes(10)
	n, _ := New()
	n.AddRoute(routes...)

	actualLoop := 0
	n.EachRoute(func(n int, r *Route) bool {
		actualLoop++
		return false
	})
	require.Equal(t, 1, actualLoop)

	actualLoop = 0
	n.EachRoute(func(n int, r *Route) bool {
		if n > 1 {
			return false
		}
		actualLoop++
		return true
	})
	require.Equal(t, 2, actualLoop)

	n.ResetRoutes()
	actualLoop = 0
	n.EachRoute(func(n int, r *Route) bool {
		actualLoop++
		return true
	})
	require.Zero(t, actualLoop)
}

func TestNavigator_RemoveRoute(t *testing.T) {
	expectedRoutes := routes(10)

	type args struct {
		routeID string
	}
	tests := []struct {
		name   string
		args   args
		assert func(*Navigator)
		wantOk bool
	}{
		{
			name:   "should return false when routeID is empty",
			wantOk: false,
		},
		{
			name: "route deleted successfully",
			args: args{
				routeID: expectedRoutes[0].ID(),
			},
			assert: func(n *Navigator) {
				assert.Equal(t, len(expectedRoutes)-1, n.NumRoutes())
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			n.AddRoute(expectedRoutes...)
			if gotOk := n.RemoveRoute(tt.args.routeID); gotOk != tt.wantOk {
				t.Errorf("Navigator.RemoveRoute() = %v, want %v", gotOk, tt.wantOk)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_RemoveRoutes(t *testing.T) {
	n, _ := New()
	routes := routes(10)
	n.AddRoute(routes...)

	for i := 0; i < len(routes); i++ {
		route := routes[i]
		assert.True(t, n.RemoveRoute(route.ID()))
		assert.Zero(t, n.RouteIndex())
		assert.Zero(t, n.TrackIndex())
		assert.Zero(t, n.SegmentIndex())
		assert.Zero(t, n.CurrentDistance())
		assert.Zero(t, n.CurrentRouteDistance())
		assert.Zero(t, n.CurrentTrackDistance())
		assert.Zero(t, n.CurrentSegmentDistance())
	}

	assert.Zero(t, n.Distance())
}

func TestNavigator_RemoveTrack(t *testing.T) {
	route := RouteFromTracks(track3km7segments)
	routeID := route.ID()
	trackID := route.TrackAt(0).ID()
	numTracks := route.NumTracks()
	require.NotZero(t, numTracks)

	type args struct {
		routeID string
		trackID string
	}
	tests := []struct {
		name   string
		args   args
		assert func(*Navigator)
		want   bool
	}{
		{
			name: "should return false when routeID is empty",
			want: false,
		},
		{
			name: "should return false when trackID is empty",
			want: false,
		},
		{
			name: "should return false if the route is not found",
			args: args{
				routeID: "someid",
				trackID: "someid",
			},
			want: false,
		},
		{
			name: "track deleted successfully",
			args: args{
				routeID: routeID,
				trackID: trackID,
			},
			assert: func(n *Navigator) {
				assert.Equal(t, numTracks-1, route.NumTracks())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			n.AddRoute(route)
			if got := n.RemoveTrack(tt.args.routeID, tt.args.trackID); got != tt.want {
				t.Errorf("Navigator.RemoveTrack() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_NextLocation(t *testing.T) {
	type args struct {
		tick        float64 // seconds
		speed       float64 // m/s
		skipOffline bool
	}
	tests := []struct {
		name    string
		args    args
		arrange func(*Navigator)
		assert  func(*Navigator)
		wantOk  bool
	}{
		{
			name:   "should return false if there are no routes",
			wantOk: false,
		},
		{
			name: "should return false when navigator is offline",
			args: args{
				tick:  float64(1),
				speed: float64(3),
			},
			arrange: func(n *Navigator) {
				route := RouteFromTracks(track300m1segment)
				n.AddRoute(route)
				n.ToOffline()
			},
			assert: func(n *Navigator) {
				assert.True(t, n.IsOffline())
				assert.NotZero(t, n.OfflineDuration())
			},
			wantOk: false,
		},
		{
			name: "should walk successfully 2 meters in one segment",
			args: args{
				tick:  float64(1),
				speed: float64(2),
			},
			arrange: func(n *Navigator) {
				route := RouteFromTracks(track300m1segment)
				n.AddRoute(route)
			},
			assert: func(n *Navigator) {
				newDist := float64(2)
				assert.Equal(t, newDist, n.CurrentDistance())
				assert.Equal(t, newDist, n.CurrentRouteDistance())
				assert.Equal(t, newDist, n.CurrentTrackDistance())
				assert.Equal(t, newDist, n.CurrentSegmentDistance())
				assert.Equal(t, 0, n.RouteIndex())
				assert.Equal(t, 0, n.TrackIndex())
				assert.Equal(t, 0, n.SegmentIndex())
				assert.NotZero(t, n.Distance())
				assert.NotZero(t, n.RouteDistance())
				assert.NotZero(t, n.TrackDistance())
				assert.NotZero(t, n.SegmentDistance())
				assert.False(t, n.IsOffline())
				assert.False(t, n.IsFinish())
				assert.Zero(t, n.OfflineDuration())

				segment := n.RouteAt(0).TrackAt(0).SegmentAt(0)
				assert.Equal(t, segment.Bearing(), n.CurrentBearing())
				assert.NotZero(t, n.Location().Lat)
				assert.NotZero(t, n.Location().Lon)
			},
			wantOk: true,
		},
		{
			name: "should successfully switch next segment without offline mode",
			args: args{
				tick:  1,
				speed: 100,
			},
			arrange: func(n *Navigator) {
				route := RouteFromTracks(track300m1segment, track1km2segment)
				n.AddRoute(route)

				n.DestinationTo(600) // move to 600 meters

				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 1, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
			},
			assert: func(n *Navigator) {
				expectedDist := float64(700)
				require.Equal(t, expectedDist, n.CurrentDistance())
				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 1, n.TrackIndex())
				require.Equal(t, 1, n.SegmentIndex())
				require.False(t, n.IsOffline())
				require.False(t, n.IsFinish())
			},
			wantOk: true,
		},
		{
			name: "should successfully switch next segment with offline mode",
			args: args{
				tick:  1,
				speed: 10,
			},
			arrange: func(n *Navigator) {
				route := RouteFromTracks(track300m1segment, track1km2segment)
				n.AddRoute(route)

				n.DestinationTo(300) // move to 300 meters => route:0,track:0,segment:0

				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 0, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
			},
			assert: func(n *Navigator) {
				// move to 10 meters => route:0,track:1,segment:0
				// when switching between tracks or unrelated segments, offline mode is activated
				expectedDist := float64(300)
				require.Equal(t, expectedDist, n.CurrentDistance())
				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 1, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
				require.True(t, n.IsOffline())
				require.False(t, n.IsFinish())
			},
			wantOk: false,
		},
		{
			name: "should successfully switch next route with offline mode",
			args: args{
				tick:  1,
				speed: 10,
			},
			arrange: func(n *Navigator) {
				route1 := RouteFromTracks(track300m1segment)
				route2 := RouteFromTracks(track300m1segment)

				n.AddRoute(route1, route2)

				n.DestinationTo(299) // move to 299 meters

				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 0, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
			},
			assert: func(n *Navigator) {
				expectedDist := float64(299)
				require.Equal(t, expectedDist, n.CurrentDistance())
				require.Equal(t, 1, n.RouteIndex())
				require.Equal(t, 0, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
				require.True(t, n.IsOffline())
				require.False(t, n.IsFinish())
			},
			wantOk: false,
		},
		{
			name: "should successfully switch start route with offline mode",
			args: args{
				tick:        1,
				speed:       10,
				skipOffline: true,
			},
			arrange: func(n *Navigator) {
				route := RouteFromTracks(track1km2segment)

				n.AddRoute(route)

				n.DestinationTo(998)

				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 0, n.TrackIndex())
				require.Equal(t, 1, n.SegmentIndex())
			},
			assert: func(n *Navigator) {
				expectedStartDist := float64(0.1)
				require.Equal(t, expectedStartDist, n.CurrentDistance())
				require.Equal(t, 0, n.RouteIndex())
				require.Equal(t, 0, n.TrackIndex())
				require.Equal(t, 0, n.SegmentIndex())
				require.False(t, n.IsOffline())
				require.True(t, n.IsFinish())
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]Option, 0)
			if tt.args.skipOffline {
				opts = append(opts, SkipOfflineMode())
			} else {
				opts = append(opts, WithOffline(3, 10))
			}
			n, _ := New(opts...)
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if gotOk := n.NextLocation(tt.args.tick, tt.args.speed); gotOk != tt.wantOk {
				t.Errorf("Navigator.NextLocation() = %v, want %v", gotOk, tt.wantOk)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_Update(t *testing.T) {
	route := RouteFromTracks(track1km2segment)
	n, _ := New()
	n.AddRoute(route)
	n.NextLocation(1, 100)
	state := new(proto.Device)
	n.Update(state)

	require.Equal(t, n.Distance(), state.Distance.Distance)
	require.Equal(t, n.CurrentDistance(), state.Distance.CurrentDistance)
	require.Equal(t, n.RouteDistance(), state.Distance.RouteDistance)
	require.Equal(t, n.CurrentRouteDistance(), state.Distance.CurrentRouteDistance)
	require.Equal(t, n.TrackDistance(), state.Distance.TrackDistance)
	require.Equal(t, n.CurrentTrackDistance(), state.Distance.CurrentTrackDistance)
	require.Equal(t, n.SegmentDistance(), state.Distance.SegmentDistance)
	require.Equal(t, n.CurrentSegmentDistance(), state.Distance.CurrentSegmentDistance)

	require.Equal(t, int64(n.RouteIndex()), state.Navigator.CurrentRouteIndex)
	require.Equal(t, int64(n.TrackIndex()), state.Navigator.CurrentTrackIndex)
	require.Equal(t, int64(n.SegmentIndex()), state.Navigator.CurrentSegmentIndex)
	require.NotEmpty(t, state.Navigator.CurrentRouteId)
	require.NotEmpty(t, state.Navigator.CurrentTrackId)

	require.Equal(t, n.CurrentBearing(), state.Location.Bearing)
	require.Equal(t, n.Elevation(), state.Location.Elevation)
	require.Equal(t, n.Location().Lon, state.Location.Lon)
	require.Equal(t, n.Location().Lat, state.Location.Lat)
	require.NotZero(t, state.Location.Utm.CentralMeridian)
	require.NotZero(t, state.Location.Utm.Easting)
	require.NotZero(t, state.Location.Utm.Hemisphere)
	require.NotZero(t, state.Location.Utm.LatZone)
	require.NotZero(t, state.Location.Utm.LongZone)
	require.NotZero(t, state.Location.Utm.Northing)
	require.NotZero(t, state.Location.Utm.Srid)

	require.NotZero(t, state.Location.LonDms.Degrees)
	require.NotZero(t, state.Location.LonDms.Direction)
	require.NotZero(t, state.Location.LonDms.Minutes)
	require.NotZero(t, state.Location.LonDms.Seconds)

	require.NotZero(t, state.Location.LatDms.Degrees)
	require.NotZero(t, state.Location.LatDms.Direction)
	require.NotZero(t, state.Location.LatDms.Minutes)
	require.NotZero(t, state.Location.LatDms.Seconds)

	require.Equal(t, n.IsOffline(), state.IsOffline)
	require.Equal(t, int64(n.OfflineDuration()), state.OfflineDuration)
}

func TestNavigator_CurrentBearing(t *testing.T) {
	n, _ := New()
	require.Zero(t, n.CurrentBearing())

	route := RouteFromTracks(track1km2segment)
	n.AddRoute(route)
	n.MoveToSegment(0, 0, 1)
	expectedBearing := route.TrackAt(0).SegmentAt(1).Bearing()
	require.Equal(t, expectedBearing, n.CurrentBearing())
}

func TestNavigator_NextRoute(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(*Navigator)
		assert  func(*Navigator)
		want    bool
	}{
		{
			name: "should return false when last route",
			arrange: func(n *Navigator) {
				n.AddRoute(routes(3)...)
				n.MoveToRoute(2)
			},
			assert: func(n *Navigator) {
				require.Equal(t, 2, n.RouteIndex())
			},
			want: false,
		},
		{
			name: "should switch to next route",
			arrange: func(n *Navigator) {
				n.AddRoute(routes(3)...)
				require.Equal(t, 0, n.RouteIndex())
			},
			assert: func(n *Navigator) {
				require.Equal(t, 1, n.RouteIndex())
				require.Zero(t, n.TrackIndex())
				require.Zero(t, n.SegmentIndex())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if got := n.NextRoute(); got != tt.want {
				t.Errorf("Navigator.NextRoute() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_PrevRoute(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(*Navigator)
		assert  func(*Navigator)
		want    bool
	}{
		{
			name: "should return false when first route",
			arrange: func(n *Navigator) {
				n.AddRoute(routes(3)...)
			},
			assert: func(n *Navigator) {
				require.Equal(t, 0, n.RouteIndex())
			},
			want: false,
		},
		{
			name: "should switch to prev route",
			arrange: func(n *Navigator) {
				n.AddRoute(routes(3)...)
				n.MoveToRoute(2)
				require.Equal(t, 2, n.RouteIndex())
			},
			assert: func(n *Navigator) {
				require.Equal(t, 1, n.RouteIndex())
				require.Zero(t, n.TrackIndex())
				require.Zero(t, n.SegmentIndex())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if got := n.PrevRoute(); got != tt.want {
				t.Errorf("Navigator.PrevRoute() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_MoveToRouteByID(t *testing.T) {
	type args struct {
		routeID string
	}
	routes := routes(10)
	tests := []struct {
		name    string
		args    args
		arrange func(*Navigator)
		assert  func(*Navigator)
		want    bool
	}{
		{
			name: "should return false when routeID is empty",
			args: args{
				routeID: "",
			},
			want: false,
		},
		{
			name: "should return false when route not found",
			args: args{
				routeID: uuid.NewString(),
			},
			want: false,
		},
		{
			name: "successfully moved to route",
			args: args{
				routeID: routes[5].ID(),
			},
			arrange: func(n *Navigator) {
				n.AddRoute(routes...)
			},
			assert: func(n *Navigator) {
				require.Equal(t, routes[5].ID(), n.CurrentRoute().ID())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if got := n.MoveToRouteByID(tt.args.routeID); got != tt.want {
				t.Errorf("Navigator.MoveToRouteByID() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_MoveToTrackByID(t *testing.T) {
	type args struct {
		routeID string
		trackID string
	}

	routes := routes(10)

	tests := []struct {
		name    string
		args    args
		arrange func(*Navigator)
		assert  func(*Navigator)
		want    bool
	}{
		{
			name: "should return false when routeID is empty",
			args: args{
				routeID: "",
			},
			want: false,
		},
		{
			name: "should return false when trackID is empty",
			args: args{
				routeID: routes[0].ID(),
				trackID: "",
			},
			arrange: func(n *Navigator) {
				n.AddRoute(routes...)
			},
			want: false,
		},
		{
			name: "successfully moved to route",
			args: args{
				routeID: routes[8].ID(),
				trackID: routes[1].TrackAt(1).ID(),
			},
			arrange: func(n *Navigator) {
				n.AddRoute(routes...)
			},
			assert: func(n *Navigator) {
				require.Equal(t, routes[8].ID(), n.CurrentRoute().ID())
				require.Equal(t, routes[8].TrackAt(1).ID(), n.CurrentTrack().ID())
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, _ := New()
			if tt.arrange != nil {
				tt.arrange(n)
			}
			if got := n.MoveToTrackByID(tt.args.routeID, tt.args.trackID); got != tt.want {
				t.Errorf("Navigator.MoveToTrackByID() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(n)
			}
		})
	}
}

func TestNavigator_SelectAt(t *testing.T) {
	n, _ := New()
	n.AddRoute(routes(10)...)

	require.Nil(t, n.RouteAt(10))
	require.Nil(t, n.RouteAt(-1))
	require.Nil(t, n.RouteAt(1).TrackAt(11))
	require.Nil(t, n.RouteAt(1).TrackAt(-1))
	require.True(t, n.RouteAt(1).TrackAt(1).SegmentAt(-1).IsEmpty())
	require.True(t, n.RouteAt(1).TrackAt(1).SegmentAt(111).IsEmpty())

	require.NotNil(t, n.RouteAt(0))
	require.NotNil(t, n.RouteAt(9))
	require.NotNil(t, n.RouteAt(0).TrackAt(1))
	require.NotNil(t, n.RouteAt(2).TrackAt(1))
	require.False(t, n.RouteAt(1).TrackAt(1).SegmentAt(0).IsEmpty())
	require.False(t, n.RouteAt(1).TrackAt(1).SegmentAt(1).IsEmpty())
}

func TestNavigator_WithoutRoutes(t *testing.T) {
	n, _ := New()
	require.Nil(t, n.CurrentRoute())
	require.Nil(t, n.CurrentTrack())
	require.True(t, n.CurrentSegment().IsEmpty())
	require.Zero(t, n.RouteDistance())
	require.Zero(t, n.TrackDistance())
	require.Zero(t, n.SegmentDistance())
}

func TestNavigator_Snapshot(t *testing.T) {
	n, _ := New()
	n.AddRoute(RouteFromTracks(track300m1segment))
	n.AddRoute(RouteFromTracks(track300m1segment))
	n.DestinationTo(500)
	snapshot := n.Snapshot()
	require.NotNil(t, snapshot)
	require.Len(t, snapshot.Routes, n.NumRoutes())
	require.Equal(t, int64(n.RouteIndex()), snapshot.RouteIndex)
	require.Equal(t, int64(n.TrackIndex()), snapshot.TrackIndex)
	require.Equal(t, int64(n.SegmentIndex()), snapshot.SegmentIndex)
	require.Equal(t, n.CurrentSegmentDistance(), snapshot.CurrentSegmentDistance)
	require.Equal(t, n.CurrentRouteDistance(), snapshot.CurrentRouteDistance)
	require.Equal(t, n.CurrentTrackDistance(), snapshot.CurrentTrackDistance)
	require.Equal(t, n.CurrentDistance(), snapshot.CurrentDistance)
	require.Equal(t, int64(n.offlineIndex), snapshot.OfflineIndex)
	require.Equal(t, n.point.Lon, snapshot.Point.Lon)
	require.Equal(t, n.point.Lat, snapshot.Point.Lat)
	require.NotNil(t, snapshot.Elevation)
	require.Equal(t, int64(n.offline.Min()), snapshot.OfflineMin)
	require.Equal(t, int64(n.offline.Max()), snapshot.OfflineMax)
	require.Equal(t, n.Distance(), snapshot.Distance)
	require.Equal(t, n.skipOffline, snapshot.SkipOffline)
}

func TestNavigator_FromSnapshot(t *testing.T) {
	n1, _ := New()
	n1.AddRoute(RouteFromTracks(track300m1segment))
	n1.AddRoute(RouteFromTracks(track300m1segment))
	n1.DestinationTo(500)
	snapshot := n1.Snapshot()
	n := new(Navigator)
	n.FromSnapshot(snapshot)
	require.Len(t, snapshot.Routes, n.NumRoutes())
	require.Equal(t, snapshot.RouteIndex, int64(n.RouteIndex()))
	require.Equal(t, snapshot.TrackIndex, int64(n.TrackIndex()))
	require.Equal(t, snapshot.SegmentIndex, int64(n.SegmentIndex()))
	require.Equal(t, snapshot.CurrentSegmentDistance, n.CurrentSegmentDistance())
	require.Equal(t, snapshot.CurrentRouteDistance, n.CurrentRouteDistance())
	require.Equal(t, snapshot.CurrentTrackDistance, n.CurrentTrackDistance())
	require.Equal(t, snapshot.CurrentDistance, n.CurrentDistance())
	require.Equal(t, snapshot.OfflineIndex, int64(n.offlineIndex))
	require.Equal(t, snapshot.Point.Lon, n.point.Lon)
	require.Equal(t, snapshot.Point.Lat, n.point.Lat)
	require.Equal(t, snapshot.OfflineMin, int64(n.offline.Min()))
	require.Equal(t, snapshot.OfflineMax, int64(n.offline.Max()))
	require.Equal(t, snapshot.Distance, n.Distance())
	require.Equal(t, snapshot.SkipOffline, n.skipOffline)
}

func TestNavigator_NextElevation(t *testing.T) {
	n1, _ := New(WithElevation(1, 10, 8, 0))
	for i := 0; i < 10; i++ {
		v := float64(i) * 0.1
		n1.NextElevation(v)
		val := n1.Elevation()
		if val < 1 || val > 10 {
			t.Fatalf("navigator.Elevation() => %f, want > 1 and < 10", val)
		}
	}
}
