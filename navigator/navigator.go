package navigator

import (
	"github.com/mmadfox/go-gpsgen/types"
)

type Navigator struct {
	routes          []*Route
	routeIndex      int
	trackIndex      int
	segmentIndex    int
	segmentDistance float64
	currentDistance float64
	offlineIndex    int
	point           Point
	location        *Location
	elevation       *types.Sensor
	offline         *types.Random
	totalDist       float64
}

func New(opts ...Option) (*Navigator, error) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	elevation, err := types.NewSensor("elevation",
		o.elevationMin,
		o.elevationMax,
		int(o.elevationAmplitude),
	)
	if err != nil {
		return nil, err
	}

	nav := &Navigator{
		routes:    make([]*Route, 0),
		location:  new(Location),
		elevation: elevation,
		offline:   types.NewRandom(o.minOffline, o.maxOffline),
	}

	return nav, nil
}

func (n *Navigator) AddRoutes(routes []*Route) {
	n.routes = append(n.routes, routes...)
	n.calcRouteDistance()
}

func (n *Navigator) AddRoute(route *Route) {
	n.routes = append(n.routes, route)
	n.calcRouteDistance()
}

func (n *Navigator) AddRouteFromPoints(points [][]Point) error {
	newRoute, err := NewRoute(points)
	if err != nil {
		return err
	}
	n.routes = append(n.routes, newRoute)
	n.calcRouteDistance()
	return nil
}

func (n *Navigator) CurrentDistance() float64 {
	return n.currentDistance
}

func (n *Navigator) TotalDistance() float64 {
	return n.totalDist
}

func (n *Navigator) IsFinish() bool {
	if n.routeIndex == 0 &&
		n.trackIndex == 0 &&
		n.segmentIndex == 0 &&
		n.currentDistance == 0.1 {
		n.currentDistance = 0
		return true
	}
	return false
}

func (n *Navigator) Segment() *Segment {
	return n.routes[n.routeIndex].tracks[n.trackIndex][n.segmentIndex]
}

func (n *Navigator) Route() *Route {
	return n.routes[n.routeIndex]
}

func (n *Navigator) TrackIndex() int {
	return n.trackIndex
}

func (n *Navigator) RouteIndex() int {
	return n.routeIndex
}

func (n *Navigator) SegmentIndex() int {
	return n.segmentIndex
}

func (n *Navigator) IsOnline() bool {
	return n.offlineIndex == 0
}

func (n *Navigator) Location() (loc Location) {
	loc.Lat = n.point.X
	loc.Lon = n.point.Y
	loc.Alt = n.elevation.ValueY()
	loc.Bearing = n.Segment().bearing
	loc.CurrentDistance = n.currentDistance
	loc.TotalDistance = n.totalDist
	loc.UTM = ToUTM(n.point.X, n.point.Y)
	loc.LatDMS, loc.LonDMS = ToDMS(n.point.X, n.point.Y)
	return loc
}

func (n *Navigator) NextOffline() {
	if n.offlineIndex > 0 {
		n.offlineIndex--
	}
}

func (n *Navigator) NextSensors(t float64) {
	n.elevation.Next(t)
}

func (n *Navigator) ToOffline() {
	n.offlineIndex = n.offline.Value()
}

func (n *Navigator) Next(tick float64, speed float64) (ok bool) {
	if len(n.routes) == 0 {
		return
	}

	if n.offlineIndex > 0 {
		n.offlineIndex--
		return
	}

	newDist := tick * speed

loop:
	for {
		if n.isValidDistance(newDist) {
			n.updateDist(newDist)
			if n.isStartSegment() {
				n.point = SegmentPoint(n.Segment(), 1)
			} else {
				n.point = SegmentPoint(n.Segment(), n.segmentDistance)
			}
			ok = true
			break
		} else {
			seg := n.nextSegment()
			foundSegment := seg != nil
			notFoundSegment := seg == nil
			switch {
			case foundSegment:
				if !seg.hasRelation() {
					n.offlineIndex = n.offline.Value()
					break loop
				}
				continue
			case notFoundSegment:
				if !n.nextTrack() {
					n.nextRoute()
				}
				n.offlineIndex = n.offline.Value()
				break loop
			}
		}
	}
	return
}

func (n *Navigator) updateDist(newDist float64) {
	n.segmentDistance += newDist
	n.currentDistance += newDist
}

func (n *Navigator) isValidDistance(dist float64) bool {
	segDist := n.Segment().dist
	return n.segmentDistance+dist <= segDist
}

func (n *Navigator) isStartSegment() bool {
	return n.segmentIndex == 0 && n.segmentDistance == 0
}

func (n *Navigator) nextSegment() *Segment {
	if n.segmentIndex == len(n.routes[n.routeIndex].tracks[n.trackIndex])-1 {
		return nil
	}
	segment := n.Segment()
	n.segmentDistance = 0
	n.segmentIndex++
	return segment
}

func (n *Navigator) nextTrack() bool {
	n.segmentDistance = 0
	n.segmentIndex = 0
	if n.trackIndex < len(n.routes[n.routeIndex].tracks)-1 {
		n.trackIndex++
		return true
	}
	n.trackIndex = 0
	return false
}

func (n *Navigator) nextRoute() {
	if n.routeIndex < len(n.routes)-1 {
		n.routeIndex++
		return
	}
	n.routeIndex = 0
	n.currentDistance = 0.1
}

func (n *Navigator) calcRouteDistance() {
	for i := 0; i < len(n.routes); i++ {
		n.totalDist += n.routes[i].dist
	}
}
