package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/properties"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	ErrNoRoutes     = errors.New("gpsgen/gpx: no routes")
	ErrInvalidRoute = errors.New("gpsgen/gpx: invalid route")
)

// Encode converts a slice of navigator.Routes to GPX format.
func Encode(routes []*navigator.Route) ([]byte, error) {
	if len(routes) == 0 {
		return nil, ErrNoRoutes
	}

	gpxData := new(gpx.GPX)
	gpxData.Tracks = make([]gpx.GPXTrack, 0)
	gpxData.AuthorName = "go-gpsgen"

	routesNode := gpx.ExtensionNode{
		XMLName: xml.Name{Local: "gpsgen"},
	}

	// map: tracks=routes, segments=tracks
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		gpxTrack := gpx.GPXTrack{
			Segments: make([]gpx.GPXTrackSegment, 0, route.NumTracks()),
		}

		routeNode := gpx.ExtensionNode{
			XMLName: xml.Name{Local: "route"},
		}
		routeNode.SetAttr("id", route.ID())
		routeNode.SetAttr("name", route.Name().String())
		routeNode.SetAttr("color", route.Color())
		routeNode.SetAttr("distance", f2s(route.Distance()))
		routeNode.SetAttr("units", "meters")
		if attrs := makeAttrs(route.Props()); len(attrs) > 0 {
			routeNode.Attrs = append(routeNode.Attrs, attrs...)
		}

		for j := 0; j < route.NumTracks(); j++ {
			track := route.TrackAt(j)
			gpxSegment := gpx.GPXTrackSegment{
				Points: copyPoints(track),
			}

			trackNode := gpx.ExtensionNode{
				XMLName: xml.Name{Local: "track"},
			}
			trackNode.SetAttr("name", track.Name().String())
			trackNode.SetAttr("color", track.Color())
			trackNode.SetAttr("distance", f2s(track.Distance()))
			trackNode.Data = track.ID()
			routeNode.Nodes = append(routeNode.Nodes, trackNode)
			if attrs := makeAttrs(track.Props()); len(attrs) > 0 {
				trackNode.Attrs = append(trackNode.Attrs, attrs...)
			}

			gpxTrack.Segments = append(gpxTrack.Segments, gpxSegment)
		}

		routesNode.Nodes = append(routesNode.Nodes, routeNode)
		gpxData.Tracks = append(gpxData.Tracks, gpxTrack)
	}

	gpxData.Extensions.Nodes = append(gpxData.Extensions.Nodes, routesNode)

	return gpxData.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
}

// Decode converts GPX data into a slice of navigator.Routes.
func Decode(data []byte) ([]*navigator.Route, error) {
	if len(data) == 0 {
		return nil, ErrNoRoutes
	}
	gpxData, err := gpx.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	var routes []*navigator.Route
	if routesExists(gpxData) {
		// restore routes
		if err := validateRoutes(gpxData); err != nil {
			return nil, err
		}
		root := gpxData.Extensions.Nodes[0]
		routes = make([]*navigator.Route, 0, len(gpxData.Tracks))
		for i := 0; i < len(gpxData.Tracks); i++ {
			gpxTrack := gpxData.Tracks[i]
			if len(gpxTrack.Segments) == 0 {
				continue
			}
			routeInfo := root.Nodes[i]
			routeProps := makeProps(routeInfo.Attrs)
			routeID, ok := routeProps.String("id")
			if !ok {
				return nil, ErrInvalidRoute
			}
			color, _ := routeProps.String("color")
			routeName, _ := routeProps.String("name")
			resetProps(routeProps)
			route := navigator.RestoreRoute(routeID, color, routeProps)
			if len(routeName) > 0 {
				_ = route.ChangeName(routeName)
			}
			var trackErr error
		loop:
			for j := 0; j < len(gpxTrack.Segments); j++ {
				gpxSegment := gpxTrack.Segments[j]
				trackInfo := root.Nodes[i].Nodes[j]
				trackID := trackInfo.Data
				trackProps := makeProps(trackInfo.Attrs)
				color, _ := trackProps.String("color")
				trackName, _ := trackProps.String("name")
				resetProps(trackProps)
				track, err := navigator.RestoreTrack(
					trackID,
					color,
					makePoints(gpxSegment.Points),
				)
				if err != nil {
					trackErr = err
					break loop
				}
				track.Props().Merge(trackProps)
				if len(trackName) > 0 {
					_ = track.ChangeName(trackName)
				}
				route.AddTrack(track)
			}
			if trackErr != nil {
				return nil, trackErr
			}
			routes = append(routes, route)
		}

	} else {
		// new routes
		if len(gpxData.Routes) > 0 {
			routes = make([]*navigator.Route, 0, len(gpxData.Routes))
			for i := 0; i < len(gpxData.Routes); i++ {
				route := navigator.NewRoute()
				track, err := navigator.NewTrack(makePoints(gpxData.Routes[i].Points))
				if err != nil {
					return nil, err
				}
				route.AddTrack(track)
				routes = append(routes, route)
			}
		}

		if len(gpxData.Tracks) > 0 {
			routes = make([]*navigator.Route, 0, len(gpxData.Tracks))
			for i := 0; i < len(gpxData.Tracks); i++ {
				route := navigator.NewRoute()
				var trackErr error
				for j := 0; j < len(gpxData.Tracks[i].Segments); j++ {
					gpxPoints := gpxData.Tracks[i].Segments[j].Points
					track, err := navigator.NewTrack(makePoints(gpxPoints))
					if err != nil {
						trackErr = err
						break
					}
					route.AddTrack(track)
				}
				if trackErr != nil {
					return nil, trackErr
				}
				routes = append(routes, route)
			}
		}

		if len(gpxData.Waypoints) > 0 {
			route := navigator.NewRoute()
			track, err := navigator.NewTrack(makePoints(gpxData.Waypoints))
			if err != nil {
				return nil, err
			}
			route.AddTrack(track)
			routes = []*navigator.Route{route}
		}
	}

	return routes, nil
}

func routesExists(gpxData *gpx.GPX) bool {
	if len(gpxData.Extensions.Nodes) == 0 || gpxData.Extensions.Nodes[0].XMLName.Local != "gpsgen" {
		return false
	}
	return true
}

func f2s(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func makeAttrs(props properties.Properties) []xml.Attr {
	if len(props) == 0 {
		return nil
	}
	attrs := make([]xml.Attr, 0, len(props))
	for k, v := range props {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: fmt.Sprintf("%v", v),
		})
	}
	return attrs
}

func makeProps(attrs []xml.Attr) properties.Properties {
	props := properties.Make()
	for i := 0; i < len(attrs); i++ {
		key := attrs[i].Name.Local
		val := value(attrs[i].Value)
		props[key] = val
	}
	return props
}

func resetProps(props properties.Properties) {
	delete(props, "id")
	delete(props, "color")
	delete(props, "distance")
	delete(props, "tracksInfo")
	delete(props, "numTracks")
	delete(props, "name")
	delete(props, "units")
}

func value(s string) interface{} {
	intVal, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return intVal
	}

	boolVal, err := strconv.ParseBool(s)
	if err == nil {
		return boolVal
	}

	floatVal, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return floatVal
	}

	return s
}

func validateRoutes(gpxData *gpx.GPX) error {
	if len(gpxData.Tracks) == 0 || len(gpxData.Extensions.Nodes) == 0 {
		return ErrInvalidRoute
	}
	if len(gpxData.Extensions.Nodes[0].Nodes) != len(gpxData.Tracks) {
		return ErrInvalidRoute
	}
	root := gpxData.Extensions.Nodes[0].Nodes
	for i := 0; i < len(gpxData.Tracks); i++ {
		a := len(gpxData.Tracks[i].Segments)
		b := len(root[i].Nodes)
		if a != b {
			return ErrInvalidRoute
		}
	}
	return nil
}

func copyPoints(track *navigator.Track) []gpx.GPXPoint {
	points := make([]gpx.GPXPoint, 0, track.NumSegments()+1)
	for i := 0; i < track.NumSegments(); i++ {
		seg := track.SegmentAt(i)
		points = append(points, gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  seg.PointA().Lat,
				Longitude: seg.PointA().Lon,
			},
		})
		if i == track.NumSegments()-1 {
			points = append(points, gpx.GPXPoint{
				Point: gpx.Point{
					Latitude:  seg.PointB().Lat,
					Longitude: seg.PointB().Lon,
				},
			})
		}
	}
	return points
}

func makePoints(gpxPoints []gpx.GPXPoint) []geo.LatLonPoint {
	points := make([]geo.LatLonPoint, 0, len(gpxPoints))
	for i := 0; i < len(gpxPoints); i++ {
		points = append(points, geo.LatLonPoint{
			Lat: gpxPoints[i].Latitude,
			Lon: gpxPoints[i].Longitude,
		})
	}
	return points
}
