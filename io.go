package gpsgen

import (
	"github.com/mmadfox/go-gpsgen/geojson"
	"github.com/mmadfox/go-gpsgen/gpx"
	"github.com/mmadfox/go-gpsgen/navigator"
)

// GeoJSONEncode encodes a slice of navigator routes into GeoJSON format.
// Returns the encoded GeoJSON data as a byte slice and an error if encoding fails.
func GeoJSONEncode(routes []*navigator.Route) ([]byte, error) {
	return geojson.Encode(routes)
}

// GeoJSONEncode encodes a slice of navigator routes into GeoJSON format.
// Returns the encoded GeoJSON data as a byte slice and an error if encoding fails.
func GeoJSONDecode(data []byte) ([]*navigator.Route, error) {
	return geojson.Decode(data)
}

// GPXEncode encodes a slice of navigator routes into GPX (GPS Exchange Format) format.
// Returns the encoded GPX data as a byte slice and an error if encoding fails.
func GPXEncode(routes []*navigator.Route) ([]byte, error) {
	return gpx.Encode(routes)
}

// GPXDecode decodes GPX (GPS Exchange Format) data into a slice of navigator routes.
// Returns the decoded navigator routes and an error if decoding fails.
func GPXDecode(data []byte) ([]*navigator.Route, error) {
	return gpx.Decode(data)
}
