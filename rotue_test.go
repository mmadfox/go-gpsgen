package gpsgen

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/stretchr/testify/require"
)

func TestRandomRoute(t *testing.T) {
	type args struct {
		lat, lon float64
		numTrack int
		level    int
	}
	tests := []struct {
		name   string
		args   args
		assert func(d *navigator.Route)
	}{
		{
			name: "should return route with one track when numTracks is negative",
			args: args{
				numTrack: -1,
				lon:      101.66733912613665,
				lat:      35.67896024639478,
				level:    RouteLevelXS - 1,
			},
			assert: func(r *navigator.Route) {
				require.Equal(t, 1, r.NumTracks())
				require.Equal(t, 16, r.TrackAt(0).NumSegments())
			},
		},
		{
			name: "should return route with three tracks",
			args: args{
				numTrack: 3,
				lon:      101.66733912613665,
				lat:      35.67896024639478,
				level:    RouteLevelXXL + 1,
			},
			assert: func(r *navigator.Route) {
				require.Equal(t, 3, r.NumTracks())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := RandomRoute(tt.args.lon, tt.args.lat, tt.args.numTrack, tt.args.level)
			if tt.assert != nil {
				tt.assert(route)
			}
		})
	}
}

func TestRandomRouteForNewYork(t *testing.T) {
	expectedNumTracks := 3

	route := RandomRouteForNewYork()
	require.NotNil(t, route)
	require.Equal(t, expectedNumTracks, route.NumTracks())
}

func TestRandomRouteForMoscow(t *testing.T) {
	expectedNumTracks := 3

	route := RandomRouteForMoscow()
	require.NotNil(t, route)
	require.Equal(t, expectedNumTracks, route.NumTracks())
}

func TestRandomRouteForParis(t *testing.T) {
	expectedNumTracks := 3

	route := RandomRouteForParis()
	require.NotNil(t, route)
	require.Equal(t, expectedNumTracks, route.NumTracks())
}
