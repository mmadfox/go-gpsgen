package navigator

import (
	"github.com/mmadfox/go-gpsgen/proto"
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
	skipOffline     bool
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
		routes:      make([]*Route, 0),
		location:    new(Location),
		elevation:   elevation,
		skipOffline: o.skipOffline,
	}
	switch o.skipOffline {
	case false:
		nav.offline = types.NewRandom(o.minOffline, o.maxOffline)
	}

	return nav, nil
}

func (n *Navigator) ToProto() *proto.NavigatorState {
	protoRoutes := make([]*proto.NavigatorState_Route, len(n.routes))
	for i := 0; i < len(n.routes); i++ {
		route := n.routes[i]
		protoRoute := &proto.NavigatorState_Route{
			Distance: route.dist,
			Tracks:   make([]*proto.NavigatorState_Route_Track, len(route.tracks)),
		}
		for j := 0; j < len(route.tracks); j++ {
			protoTrack := &proto.NavigatorState_Route_Track{
				Segmenets: make([]*proto.NavigatorState_Route_Track_Segment, 0, len(route.tracks[j])),
			}
			for s := 0; s < len(route.tracks[j]); s++ {
				protoSegment := &proto.NavigatorState_Route_Track_Segment{
					PointA: &proto.NavigatorState_Point{
						Lat: route.tracks[j][s].pointA.X,
						Lon: route.tracks[j][s].pointA.Y,
					},
					PointB: &proto.NavigatorState_Point{
						Lat: route.tracks[j][s].pointB.X,
						Lon: route.tracks[j][s].pointB.Y,
					},
					Bearing:  route.tracks[j][s].bearing,
					Distance: route.tracks[j][s].dist,
					Rel:      int64(route.tracks[j][s].rel),
				}
				protoTrack.Segmenets = append(protoTrack.Segmenets, protoSegment)
			}
			protoRoute.Tracks[i] = protoTrack
		}
		protoRoutes[i] = protoRoute
	}
	nav := &proto.NavigatorState{
		Routes:          protoRoutes,
		RouteIndex:      int64(n.routeIndex),
		TrackIndex:      int64(n.trackIndex),
		SegmentIndex:    int64(n.segmentIndex),
		SegmentDistance: n.segmentDistance,
		CurrentDistance: n.currentDistance,
		OfflineIndex:    int64(n.offlineIndex),
		Point: &proto.NavigatorState_Point{
			Lat: n.point.X,
			Lon: n.point.Y,
		},
		Elevation:     n.elevation.ToProto(),
		TotalDistance: n.totalDist,
		SkipOffline:   n.skipOffline,
	}
	if n.offline != nil {
		nav.OfflineMin = int64(n.offline.Min())
		nav.OfflineMax = int64(n.offline.Max())
	}
	return nav
}

func (n *Navigator) FromProto(nav *proto.NavigatorState) {
	n.routes = make([]*Route, 0, len(nav.Routes))
	for i := 0; i < len(nav.Routes); i++ {
		protoRoute := nav.Routes[i]
		route := Route{
			dist:   protoRoute.Distance,
			tracks: make([][]*Segment, len(protoRoute.Tracks)),
		}
		for j := 0; j < len(protoRoute.Tracks); j++ {
			route.tracks[j] = make([]*Segment, 0, len(protoRoute.Tracks[j].Segmenets))
			for s := 0; s < len(protoRoute.Tracks[j].Segmenets); s++ {
				protoSegment := protoRoute.Tracks[j].Segmenets[s]
				route.tracks[j] = append(route.tracks[j], &Segment{
					pointA:  Point{X: protoSegment.PointA.Lat, Y: protoSegment.PointA.Lon},
					pointB:  Point{X: protoSegment.PointB.Lat, Y: protoSegment.PointB.Lon},
					dist:    protoSegment.Distance,
					bearing: protoSegment.Bearing,
					rel:     int(protoSegment.Rel),
				})
			}
		}
		n.routes = append(n.routes, &route)
	}
	n.routeIndex = int(nav.RouteIndex)
	n.trackIndex = int(nav.TrackIndex)
	n.segmentIndex = int(nav.SegmentIndex)
	n.segmentDistance = nav.SegmentDistance
	n.currentDistance = nav.CurrentDistance
	n.offlineIndex = int(nav.OfflineIndex)
	n.point = Point{X: nav.Point.Lat, Y: nav.Point.Lon}
	n.location = &Location{}
	n.elevation = new(types.Sensor)
	n.elevation.FromProto(nav.Elevation)
	if !nav.SkipOffline {
		n.offline = types.NewRandom(int(nav.OfflineMin), int(nav.OfflineMax))
	}
	n.totalDist = nav.TotalDistance
	n.skipOffline = nav.SkipOffline
}

func (n *Navigator) AddRoutes(routes ...*Route) {
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

func (n *Navigator) AllRoutes() []*Route {
	routes := make([]*Route, len(n.routes))
	copy(routes, n.routes)
	return routes
}

func (n *Navigator) TotalRoutes() int {
	return len(n.routes)
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
	if n.skipOffline {
		return
	}
	if n.offlineIndex > 0 {
		n.offlineIndex--
	}
}

func (n *Navigator) NextSensors(t float64) {
	n.elevation.Next(t)
}

func (n *Navigator) ToOffline() {
	if n.skipOffline {
		return
	}
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
					n.ToOffline()
					break loop
				}
				continue
			case notFoundSegment:
				if !n.nextTrack() {
					n.nextRoute()
				}
				n.ToOffline()
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
