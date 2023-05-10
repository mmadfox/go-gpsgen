package gpx

import (
	"io/ioutil"
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/stretchr/testify/require"
)

func TestGPXEncode(t *testing.T) {
	route, err := route.China3()
	require.NoError(t, err)
	data, err := Encode([]*navigator.Route{route})
	require.NoError(t, err)
	require.NotZero(t, data)

	content, err := ioutil.ReadFile("testdata/gpx.golden")
	require.NoError(t, err)

	require.Equal(t, string(content), string(data))
}

func TestGPXDecodeRoutes(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/routes.gpx")
	require.NoError(t, err)

	routes, err := Decode(content)
	require.NoError(t, err)
	require.Equal(t, 1, len(routes))
	require.Equal(t, 1, routes[0].NumTracks())
	require.Equal(t, 131, int(routes[0].TotalDistance()))
}

func TestGPXDecodeTracks(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/tracks.gpx")
	require.NoError(t, err)

	routes, err := Decode(content)
	require.NoError(t, err)
	require.Equal(t, 2, len(routes))
	require.Equal(t, 2, routes[0].NumTracks())
	require.Equal(t, 2, routes[1].NumTracks())
}

func TestGPXDecodeWaypoint(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/waypoint.gpx")
	require.NoError(t, err)

	routes, err := Decode(content)
	require.NoError(t, err)
	require.Equal(t, 1, len(routes))
	require.Equal(t, 1, routes[0].NumTracks())
}
