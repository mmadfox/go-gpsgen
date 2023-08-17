package geo

import "math"

// Bearing calculates the bearing between two geographical points.
func Bearing(latA, lonA, latB, lonB float64) float64 {
	φ1 := latA * radians
	φ2 := latB * radians
	Δλ := (lonB - lonA) * radians
	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	θ := math.Atan2(y, x)
	return math.Mod(θ*degrees+360, 360)
}
