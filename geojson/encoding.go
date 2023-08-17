package geojson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/properties"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

var (
	ErrNoRoutes            = errors.New("gpsgen/geojson: no routes")
	ErrInvalidGeometryType = errors.New("gpsgen/geojson: invalid geometry type")
	ErrInvalidRoute        = errors.New("gpsgen/geojson: invalid route")
)

// Encode converts a slice of navigator.Routes to GeoJSON format.
func Encode(routes []*navigator.Route) ([]byte, error) {
	if len(routes) == 0 {
		return nil, ErrNoRoutes
	}
	return ToFeatureCollection(routes).MarshalJSON()
}

// Decode converts GeoJSON data into a slice of navigator.Routes.
func Decode(data []byte) ([]*navigator.Route, error) {
	fc := geojson.NewFeatureCollection()
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, err
	}
	routes := make([]*navigator.Route, 0, len(fc.Features))
	for i := 0; i < len(fc.Features); i++ {
		feature := fc.Features[i]
		if feature.Type != "Feature" {
			continue
		}
		var (
			route      *navigator.Route
			tracksInfo []trackInfo
			err        error
		)

		if routeExists(feature.Properties) {
			route, tracksInfo, err = restoreRoute(feature.Properties)
			if err != nil {
				return nil, err
			}
			resetProps(feature.Properties)
		} else {
			route = navigator.NewRoute()
		}

		parseProperties(route, feature.Properties)

		if err := parseGeometries(route, tracksInfo, feature.Geometry); err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}
	return routes, nil
}

func routeExists(props map[string]interface{}) bool {
	_, ok1 := props["routeID"]
	_, ok2 := props["color"]
	return ok1 && ok2
}

func restoreRoute(props map[string]interface{}) (*navigator.Route, []trackInfo, error) {
	p := properties.Properties(props)
	routeID, _ := p.String("routeID")
	color, _ := p.String("color")
	name, _ := p.String("name")
	tracks, ok := p["tracksInfo"].([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("gpsgen/geojson: invalid tracks info")
	}
	tracksInfo := make([]trackInfo, len(tracks))
	for i := 0; i < len(tracks); i++ {
		track, ok := tracks[i].(map[string]interface{})
		if !ok {
			continue
		}
		trackProps := properties.Properties(track)
		trackID, _ := trackProps.String("trackID")
		color, _ := trackProps.String("color")
		name, _ := trackProps.String("name")
		ti := trackInfo{
			ID:    trackID,
			Color: color,
			Name:  name,
		}
		properties, ok := track["properties"].(map[string]interface{})
		if ok {
			ti.Props = properties
		}
		tracksInfo[i] = ti
	}
	route := navigator.RestoreRoute(routeID, color, p)
	if len(name) > 0 {
		_ = route.ChangeName(name)
	}
	return route, tracksInfo, nil
}

func resetProps(props map[string]interface{}) {
	delete(props, "routeID")
	delete(props, "color")
	delete(props, "distance")
	delete(props, "tracksInfo")
	delete(props, "numTracks")
	delete(props, "name")
	delete(props, "units")
}

func parseProperties(route *navigator.Route, props map[string]interface{}) {
	if props == nil {
		return
	}
	props = properties.Properties(props)
	route.Props().Merge(props)
}

func parseGeometries(route *navigator.Route, tracks []trackInfo, geom orb.Geometry) error {
	collection, ok := geom.(orb.Collection)
	if !ok {
		if tracks != nil && len(tracks) != 1 {
			return ErrInvalidRoute
		}
		return parseGeometry(route, tracks, geom)
	}

	if tracks != nil && len(tracks) != len(collection) {
		return ErrInvalidRoute
	}

	for i := 0; i < len(collection); i++ {
		if err := parseGeometry(route, tracks, collection[i]); err != nil {
			return err
		}
	}
	return nil
}

func parseGeometry(route *navigator.Route, tracks []trackInfo, geometry orb.Geometry) (err error) {
	trackExists := len(tracks) > 0
	switch geom := geometry.(type) {
	case orb.LineString:
		var track *navigator.Track
		var err error
		if !trackExists {
			track, err = navigator.NewTrack(toPoints(geom))
		} else {
			trackInfo := tracks[0]
			track, err = navigator.RestoreTrack(trackInfo.ID, trackInfo.Color, toPoints(geom))
			track.Props().Merge(trackInfo.Props)
		}
		if err != nil {
			return err
		}
		route.AddTrack(track)
	case orb.MultiLineString:
		if trackExists {
			return ErrInvalidRoute
		}
		for i := 0; i < len(geom); i++ {
			track, err := navigator.NewTrack(toPoints(geom[i]))
			if err != nil {
				return err
			}
			route.AddTrack(track)
		}
	case orb.Polygon:
		if len(geom) == 0 {
			return
		}
		var track *navigator.Track
		var err error
		if !trackExists {
			track, err = navigator.NewTrack(toPoints(geom[0]))
		} else {
			trackInfo := tracks[0]
			track, err = navigator.RestoreTrack(trackInfo.ID, trackInfo.Color, toPoints(geom[0]))
			track.Props().Merge(trackInfo.Props)
		}
		if err != nil {
			return err
		}
		route.AddTrack(track)
	case orb.MultiPolygon:
		if trackExists {
			return ErrInvalidRoute
		}
		for i := 0; i < len(geom); i++ {
			if len(geom[i]) > 0 {
				track, err := navigator.NewTrack(toPoints(geom[i][0]))
				if err != nil {
					return err
				}
				route.AddTrack(track)
			}
		}
	case orb.MultiPoint:
		if trackExists {
			return ErrInvalidRoute
		}
		track, err := navigator.NewTrack(toPoints(geom))
		if err != nil {
			return err
		}
		route.AddTrack(track)
	default:
		err = fmt.Errorf("%w - %s", ErrInvalidGeometryType, geom.GeoJSONType())
	}
	return
}

func toPoints(points []orb.Point) []geo.LatLonPoint {
	latLons := make([]geo.LatLonPoint, len(points))
	for i := 0; i < len(points); i++ {
		latLons[i] = geo.LatLonPoint{
			Lon: points[i].Lon(),
			Lat: points[i].Lat(),
		}
	}
	return latLons
}
