package gpx

import (
	"bytes"
	"math"
	"os"
	"testing"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/stretchr/testify/require"
	"github.com/tkrajina/gpxgo/gpx"
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
	track4Line, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.46609599041324, Lat: 29.528233799305895},
		{Lon: 106.47185128721964, Lat: 29.526671734226028},
		{Lon: 106.46691440417999, Lat: 29.523317807121842},
	})
)

func file(filename string) []byte {
	data, err := os.ReadFile("./testdata/" + filename + ".gpx")
	if err != nil {
		panic(err)
	}
	return data
}

func TestEncode(t *testing.T) {
	_, err := Encode(nil)
	require.ErrorIs(t, err, ErrNoRoutes)

	route1 := navigator.RouteFromTracks(track1Poly, track2Poly, track4Line)
	route1.Props().Set("foo", 1).Set("baz", true).Set("ing", 123).Set("val", "val")
	route2 := navigator.RouteFromTracks(track4Line, track2Poly)
	data, err := Encode([]*navigator.Route{route1, route2})
	require.NoError(t, err)
	require.NotEmpty(t, data)

	xmlData, err := gpx.Parse(bytes.NewReader(data))
	require.NoError(t, err)
	require.NotEmpty(t, xmlData.Extensions)
	require.Equal(t, "gpsgen", xmlData.Extensions.Nodes[0].XMLName.Local)
	require.NotEmpty(t, xmlData.Extensions)
	require.Len(t, xmlData.Tracks, 2)

	// routes
	require.Len(t, xmlData.Extensions.Nodes[0].Nodes, 2)
	// routes -> tracks
	require.Len(t, xmlData.Extensions.Nodes[0].Nodes[0].Nodes, 3)
	require.Len(t, xmlData.Extensions.Nodes[0].Nodes[1].Nodes, 2)
	require.Len(t, xmlData.Tracks[0].Segments, 3)
	require.Len(t, xmlData.Tracks[0].Segments[0].Points, track1Poly.NumSegments()+1)
	require.Len(t, xmlData.Tracks[0].Segments[1].Points, track2Poly.NumSegments()+1)
	require.Len(t, xmlData.Tracks[0].Segments[2].Points, track4Line.NumSegments()+1)
	require.Len(t, xmlData.Tracks[1].Segments, 2)
	require.Len(t, xmlData.Tracks[1].Segments[0].Points, track4Line.NumSegments()+1)
	require.Len(t, xmlData.Tracks[1].Segments[1].Points, track2Poly.NumSegments()+1)
}

func TestDecodeExistingRoutes(t *testing.T) {
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
			name: "should return error when no points",
			args: args{
				data: file("restore_no_segments"),
			},
			wantErr: true,
		},
		{
			name: "should return error when route props is invalid",
			args: args{
				data: file("restore_invalid_route_props"),
			},
			wantErr: true,
		},
		{
			name: "should return empty routes when no tracks",
			args: args{
				data: file("restore_no_tracks"),
			},
			assert: func(routes []*navigator.Route) {
				require.Empty(t, routes)
			},
		},
		{
			name: "should return error when data is empty",
			args: args{
				data: nil,
			},
			wantErr: true,
		},
		{
			name: "should return error when gpx is invalid",
			args: args{
				data: []byte(`somedata`),
			},
			wantErr: true,
		},
		{
			name: "should return error when gpx format is broken",
			args: args{
				data: file("restore_broken_format"),
			},
			wantErr: true,
		},
		{
			name: "should restore multi routes from gpx",
			args: args{
				data: file("restore"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 2)
				require.Equal(t, 3, routes[0].NumTracks())
				require.Equal(t, 2, routes[1].NumTracks())
				require.Equal(t, "some", routes[0].Name().String())

				require.Equal(t, "df1d84dc-2516-4384-b477-a271c9b3fadb", routes[0].ID())
				require.Equal(t, "#b730a1", routes[0].Color())
				require.Equal(t, 5300.0, math.Floor(routes[0].Distance()))

				route := routes[0]
				// props
				foo, ok := route.Props().Int("foo")
				require.True(t, ok)
				require.Equal(t, 1, foo)
				baz, ok := route.Props().Bool("baz")
				require.True(t, ok)
				require.True(t, baz)
				ing, ok := route.Props().Int("ing")
				require.True(t, ok)
				require.Equal(t, 123, ing)
				val, ok := route.Props().String("val")
				require.True(t, ok)
				require.Equal(t, "val", val)

				// track
				track := routes[0].TrackAt(0)
				require.NotEmpty(t, track.NumSegments())
				require.Equal(t, "03a902c3-43a4-4f73-9bbe-b82a344e2ddc", track.ID())
				require.Equal(t, "#c51520", track.Color())
			},
			wantErr: false,
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

func TestDecodeNewRoutes(t *testing.T) {
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
			name: "should return error when gpx.Routes without points",
			args: args{
				data: file("new_tracks_without_points"),
			},
			wantErr: true,
		},
		{
			name: "should return new routes from gpx.Tracks",
			args: args{
				data: file("tracks"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 2)
				require.Equal(t, routes[0].NumTracks(), 2)
				require.Equal(t, routes[1].NumTracks(), 2)
			},
		},
		{
			name: "should return error when gpx.Routes without points",
			args: args{
				data: file("new_routes_without_points"),
			},
			wantErr: true,
		},
		{
			name: "should return error when gpx.Waypoint without points",
			args: args{
				data: file("new_waypoints_without_points"),
			},
			wantErr: true,
		},
		{
			name: "should return new routes from gpx.Routes",
			args: args{
				data: file("routes"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 1)
				require.Equal(t, routes[0].NumTracks(), 1)
			},
		},
		{
			name: "should return new routes from gpx.Waypoints",
			args: args{
				data: file("waypoints"),
			},
			assert: func(routes []*navigator.Route) {
				require.Len(t, routes, 1)
				require.Equal(t, routes[0].NumTracks(), 1)
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
