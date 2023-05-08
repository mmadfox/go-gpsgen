package geojson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mmadfox/go-gpsgen/navigator"
	gj "github.com/paulmach/go.geojson"
)

func Encode(routes []*navigator.Route) ([]byte, error) {
	featureCollection := gj.NewFeatureCollection()
	for i := 0; i < len(routes); i++ {
		route := routes[i]

		for j := 0; j < route.NumTracks(); j++ {
			isPolygon := route.IsPolygon(j)
			var geom *gj.Geometry
			switch isPolygon {
			case true:
				geom = makePolygonGeometry(j, route)
			case false:
				geom = makeLineGeometry(j, route)
			}
			feature := gj.NewFeature(geom)
			featureCollection.AddFeature(feature)
		}
	}
	return featureCollection.MarshalJSON()
}

func Decode(data []byte) ([]*navigator.Route, error) {
	data, err := io.ReadAll(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("geojson: file is empty")
	}
	switch {
	case is(data, "Feature"):
		return featureToRoutes(data)
	case is(data, "FeatureCollection"):
		return featureCollectionToRoutes(data)
	case is(data, "GeometryCollection"):
		return collectionToRoutes(data)
	default:
		return toRoutes(data)
	}
}

func makePolygonGeometry(track int, route *navigator.Route) *gj.Geometry {
	geom := gj.Geometry{
		Type:    gj.GeometryPolygon,
		Polygon: make([][][]float64, 0),
	}
	points := make([][]float64, 0)
	route.EachSegment(track, func(seg *navigator.Segment) {
		points = append(points,
			[]float64{seg.PointA().Y, seg.PointA().X},
			[]float64{seg.PointB().Y, seg.PointB().X},
		)
	})
	geom.Polygon = append(geom.Polygon, points)
	return &geom
}

func makeLineGeometry(track int, route *navigator.Route) *gj.Geometry {
	geom := gj.Geometry{
		Type:       gj.GeometryLineString,
		LineString: make([][]float64, 0),
	}
	route.EachSegment(track, func(seg *navigator.Segment) {
		geom.LineString = append(geom.LineString,
			[]float64{seg.PointA().Y, seg.PointA().X},
			[]float64{seg.PointB().Y, seg.PointB().X},
		)
	})
	return &geom
}

func featureCollectionToRoutes(data []byte) ([]*navigator.Route, error) {
	fc, err := gj.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}
	if len(fc.Features) == 0 {
		return nil, fmt.Errorf("geojson: routes not found")
	}
	var routes []*navigator.Route
	for i := 0; i < len(fc.Features); i++ {
		geom := fc.Features[i].Geometry
		localRoutes, err := geomToRoutes(geom)
		if err != nil {
			return nil, err
		}
		if len(localRoutes) == 0 {
			continue
		}
		routes = append(routes, localRoutes...)
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("geojson: routes not found")
	}
	return routes, nil
}

func collectionToRoutes(data []byte) ([]*navigator.Route, error) {
	cc := new(gj.Geometry)
	if err := json.Unmarshal(data, cc); err != nil {
		return nil, err
	}
	var routes []*navigator.Route
	for i := 0; i < len(cc.Geometries); i++ {
		geom := cc.Geometries[i]
		localRoutes, err := geomToRoutes(geom)
		if err != nil {
			return nil, err
		}
		if len(localRoutes) == 0 {
			continue
		}
		routes = append(routes, localRoutes...)
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("geojson: routes not found")
	}
	return routes, nil
}

func featureToRoutes(data []byte) ([]*navigator.Route, error) {
	f := new(gj.Feature)
	if err := json.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return geomToRoutes(f.Geometry)
}

func toRoutes(data []byte) ([]*navigator.Route, error) {
	g, err := gj.UnmarshalGeometry(data)
	if err != nil {
		return nil, err
	}
	routes, err := geomToRoutes(g)
	if err != nil {
		return nil, err
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("geojson: routes not found")
	}
	return routes, nil
}

func geomToRoutes(geom *gj.Geometry) (routes []*navigator.Route, err error) {
	if geom == nil {
		return nil, fmt.Errorf("geosjon: invalid format")
	}
	switch {
	case geom.IsCollection():
		routes = make([]*navigator.Route, 0, len(geom.Geometries))
		for j := 0; j < len(geom.Geometries); j++ {
			g2 := geom.Geometries[j]
			switch {
			case g2.IsLineString():
				r, err := value2dToRoutes(g2.LineString)
				if err != nil {
					return nil, err
				}
				routes = append(routes, r...)
			case g2.IsMultiLineString():
				r, err := value3dToRoutes(g2.MultiLineString)
				if err != nil {
					return nil, err
				}
				routes = append(routes, r...)
			case g2.IsMultiPolygon():
				r, err := value4dToRoutes(g2.MultiPolygon)
				if err != nil {
					return nil, err
				}
				routes = append(routes, r...)
			case g2.IsPolygon():
				r, err := value3dToRoutes(g2.Polygon)
				if err != nil {
					return nil, err
				}
				routes = append(routes, r...)
			}
		}
	case geom.IsLineString():
		routes, err = value2dToRoutes(geom.LineString)
		if err != nil {
			return nil, err
		}
	case geom.IsMultiLineString():
		routes, err = value3dToRoutes(geom.MultiLineString)
		if err != nil {
			return nil, err
		}
	case geom.IsMultiPolygon():
		routes, err = value4dToRoutes(geom.MultiPolygon)
		if err != nil {
			return nil, err
		}
	case geom.IsPolygon():
		routes, err = value3dToRoutes(geom.Polygon)
		if err != nil {
			return nil, err
		}
	}
	return routes, nil
}

func is(data []byte, typ string) bool {
	if len(data) < 60 {
		return false
	}
	return strings.Contains(string(data[:60]), `"`+typ+`"`)
}

func value2dToRoutes(values [][]float64) ([]*navigator.Route, error) {
	points := make([]navigator.Point, 0, len(values))
	for i := 0; i < len(values); i++ {
		pt := values[i]
		points = append(points, navigator.Point{
			X: pt[1], // lat
			Y: pt[0], // lon
		})
	}
	route, err := navigator.NewRoute([][]navigator.Point{points})
	if err != nil {
		return nil, err
	}
	return []*navigator.Route{route}, nil
}

func value3dToRoutes(values [][][]float64) ([]*navigator.Route, error) {
	routes := make([]*navigator.Route, 0, len(values))
	for i := 0; i < len(values); i++ {
		localRoutes, err := value2dToRoutes(values[i])
		if err != nil {
			return nil, err
		}
		routes = append(routes, localRoutes...)
	}
	return routes, nil
}

func value4dToRoutes(values [][][][]float64) ([]*navigator.Route, error) {
	routes := make([]*navigator.Route, 0, len(values))
	for i := 0; i < len(values); i++ {
		localRoutes, err := value3dToRoutes(values[i])
		if err != nil {
			return nil, err
		}
		routes = append(routes, localRoutes...)
	}
	return routes, nil
}
