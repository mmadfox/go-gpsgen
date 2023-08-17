package geojson

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/stretchr/testify/require"
)

var (
	track1Poly, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.45792276464522, Lat: 29.533628561797286},
		{Lon: 106.45792276464522, Lat: 29.529399791830542},
		{Lon: 106.46334808393323, Lat: 29.529399791830542},
		{Lon: 106.46334808393323, Lat: 29.533628561797286},
		{Lon: 106.45792276464522, Lat: 29.533628561797286},
	})

	track2Poly, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.46572166112287, Lat: 29.533628561797286},
		{Lon: 106.46572166112287, Lat: 29.52910475476996},
		{Lon: 106.47148606286629, Lat: 29.52910475476996},
		{Lon: 106.47148606286629, Lat: 29.533628561797286},
		{Lon: 106.46572166112287, Lat: 29.533628561797286},
	})
	track3Line, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.45688223543743, Lat: 29.52481100757835},
		{Lon: 106.46224152558176, Lat: 29.528509455345258},
	})
	track4Line, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.46609599041324, Lat: 29.528233799305895},
		{Lon: 106.47185128721964, Lat: 29.526671734226028},
		{Lon: 106.46691440417999, Lat: 29.523317807121842},
	})
)

func TestToFeatureCollection(t *testing.T) {
	route1 := navigator.RouteFromTracks(track1Poly, track3Line)
	route2 := navigator.RouteFromTracks(track2Poly)
	route3 := navigator.NewRoute() // empty route
	route1.Props().Set("foo", "foo")
	expectedRoutes := 2

	fc := ToFeatureCollection([]*navigator.Route{route1, route2, route3})
	require.NotNil(t, fc)

	require.Equal(t, expectedRoutes, len(fc.Features))
	feature1 := fc.Features[0]
	feature2 := fc.Features[1]
	require.Equal(t, route1.ID(), feature1.Properties["routeID"])
	require.Equal(t, route1.Color(), feature1.Properties["color"])
	require.Equal(t, route1.NumTracks(), feature1.Properties["numTracks"])
	require.Equal(t, route1.Distance(), feature1.Properties["distance"])
	require.Len(t, feature1.Properties["tracksInfo"], route1.NumTracks())
	require.NotEmpty(t, feature1.Properties["foo"])
	require.Equal(t, "GeometryCollection", feature1.Geometry.GeoJSONType())

	require.Equal(t, route2.ID(), feature2.Properties["routeID"])
	require.Equal(t, route2.Color(), feature2.Properties["color"])
	require.Equal(t, route2.NumTracks(), feature2.Properties["numTracks"])
	require.Equal(t, route2.Distance(), feature2.Properties["distance"])
	require.Len(t, feature2.Properties["tracksInfo"], route2.NumTracks())
	require.Equal(t, "Polygon", feature2.Geometry.GeoJSONType())
}
