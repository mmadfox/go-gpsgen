package navigator

import (
	"errors"
)

var (
	ErrNoPoints         = errors.New("no route points")
	ErrInvalidRoutePath = errors.New("invalid route path")
)

type Point struct {
	X float64 // Lat
	Y float64 // Lon
}

type Route struct {
	dist   float64
	tracks [][]*Segment
}

func NewRoute(points [][]Point) (*Route, error) {
	if len(points) == 0 {
		return nil, ErrNoPoints
	}

	newRoute := Route{
		tracks: make([][]*Segment, 0, len(points)),
	}

	for i := 0; i < len(points); i++ {
		if err := newRoute.AddTrack(points[i]); err != nil {
			return nil, err
		}
	}

	return &newRoute, nil
}

func (r *Route) TotalDistance() float64 {
	return r.dist
}

func (r *Route) IsPolygon(track int) bool {
	if track > len(r.tracks) {
		return false
	}
	return r.tracks[track][0].pointA == r.tracks[track][len(r.tracks[track])-1].pointB
}

func (r *Route) AddTrack(points []Point) error {
	if len(points) < 2 {
		return ErrInvalidRoutePath
	}

	numSeg := len(points) - 1
	segments := make([]*Segment, 0, numSeg)
	for j := 0; j < numSeg; j++ {
		newSeg := &Segment{
			pointA: points[j],
			pointB: points[j+1],
			rel:    -1,
		}
		newSeg.bearing = BearingTo(
			newSeg.pointA.X,
			newSeg.pointA.Y,
			newSeg.pointB.X,
			newSeg.pointB.Y)
		newSeg.dist = DistanceTo(
			newSeg.pointA.X,
			newSeg.pointA.Y,
			newSeg.pointB.X,
			newSeg.pointB.Y)
		segments = append(segments, newSeg)
		r.dist += newSeg.dist
	}

	for s := 0; s < len(segments); s++ {
		if s == 0 {
			continue
		}
		if segments[s-1].pointB == segments[s].pointA {
			segments[s-1].rel = s - 1
			segments[s].rel = s
		}
	}

	r.tracks = append(r.tracks, segments)
	return nil
}

func (r *Route) Distance() float64 {
	return r.dist
}

func (r *Route) NumTracks() int {
	return len(r.tracks)
}

func (r *Route) NumSegments(trackIndex int) int {
	if trackIndex > len(r.tracks) {
		return -1
	}
	return len(r.tracks[trackIndex])
}

func (r *Route) EachSegment(track int, fn func(seg *Segment)) {
	if track > len(r.tracks) {
		return
	}
	for j := 0; j < len(r.tracks[track]); j++ {
		fn(r.tracks[track][j])
	}
}

type Segment struct {
	pointA  Point
	pointB  Point
	dist    float64
	bearing float64
	rel     int
}

func (s *Segment) PointA() Point {
	return s.pointA
}

func (s *Segment) PointB() Point {
	return s.pointB
}

func (s *Segment) Dist() float64 {
	return s.dist
}

func (s *Segment) Bearing() float64 {
	return s.bearing
}

func (s *Segment) hasRelation() bool {
	return s.rel != -1
}

func SegmentPoint(s *Segment, dist float64) Point {
	lat, lon := DestinationPoint(s.pointA.X, s.pointA.Y, dist, s.bearing)
	return Point{X: lat, Y: lon}
}
