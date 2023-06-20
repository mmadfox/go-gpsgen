package navigator

import (
	"errors"

	pb "github.com/mmadfox/go-gpsgen/proto"
	"google.golang.org/protobuf/proto"
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

func (r *Route) SegmentAt(track int, segment int) *Segment {
	return r.tracks[track][segment]
}

func (r *Route) EachSegment(track int, fn func(seg *Segment)) {
	if track > len(r.tracks) {
		return
	}
	for j := 0; j < len(r.tracks[track]); j++ {
		fn(r.tracks[track][j])
	}
}

func (r *Route) MarshalBinary() ([]byte, error) {
	return proto.Marshal(r.ToProto())
}

func (r *Route) UnmarshalBinary(data []byte) error {
	p := new(pb.Route)
	if err := proto.Unmarshal(data, p); err != nil {
		return err
	}
	r.FromProto(p)
	return nil
}

func (r *Route) ToProto() *pb.Route {
	protoRoute := &pb.Route{
		Distance: r.dist,
		Tracks:   make([]*pb.Route_Track, len(r.tracks)),
	}
	for j := 0; j < len(r.tracks); j++ {
		protoTrack := &pb.Route_Track{
			Segmenets: make([]*pb.Route_Track_Segment, 0, len(r.tracks[j])),
		}
		for s := 0; s < len(r.tracks[j]); s++ {
			protoSegment := &pb.Route_Track_Segment{
				PointA: &pb.Point{
					Lat: r.tracks[j][s].pointA.X,
					Lon: r.tracks[j][s].pointA.Y,
				},
				PointB: &pb.Point{
					Lat: r.tracks[j][s].pointB.X,
					Lon: r.tracks[j][s].pointB.Y,
				},
				Bearing:  r.tracks[j][s].bearing,
				Distance: r.tracks[j][s].dist,
				Rel:      int64(r.tracks[j][s].rel),
			}
			protoTrack.Segmenets = append(protoTrack.Segmenets, protoSegment)
		}
		protoRoute.Tracks[j] = protoTrack
	}
	return protoRoute
}

func (r *Route) FromProto(route *pb.Route) {
	r.dist = route.Distance
	r.tracks = make([][]*Segment, len(route.Tracks))
	for j := 0; j < len(route.Tracks); j++ {
		r.tracks[j] = make([]*Segment, 0, len(route.Tracks[j].Segmenets))
		for s := 0; s < len(route.Tracks[j].Segmenets); s++ {
			protoSegment := route.Tracks[j].Segmenets[s]
			r.tracks[j] = append(r.tracks[j], &Segment{
				pointA:  Point{X: protoSegment.PointA.Lat, Y: protoSegment.PointA.Lon},
				pointB:  Point{X: protoSegment.PointB.Lat, Y: protoSegment.PointB.Lon},
				dist:    protoSegment.Distance,
				bearing: protoSegment.Bearing,
				rel:     int(protoSegment.Rel),
			})
		}
	}
}

type Segment struct {
	pointA  Point
	pointB  Point
	dist    float64
	bearing float64
	rel     int
}

func (s *Segment) Rel() int {
	return s.rel
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
