package geojson

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	data := loadData("geometry_collection")
	routes, err := Decode(data)
	require.NoError(t, err)
	data, err = Encode(routes)
	require.NoError(t, err)
	require.NotZero(t, data)
}

func TestDecode(t *testing.T) {
	type args struct {
		r []byte
	}
	tests := []struct {
		name       string
		args       args
		wantRoutes int
		wantErr    bool
	}{
		{
			name: "minified feature collection",
			args: args{
				r: loadData("minified"),
			},
			wantRoutes: 1,
		},
		{
			name: "geometry collection",
			args: args{
				r: loadData("geometry_collection"),
			},
			wantRoutes: 2,
		},
		{
			name: "circle",
			args: args{
				r: loadData("circle"),
			},
			wantRoutes: 1,
		},
		{
			name: "feature collection",
			args: args{
				r: loadData("feature_collection"),
			},
			wantRoutes: 1,
		},
		{
			name: "feature linestring",
			args: args{
				r: loadData("feature_linestring"),
			},
			wantRoutes: 1,
		},
		{
			name: "feature polygon",
			args: args{
				r: loadData("feature_polygon"),
			},
			wantRoutes: 1,
		},
		{
			name: "polygon",
			args: args{
				r: loadData("polygon"),
			},
			wantRoutes: 1,
		},
		{
			name: "multi polygon",
			args: args{
				r: loadData("multi_polygon"),
			},
			wantRoutes: 2,
		},
		{
			name: "linestring",
			args: args{
				r: loadData("linestring"),
			},
			wantRoutes: 1,
		},
		{
			name: "multi linestring",
			args: args{
				r: loadData("multi_linestring"),
			},
			wantRoutes: 4,
		},
		{
			name: "collection collection",
			args: args{
				r: loadData("geometry_collection_collection"),
			},
			wantRoutes: 12,
		},
		{
			name: "point",
			args: args{
				r: loadData("point"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoutesFromGeoJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Len(t, got, tt.wantRoutes)
		})
	}
}

func loadData(filename string) []byte {
	data, err := os.ReadFile("./testdata/geojson/" + filename + ".geojson")
	if err != nil {
		panic(err)
	}
	return data
}
