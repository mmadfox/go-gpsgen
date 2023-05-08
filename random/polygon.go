package random

import (
	"math"
	"math/rand"
	"time"
)

const (
	mp            = 8
	minFactor     = 1.0 / 10000
	maxFactor     = 5.0 / 10000
	defaultPoints = 16
	minZoom       = 0
	maxZoom       = 100
)

var genrand *rand.Rand

func init() {
	genrand = rand.New(rand.NewSource(time.Now().UnixNano()))

}

func Polygon(points int, zoom float64) [][2]float64 {
	if points < 0 {
		points = defaultPoints
	}
	if zoom < minZoom {
		zoom = minZoom
	}
	if zoom > maxZoom {
		zoom = maxZoom
	}
	coordinates := make([][2]float64, points+1)
	offsets := make([]float64, points)
	for i := 0; i < points; i++ {
		v := genrand.Float64()
		if i == 0 {
			offsets[i] = v
		} else {
			offsets[i] = offsets[i-1] + v
		}
	}
	last := offsets[len(offsets)-1]
	for i := 0; i < points; i++ {
		cur := (offsets[i] * mp * math.Phi) / last
		factor := randFactor(minFactor, maxFactor)
		var p1, p2 float64
		if zoom > 0 {
			p1 = factor * zoom * math.Sin(cur)
			p2 = factor * zoom * math.Cos(cur)
		} else {
			p1 = factor * math.Sin(cur)
			p2 = factor * math.Cos(cur)
		}
		coordinates[i] = [2]float64{p1, p2}
	}
	coordinates[len(coordinates)-1] = coordinates[0]
	return coordinates
}

func randFactor(min, max float64) float64 {
	return min + rand.ExpFloat64()*(max-min)
}
