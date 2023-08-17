package gpsgen

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/stretchr/testify/require"
)

func testRoutes() []*navigator.Route {
	lon := 104.06827900598023
	lat := 30.665504435503408
	routes := make([]*navigator.Route, 0, 1)
	for i := 0; i < 1; i++ {
		route := RandomRoute(lon, lat, 1, RouteLevelXS)
		routes = append(routes, route)
	}
	return routes
}

func TestGeoJSON(t *testing.T) {
	routes := testRoutes()
	data, err := GeoJSONEncode(routes)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	routes2, err := GeoJSONDecode(data)
	require.NoError(t, err)
	assertRotues(t, routes, routes2)
}

func TestGPX(t *testing.T) {
	routes := testRoutes()
	data, err := GPXEncode(routes)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	routes2, err := GPXDecode(data)
	require.NoError(t, err)
	assertRotues(t, routes, routes2)
}

func assertRotues(t *testing.T, expected, actual []*navigator.Route) {
	require.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		r1 := expected[i]
		r2 := actual[i]
		require.Equal(t, r1.NumTracks(), r2.NumTracks())
		for j := 0; j < r1.NumTracks(); j++ {
			t1 := r1.TrackAt(j)
			t2 := r2.TrackAt(j)
			require.Equal(t, t1.NumSegments(), t2.NumSegments())
		}
	}
}
