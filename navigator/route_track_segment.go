package navigator

import (
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/proto"
)

// Segment represents a segment between two geographical points in a track.
type Segment struct {
	pointA  geo.LatLonPoint
	pointB  geo.LatLonPoint
	dist    float64
	bearing float64
	index   int
	rel     int
}

// SegmentSnapshot is a serialized snapshot of a Segment.
type SegmentSnapshot = proto.Snapshot_Navigator_Route_Track_Segment

// Index returns the index of the segment.
func (s Segment) Index() int {
	return s.index
}

// PointA returns the starting point of the segment.
func (s Segment) PointA() geo.LatLonPoint {
	return s.pointA
}

// PointB returns the ending point of the segment.
func (s Segment) PointB() geo.LatLonPoint {
	return s.pointB
}

// Distance returns the distance of the segment.
func (s Segment) Distance() float64 {
	return s.dist
}

// DestinationTo calculates the destination point from the starting point
// given a certain distance and bearing.
func (s Segment) DestinationTo(meters float64) geo.LatLonPoint {
	lat, lon := geo.Destination(s.pointA.Lat, s.pointA.Lon, meters, s.bearing)
	return geo.LatLonPoint{Lat: lat, Lon: lon}
}

// Bearing returns the bearing (direction) of the segment in degrees.
func (s Segment) Bearing() float64 {
	return s.bearing
}

// IsEmpty checks if the segment is empty (distance is zero).
func (s Segment) IsEmpty() bool {
	return s.dist == 0
}

// Snapshot returns a serialized snapshot of the segment.
func (s Segment) Snapshot() *SegmentSnapshot {
	return &SegmentSnapshot{
		PointA: &proto.Snapshot_PointLatLon{
			Lon: s.pointA.Lon,
			Lat: s.pointA.Lat,
		},
		PointB: &proto.Snapshot_PointLatLon{
			Lon: s.pointB.Lon,
			Lat: s.pointB.Lat,
		},
		Distance: s.dist,
		Bearing:  s.bearing,
		Index:    int64(s.index),
		Rel:      int64(s.rel),
	}
}

// SegmentFromSnapshot restores a segment from a snapshot.
func SegmentFromSnapshot(snap *SegmentSnapshot) Segment {
	if snap == nil {
		return Segment{}
	}
	return Segment{
		pointA:  geo.LatLonPoint{Lon: snap.PointA.Lon, Lat: snap.PointA.Lat},
		pointB:  geo.LatLonPoint{Lon: snap.PointB.Lon, Lat: snap.PointB.Lat},
		dist:    snap.Distance,
		bearing: snap.Bearing,
		index:   int(snap.Index),
		rel:     int(snap.Rel),
	}
}

func (s Segment) hasRelation() bool {
	return s.rel != -1
}

func makeSegments(points []geo.LatLonPoint) ([]Segment, float64, error) {
	if len(points) < 2 {
		return nil, 0, ErrInvalidRoutePath
	}

	numSeg := len(points) - 1
	segments := make([]Segment, 0, numSeg)
	dist := 0.0
	for i := 0; i < numSeg; i++ {
		newSeg := Segment{
			index:  i,
			pointA: points[i],
			pointB: points[i+1],
			rel:    -1,
		}
		newSeg.bearing = geo.Bearing(
			newSeg.pointA.Lat,
			newSeg.pointA.Lon,
			newSeg.pointB.Lat,
			newSeg.pointB.Lon)
		newSeg.dist = geo.Distance(
			newSeg.pointA.Lat,
			newSeg.pointA.Lon,
			newSeg.pointB.Lat,
			newSeg.pointB.Lon)
		segments = append(segments, newSeg)
		dist += newSeg.dist
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

	return segments, dist, nil
}
