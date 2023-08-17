package geo

import "math"

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
