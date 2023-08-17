package geo

import "math"

// Destination calculates the destination point given an initial point, distance, and bearing.
func Destination(lat, lon, meters, bearingDegrees float64) (
	destLat, destLon float64,
) {
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
