package geojson

import (
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// ParseFeatureCollection parses a GeoJSON FeatureCollection from the provided data.
func ParseFeatureCollection(data []byte) (*geojson.FeatureCollection, error) {
	fc := geojson.NewFeatureCollection()
	if err := fc.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return fc, nil
}

// ToFeatureCollection converts a slice of navigator.Routes to a GeoJSON FeatureCollection.
func ToFeatureCollection(routes []*navigator.Route) *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()
	for i := 0; i < len(routes); i++ {
		route := routes[i]
		collection := make(orb.Collection, 0, route.NumTracks())
		tracksInfo := make([]trackInfo, route.NumTracks())
		for j := 0; j < route.NumTracks(); j++ {
			track := route.TrackAt(j)
			var geometry orb.Geometry
			points := copyPoints(track)
			if track.IsClosed() {
				geometry = orb.Polygon{orb.Ring(points)}
			} else {
				geometry = orb.LineString(points)
			}
			collection = append(collection, geometry)
			tracksInfo[j] = trackInfo{
				ID:          track.ID(),
				Name:        track.Name().String(),
				Color:       track.Color(),
				Distance:    track.Distance(),
				NumSegments: track.NumSegments(),
				Props:       track.Props(),
			}
		}
		if len(collection) == 0 {
			continue
		}

		var feature *geojson.Feature
		if len(collection) > 1 {
			feature = geojson.NewFeature(collection)
		} else {
			feature = geojson.NewFeature(collection[0])
		}

		feature.Properties["routeID"] = route.ID()
		feature.Properties["color"] = route.Color()
		feature.Properties["name"] = route.Name().String()
		feature.Properties["numTracks"] = route.NumTracks()
		feature.Properties["tracksInfo"] = tracksInfo
		feature.Properties["distance"] = route.Distance()
		feature.Properties["units"] = "meters"
		for k, v := range route.Props() {
			feature.Properties[k] = v
		}
		fc.Append(feature)
	}
	return fc
}

type trackInfo struct {
	ID          string                 `json:"trackID"`
	Name        string                 `json:"name"`
	Color       string                 `json:"color"`
	Distance    float64                `json:"distance"`
	NumSegments int                    `json:"numSegments"`
	Props       map[string]interface{} `json:"properties"`
}

func copyPoints(track *navigator.Track) []orb.Point {
	points := make([]orb.Point, 0, track.NumSegments()+1)
	for i := 0; i < track.NumSegments(); i++ {
		segment := track.SegmentAt(i)
		points = append(points, orb.Point{segment.PointA().Lon, segment.PointA().Lat})
		if i == track.NumSegments()-1 {
			points = append(points, orb.Point{segment.PointB().Lon, segment.PointB().Lat})
		}
	}
	return points
}
