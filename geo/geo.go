package geo

import (
	"fmt"
	"math"
)

const (
	earthRadius = 6371e3
	radians     = math.Pi / 180
	degrees     = 180 / math.Pi
	piR         = math.Pi * earthRadius
	twoPiR      = 2 * piR
)

type LatLonPoint struct {
	Lat float64
	Lon float64
}

func (s LatLonPoint) String() string {
	return fmt.Sprintf("Point{Lon: %f, Lat: %f}",
		s.Lon, s.Lat)
}

type BBox struct {
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
}

func (b BBox) In(pt LatLonPoint) bool {
	return b.MinLon <= pt.Lon &&
		b.MinLat <= pt.Lat &&
		b.MaxLon >= pt.Lon &&
		b.MaxLat >= pt.Lat
}

func (b BBox) Center() LatLonPoint {
	return LatLonPoint{
		Lat: (b.MinLat + b.MaxLat) / 2.0,
		Lon: (b.MinLon + b.MaxLon) / 2.0,
	}
}
