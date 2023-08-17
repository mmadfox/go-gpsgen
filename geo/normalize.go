package geo

// NormalizeCoordinates normalizes a set of coordinates by adding the given latitude and longitude offsets.
func NormalizeCoordinates(lat, lon float64, points [][2]float64) []LatLonPoint {
	coords := make([]LatLonPoint, len(points))
	for i := 0; i < len(points); i++ {
		coords[i] = LatLonPoint{
			points[i][0] + lat,
			points[i][1] + lon,
		}
	}
	return coords
}
