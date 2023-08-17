package geo

import "math"

// Distance calculates the distance in meters between two points (latA, lonA) and (latB, lonB).
func Distance(latA, lonA, latB, lonB float64) (meters float64) {
	a := Haversine(latA, lonA, latB, lonB)
	return DistanceFromHaversine(a)
}

// Haversine calculates the haversine of the angular distance between two points (latA, lonA) and (latB, lonB).
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

// NormalizeDistance normalizes a distance value by wrapping it around the Earth's circumference.
func NormalizeDistance(meters float64) float64 {
	return math.Mod(meters, twoPiR)
}

// DistanceToHaversine converts a distance in meters to its haversine value.
func DistanceToHaversine(meters float64) float64 {
	sin := math.Sin(0.5 * meters / earthRadius)
	return sin * sin
}

// DistanceFromHaversine converts a haversine value to a distance in meters.
func DistanceFromHaversine(haversine float64) float64 {
	return earthRadius * 2 * math.Asin(math.Sqrt(haversine))
}
