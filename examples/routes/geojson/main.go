package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/geojson"
)

func main() {
	geoJSONRoute := `{"type":"FeatureCollection","features":[{"type":"Feature","properties":{},"geometry":{"coordinates":[[37.73185507353088,55.59010742814678],[37.740460827380986,55.585012323942266]],"type":"LineString"}}]}`
	routes, err := geojson.Decode([]byte(geoJSONRoute))
	if err != nil {
		panic(err)
	}
	myTracker, err := gpsgen.Tracker("Tracker78-1", nil, routes...)
	if err != nil {
		panic(err)
	}
	_ = myTracker
}
