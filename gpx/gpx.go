package gpx

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/tkrajina/gpxgo/gpx"
)

func Encode(routes []*navigator.Route) ([]byte, error) {
	gpxData := new(gpx.GPX)
	gpxData.Tracks = make([]gpx.GPXTrack, 0)
	gpxData.AuthorName = "go-gpsgen"
	for i := 0; i < len(routes); i++ {
		track := gpx.GPXTrack{
			Segments: make([]gpx.GPXTrackSegment, 0),
		}
		route := routes[i]
		for j := 0; j < route.NumTracks(); j++ {
			trackPoint := make([]gpx.GPXPoint, 0)
			route.EachSegment(j, func(seg *navigator.Segment) {
				trackPoint = append(trackPoint, gpx.GPXPoint{
					Point: gpx.Point{
						Latitude:  seg.PointA().X,
						Longitude: seg.PointA().Y,
					},
				})
				trackPoint = append(trackPoint, gpx.GPXPoint{
					Point: gpx.Point{
						Latitude:  seg.PointB().X,
						Longitude: seg.PointB().Y,
					},
				})
			})
			seq := gpx.GPXTrackSegment{
				Points: trackPoint,
			}
			track.Segments = append(track.Segments, seq)
		}
		gpxData.Tracks = append(gpxData.Tracks, track)
	}
	return gpxData.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
}

func Decode(data []byte) ([]*navigator.Route, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gpx: file is empty")
	}
	gpxData, err := gpx.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	routes := make([]*navigator.Route, 0, 4)
	if len(gpxData.Routes) > 0 {
		for i := 0; i < len(gpxData.Routes); i++ {
			points := make([]navigator.Point, 0, len(gpxData.Routes[i].Points))
			for j := 0; j < len(gpxData.Routes[i].Points); j++ {
				points = append(points, navigator.Point{
					X: gpxData.Routes[i].Points[j].Latitude,
					Y: gpxData.Routes[i].Points[j].Longitude,
				})
			}
			route, err := navigator.NewRoute([][]navigator.Point{points})
			if err != nil {
				return nil, err
			}
			routes = append(routes, route)
		}
	}
	if len(gpxData.Tracks) > 0 {
		for i := 0; i < len(gpxData.Tracks); i++ {
			points := make([][]navigator.Point, 0)
			for j := 0; j < len(gpxData.Tracks[i].Segments); j++ {
				segment := make([]navigator.Point, 0, 8)
				for x := 0; x < len(gpxData.Tracks[i].Segments[j].Points); x++ {
					segment = append(segment, navigator.Point{
						X: gpxData.Tracks[i].Segments[j].Points[x].Latitude,
						Y: gpxData.Tracks[i].Segments[j].Points[x].Longitude,
					})
				}
				points = append(points, segment)
			}
			route, err := navigator.NewRoute(points)
			if err != nil {
				return nil, err
			}
			routes = append(routes, route)
		}
	}
	if len(gpxData.Waypoints) > 0 {
		points := make([]navigator.Point, 0, len(gpxData.Waypoints))
		for i := 0; i < len(gpxData.Waypoints); i++ {
			points = append(points, navigator.Point{
				X: gpxData.Waypoints[i].Latitude,
				Y: gpxData.Waypoints[i].Longitude,
			})
		}
		route, err := navigator.NewRoute([][]navigator.Point{points})
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}
