package navigator

import "math"

const (
	earthRadius = 6371e3
	radians     = math.Pi / 180
	degrees     = 180 / math.Pi
	piR         = math.Pi * earthRadius
	twoPiR      = 2 * piR
)

// BearingTo returns the (initial) bearing from point 'A' to point 'B'.
func BearingTo(latA, lonA, latB, lonB float64) float64 {
	// tanθ = sinΔλ⋅cosφ2 / cosφ1⋅sinφ2 − sinφ1⋅cosφ2⋅cosΔλ
	// see mathforum.org/library/drmath/view/55417.html for derivation

	φ1 := latA * radians
	φ2 := latB * radians
	Δλ := (lonB - lonA) * radians
	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	θ := math.Atan2(y, x)

	return math.Mod(θ*degrees+360, 360)
}

func DistanceTo(latA, lonA, latB, lonB float64) (meters float64) {
	a := Haversine(latA, lonA, latB, lonB)
	return DistanceFromHaversine(a)
}

func Haversine(latA, lonA, latB, lonB float64) float64 {
	φ1 := latA * radians
	λ1 := lonA * radians
	φ2 := latB * radians
	λ2 := lonB * radians
	Δφ := φ2 - φ1
	Δλ := λ2 - λ1
	sΔφ2 := math.Sin(Δφ / 2)
	sΔλ2 := math.Sin(Δλ / 2)
	return sΔφ2*sΔφ2 + math.Cos(φ1)*math.Cos(φ2)*sΔλ2*sΔλ2
}

func NormalizeDistance(meters float64) float64 {
	return math.Mod(meters, twoPiR)
}

func DistanceToHaversine(meters float64) float64 {
	sin := math.Sin(0.5 * meters / earthRadius)
	return sin * sin
}

func DistanceFromHaversine(haversine float64) float64 {
	return earthRadius * 2 * math.Asin(math.Sqrt(haversine))
}

func DestinationPoint(lat, lon, meters, bearingDegrees float64) (
	destLat, destLon float64,
) {
	// see http://williams.best.vwh.net/avform.htm#LL
	δ := meters / earthRadius // angular distance in radians
	θ := bearingDegrees * radians
	φ1 := lat * radians
	λ1 := lon * radians
	φ2 := math.Asin(math.Sin(φ1)*math.Cos(δ) +
		math.Cos(φ1)*math.Sin(δ)*math.Cos(θ))
	λ2 := λ1 + math.Atan2(math.Sin(θ)*math.Sin(δ)*math.Cos(φ1),
		math.Cos(δ)-math.Sin(φ1)*math.Sin(φ2))
	λ2 = math.Mod(λ2+3*math.Pi, 2*math.Pi) - math.Pi // normalise to -180..+180°
	return φ2 * degrees, λ2 * degrees
}

func NormalizeCoordinates(lat, lon float64, points [][2]float64) []Point {
	coords := make([]Point, len(points))
	for i := 0; i < len(points); i++ {
		coords[i] = Point{
			points[i][0] + lat,
			points[i][1] + lon,
		}
	}
	return coords
}
