package gpsgen

import (
	"fmt"
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/types"
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
	data, err := EncodeGeoJSONRoutes(routes)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	routes2, err := DecodeGeoJSONRoutes(data)
	require.NoError(t, err)
	assertRotues(t, routes, routes2)
}

func TestGPX(t *testing.T) {
	routes := testRoutes()
	data, err := EncodeGPXRoutes(routes)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	routes2, err := DecodeGPXRoutes(data)
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

func TestEncodeDecodeTracker(t *testing.T) {
	data, err := EncodeTracker(nil)
	require.Empty(t, data)
	require.NoError(t, err)

	trk := NewDogTracker()
	trk.AddRoute(testRoutes()...)
	trk.Update()

	data, err = EncodeTracker(trk)
	require.NoError(t, err)
	require.NotZero(t, data)

	trk2, err := DecodeTracker(data)
	require.NoError(t, err)
	require.Equal(t, trk, trk2)
}

func TestDecodeTrackerNegative(t *testing.T) {
	trk, err := DecodeTracker(nil)
	require.Nil(t, trk)
	require.Error(t, err)

	trk, err = DecodeTracker([]byte("somedata"))
	require.Nil(t, trk)
	require.Error(t, err)
}

func TestEncodeDecodeSensors(t *testing.T) {
	sensors := make([]*types.Sensor, 0)
	for i := 0; i < 10; i++ {
		sensorName := fmt.Sprintf("s-%d", i)
		sensor, err := types.NewSensor(sensorName, 0, float64(i), 8, types.WithRandom)
		require.NoError(t, err)
		sensors = append(sensors, sensor)
	}

	data, err := EncodeSensors(nil)
	require.Empty(t, data)
	require.Error(t, err)

	data, err = EncodeSensors(sensors)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	sensors2, err := DecodeSensors(data)
	require.NoError(t, err)
	require.Equal(t, sensors, sensors2)
}

func TestDecodeSensorsNegative(t *testing.T) {
	sensors, err := DecodeSensors(nil)
	require.Nil(t, sensors)
	require.Error(t, err)

	sensors, err = DecodeSensors([]byte("somedata"))
	require.Nil(t, sensors)
	require.Error(t, err)
}

func TestEncodeDecodeRoutes(t *testing.T) {
	routes := testRoutes()

	data, err := EncodeRoutes(nil)
	require.Empty(t, data)
	require.Error(t, err)

	data, err = EncodeRoutes(routes)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	routes2, err := DecodeRoutes(data)
	require.NoError(t, err)
	require.Equal(t, routes, routes2)
}

func TestDecodeRoutesNegative(t *testing.T) {
	routes, err := DecodeRoutes(nil)
	require.Nil(t, routes)
	require.Error(t, err)

	routes, err = DecodeRoutes([]byte("somedata"))
	require.Nil(t, routes)
	require.Error(t, err)
}
