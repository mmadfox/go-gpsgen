package geojson

import (
	"math"
	"os"
	"testing"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func file(filename string) []byte {
	data, err := os.ReadFile("./testdata/" + filename + ".geojson")
	if err != nil {
		panic(err)
	}
	return data
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		data    func() []*navigator.Route
		assert  func(data []byte)
		wantErr bool
	}{
		{
			name:    "should return error when no routes",
			data:    func() []*navigator.Route { return nil },
			wantErr: true,
		},
		{
			name: "should encode route props",
			data: func() []*navigator.Route {
				route := navigator.RouteFromTracks(track4Line)
				route.Props().Set("foo", "bar").
					Set("speed", 1.2).
					Set("isActive", true)
				return []*navigator.Route{route}
			},
			assert: func(data []byte) {
				require.NotEmpty(t, data)

				fc, err := ParseFeatureCollection(data)
				require.NoError(t, err)
				require.Len(t, fc.Features, 1)
				props := properties.Properties(fc.Features[0].Properties)

				color, ok := props.String("color")
				require.True(t, ok)
				require.NotEmpty(t, color)

				distance, ok := props.Float64("distance")
				require.True(t, ok)
				require.NotEmpty(t, distance)

				speed, ok := props.Float64("speed")
				require.True(t, ok)
				require.NotEmpty(t, speed)

				foo, ok := props.String("foo")
				require.True(t, ok)
				require.NotEmpty(t, foo)

				isActive, ok := props.Bool("isActive")
				require.True(t, ok)
				require.True(t, isActive)

				numTracks, ok := props.Float64("numTracks")
				require.True(t, ok)
				require.NotEmpty(t, numTracks)

				routeID, ok := props.String("routeID")
				require.True(t, ok)
				require.NotEmpty(t, routeID)

				require.NotEmpty(t, props["tracksInfo"])
			},
		},
		{
			name: "should encode multiple routes into one featureCollection",
			data: func() []*navigator.Route {
				route1 := navigator.RouteFromTracks(track1Poly, track3Line)
				route2 := navigator.RouteFromTracks(track3Line)
				route3 := navigator.RouteFromTracks(track2Poly)
				return []*navigator.Route{
					route1, route2, route3,
				}
			},
			assert: func(data []byte) {
				require.NotEmpty(t, data)
				fc, err := ParseFeatureCollection(data)
				require.NoError(t, err)
				require.Len(t, fc.Features, 3)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.data())
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.assert != nil {
				tt.assert(got)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		assert  func([]*navigator.Route)
		wantErr bool
	}{
		{
			name: "should return error when invalid route geometries",
			args: args{
				data: file("invalid_route_geometry"),
			},
			wantErr: true,
		},
		{
			name: "should return error when invalid tracks info",
			args: args{
				data: file("invalid_track_info"),
			},
			wantErr: true,
		},
		{
			name: "should return error when invalid geojson",
			args: args{
				data: []byte(`{}`),
			},
			wantErr: true,
		},
		{
			name: "should return error when LineString geometry is empty",
			args: args{
				data: file("invalid_line_string"),
			},
			wantErr: true,
		},
		{
			name: "should return error when MultiLineString geometry is empty",
			args: args{
				data: file("invalid_multilinestring"),
			},
			wantErr: true,
		},
		{
			name: "should create new routes when geometry LineString",
			args: args{
				data: file("line_string"),
			},
			assert: func(routes []*navigator.Route) {
				assert.Len(t, routes, 1)
				assert.Equal(t, 1, routes[0].NumTracks())
			},
			wantErr: false,
		},
		{
			name: "should create new routes when geometry MultiPoint",
			args: args{
				data: file("multipoint"),
			},
			assert: func(routes []*navigator.Route) {
				assert.Len(t, routes, 1)
				assert.Equal(t, 1, routes[0].NumTracks())
			},
			wantErr: false,
		},
		{
			name: "should create new routes when geometry GeometryCollection",
			args: args{
				data: file("collection"),
			},
			assert: func(routes []*navigator.Route) {
				assert.Len(t, routes, 1)
				assert.Equal(t, 2, routes[0].NumTracks())
			},
			wantErr: false,
		},
		{
			name: "should create new routes when geometry GeometryCollection",
			args: args{
				data: file("multipolygon"),
			},
			assert: func(routes []*navigator.Route) {
				assert.Len(t, routes, 1)
				assert.Equal(t, 2, routes[0].NumTracks())
			},
			wantErr: false,
		},
		{
			name: "should create new routes when geometry MultiLineString",
			args: args{
				data: file("multiline_string"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 1)
				require.Equal(t, 4, routes[0].NumTracks())
			},
			wantErr: false,
		},
		{
			name: "should restore route from geojson",
			args: args{
				data: file("route_line_string"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 1)
				route := routes[0]
				require.Equal(t, "#2dd049", route.Color())
				require.Equal(t, "9c41a8f1-49e2-46c5-a978-898d40ffada1", route.ID())
				require.Equal(t, 1, route.NumTracks())
				require.Equal(t, 1189, math.Floor(route.Distance()))
				require.Len(t, route.Props(), 3)

				track := route.TrackAt(0)
				require.NotNil(t, track)
				require.Equal(t, "1acad315-6618-4edd-8619-d72a64308564", track.ID())
				require.Equal(t, "#5935dd", track.Color())
				require.Equal(t, 1189.315916641673, math.Floor(track.Distance()))
				require.Equal(t, 2, track.NumSegments())
				require.Len(t, track.Props(), 1)
			},
			wantErr: false,
		},
		{
			name: "should restore multi routes from geojson",
			args: args{
				data: file("routes"),
			},
			assert: func(routes []*navigator.Route) {
				require.NotEmpty(t, routes)
				require.Len(t, routes, 3)

				// route1
				route1 := routes[0]
				require.Equal(t, "5a5a1a02-3882-4e9a-aa04-6ec8aed004e9", routes[0].ID())
				require.Equal(t, "#9d031c", route1.Color())
				require.Equal(t, route1.NumTracks(), 2)
				require.Equal(t, 2652, math.Floor(route1.Distance()))
				require.Equal(t, "#9d031c", route1.Color())
				require.Equal(t, 2, route1.NumTracks())
				route1track1 := route1.TrackAt(0)
				require.Equal(t, "46b275c6-c3a5-4f9f-b6d6-c25ebc077311", route1track1.ID())
				require.Equal(t, "#61aa2c", route1track1.Color())
				require.Equal(t, 1990, math.Floor(route1track1.Distance()))
				require.Equal(t, 4, route1track1.NumSegments())
				require.Len(t, route1track1.Props(), 2)
				bar, ok := route1track1.Props().String("foo")
				require.True(t, ok)
				require.Equal(t, "bar", bar)
				speed, ok := route1track1.Props().Float64("speed")
				require.True(t, ok)
				require.Equal(t, float64(3), speed)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.assert != nil {
				tt.assert(got)
			}
		})
	}
}
