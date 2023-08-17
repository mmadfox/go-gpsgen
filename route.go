package gpsgen

import (
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/random"
)

// Route level constants define different levels of route complexity.
const (
	RouteLevelXS  = 1
	RouteLevelS   = 10
	RouteLevelM   = 60
	RouteLevelL   = 120
	RouteLevelXL  = 300
	RouteLevelXXL = 600
)

// RandomRoute generates a random route with specified parameters.
// The function generates tracks within the route, using a specified number of tracks and complexity level.
// The generated route is centered around the provided latitude and longitude.
// Returns a random route with tracks, or nil if an error occurs during track creation.
func RandomRoute(lon, lat float64, numTrack int, level int) *navigator.Route {
	if numTrack < 0 {
		numTrack = 1
	}
	if level < RouteLevelS {
		level = RouteLevelS
	}
	if level > RouteLevelXXL {
		level = RouteLevelXXL
	}
	route := navigator.NewRoute()
	routeName := "Route-" + random.String(8)
	route.ChangeName(routeName)
	for i := 0; i < numTrack; i++ {
		rawPoints := random.Polygon(16, float64(level))
		points := geo.NormalizeCoordinates(lat, lon, rawPoints)
		track, _ := navigator.NewTrack(points)
		trackName := "Track-" + random.String(8)
		track.ChangeName(trackName)
		route.AddTrack(track)
	}
	return route
}
