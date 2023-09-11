package random

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/stretchr/testify/require"
)

func TestCountryName(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "should return error when country code is empty",
			args: args{
				code: "",
			},
			wantErr: ErrCountryNotFound,
		},
		{
			name: "should return error when country not found",
			args: args{
				code: "xy",
			},
			wantErr: ErrCountryNotFound,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				code: "zw",
			},
			want: "Zimbabwe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountryName(tt.args.code)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CountryName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				return
			}
			if got != tt.want {
				t.Errorf("CountryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoundingBox(t *testing.T) {
	type args struct {
		countryCode string
	}
	tests := []struct {
		name    string
		args    args
		want    geo.BBox
		wantErr bool
	}{
		{
			name: "should return error when country code is empty",
			args: args{
				countryCode: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when bounding box not found",
			args: args{
				countryCode: "xy",
			},
			wantErr: true,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				countryCode: "zw",
			},
			want: geo.BBox{
				MinLon: 25.26,
				MinLat: -22.27,
				MaxLon: 32.85,
				MaxLat: -15.51,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BoundingBox(tt.args.countryCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoundingBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoundingBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatLonByCountry(t *testing.T) {
	type args struct {
		countryCode string
	}
	bbox := func(code string) geo.BBox {
		rect, err := BoundingBox(code)
		if err != nil {
			panic(err)
		}
		return rect
	}
	tests := []struct {
		name    string
		args    args
		want    geo.BBox
		wantErr bool
	}{
		{
			name: "should return error country code is empty",
			args: args{
				countryCode: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when country not found",
			args: args{
				countryCode: "zy",
			},
			wantErr: true,
		},
		{
			name: "should not return error when AF country found",
			args: args{
				countryCode: "af",
			},
			want:    bbox("af"),
			wantErr: false,
		},
		{
			name: "should not return error when TZ country found",
			args: args{
				countryCode: "tz",
			},
			want:    bbox("tz"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LatLonByCountry(tt.args.countryCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("LatLonByCountry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, tt.want.In(got))
		})
	}
}

func TestLatLon(t *testing.T) {
	bbox := geo.BBox{
		MinLon: -180,
		MinLat: -90,
		MaxLon: 180,
		MaxLat: 90,
	}
	for i := 0; i < 10; i++ {
		pt := LatLon()
		require.True(t, bbox.In(pt))
	}
}

func TestBBoxCenter(t *testing.T) {
	bbox := geo.BBox{
		MinLon: 1.213284,
		MinLat: 18.600047,
		MaxLon: 11.587763,
		MaxLat: 25.518278,
	}
	expected := geo.LatLonPoint{
		Lat: 22.0591625,
		Lon: 6.4005235,
	}
	center := bbox.Center()
	require.Equal(t, expected, center)
}
