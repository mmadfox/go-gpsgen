package navigator

import (
	"math"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/types"
)

// Navigator manages navigation operations and state.
type Navigator struct {
	routes                 []*Route
	routeIndex             int
	trackIndex             int
	segmentIndex           int
	currentSegmentDistance float64
	currentRouteDistance   float64
	currentTrackDistance   float64
	currentDistance        float64
	offlineIndex           int
	point                  geo.LatLonPoint
	elevation              *types.Sensor
	offline                *types.Random
	distance               float64
	skipOffline            bool
	version                int
}

// New creates a new instance of Navigator with the provided options.
func New(opts ...Option) (*Navigator, error) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}

	if err := o.validate(); err != nil {
		return nil, err
	}

	elevation, err := types.NewSensor("elevation",
		o.elevationMin,
		o.elevationMax,
		o.elevationAmplitude,
		o.elevationMode,
	)
	if err != nil {
		return nil, err
	}

	nav := &Navigator{
		routes:      make([]*Route, 0),
		elevation:   elevation,
		skipOffline: o.skipOffline,
	}
	switch o.skipOffline {
	case false:
		nav.offline = types.NewRandom(o.minOffline, o.maxOffline)
	}

	return nav, nil
}

// Location returns the current geographic location of the navigator.
func (n *Navigator) Location() geo.LatLonPoint {
	return n.point
}

// CurrentBearing returns the current bearing of the navigator.
func (n *Navigator) CurrentBearing() float64 {
	segment := n.CurrentSegment()
	if segment.IsEmpty() {
		return 0
	}
	return segment.Bearing()
}

// AddRoute adds one or more routes to the navigator.
func (n *Navigator) AddRoute(routes ...*Route) error {
	if len(routes) == 0 {
		return ErrNoRoutes
	}
	ok := 0
	for i := 0; i < len(routes); i++ {
		if routes[i] == nil || n.hasRoute(routes[i].id) {
			continue
		}
		n.routes = append(n.routes, routes[i])
		ok++
	}
	if ok > 0 {
		n.Reset()
		n.nextVersion()
	}
	return nil
}

// Sum calculates and returns a summary of versions for the navigator, routes, and tracks.
func (n *Navigator) Sum() [3]int {
	a, b := 0, 0
	for i := 0; i < len(n.routes); i++ {
		route := n.routes[i]
		a += route.version
		for j := 0; j < len(route.tracks); j++ {
			track := route.tracks[j]
			b += track.version
		}
	}
	return [3]int{n.version, a, b}
}

// ResetRoutes resets all routes in the navigator.
func (n *Navigator) ResetRoutes() bool {
	if len(n.routes) == 0 {
		return false
	}
	n.routes = make([]*Route, 0)
	n.Reset()
	n.nextVersion()
	return true
}

// RouteByID retrieves a route by its ID.
func (n *Navigator) RouteByID(routeID string) (*Route, error) {
	if len(routeID) == 0 || len(routeID) > 36 {
		return nil, ErrRouteNotFound
	}
	for i := 0; i < len(n.routes); i++ {
		if n.routes[i].ID() == routeID {
			return n.routes[i], nil
		}
	}
	return nil, ErrRouteNotFound
}

// EachRoute iterates through each route and invokes the provided function.
func (n *Navigator) EachRoute(fn func(int, *Route) bool) {
	if len(n.routes) == 0 {
		return
	}
	for i := 0; i < len(n.routes); i++ {
		if ok := fn(i, n.routes[i]); !ok {
			return
		}
	}
}

// RemoveRoute removes a route by its ID.
func (n *Navigator) RemoveRoute(routeID string) (ok bool) {
	if len(routeID) == 0 || len(routeID) > 36 {
		return
	}
	for i := 0; i < len(n.routes); i++ {
		if n.routes[i].ID() == routeID {
			n.routes = append(n.routes[:i], n.routes[i+1:]...)
			ok = true
			break
		}
	}
	if ok {
		n.nextVersion()
		n.Reset()
	}
	return
}

// RemoveTrack removes a track from a route.
func (n *Navigator) RemoveTrack(routeID string, trackID string) bool {
	if len(routeID) == 0 ||
		len(routeID) > 36 ||
		len(trackID) == 0 ||
		len(trackID) > 36 {
		return false
	}
	route, err := n.RouteByID(routeID)
	if err != nil {
		return false
	}
	ok := route.RemoveTrack(trackID)
	if ok {
		n.nextVersion()
	}
	return ok
}

// Distance returns the total distance of the navigator's routes.
func (n *Navigator) Distance() float64 {
	return n.distance
}

// CurrentDistance returns the current accumulated distance.
func (n *Navigator) CurrentDistance() float64 {
	return n.currentDistance
}

// RouteDistance returns the distance of the current route.
func (n *Navigator) RouteDistance() float64 {
	route := n.CurrentRoute()
	if route == nil {
		return 0
	}
	return route.Distance()
}

// CurrentRouteDistance returns the current accumulated distance of the current route.
func (n *Navigator) CurrentRouteDistance() float64 {
	return n.currentRouteDistance
}

// TrackDistance returns the distance of the current track.
func (n *Navigator) TrackDistance() float64 {
	track := n.CurrentTrack()
	if track == nil {
		return 0
	}
	return track.Distance()
}

// CurrentTrackDistance returns the current accumulated distance of the current track.
func (n *Navigator) CurrentTrackDistance() float64 {
	return n.currentTrackDistance
}

// SegmentDistance returns the distance of the current segment.
func (n *Navigator) SegmentDistance() float64 {
	segment := n.CurrentSegment()
	if segment.IsEmpty() {
		return 0
	}
	return segment.Distance()
}

// CurrentSegmentDistance returns the current accumulated distance of the current segment.
func (n *Navigator) CurrentSegmentDistance() float64 {
	return n.currentSegmentDistance
}

// IsFinish checks if the navigator has finished its navigation path.
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

// Routes returns a copy of the navigator's routes slice.
func (n *Navigator) Routes() []*Route {
	routes := make([]*Route, len(n.routes))
	copy(routes, n.routes)
	return routes
}

// CurrentRoute returns the current route.
func (n *Navigator) CurrentRoute() *Route {
	if len(n.routes) == 0 {
		return nil
	}
	return n.routes[n.routeIndex]
}

// CurrentTrack returns the current track.
func (n *Navigator) CurrentTrack() *Track {
	route := n.CurrentRoute()
	if route == nil {
		return nil
	}
	return route.TrackAt(n.trackIndex)
}

// CurrentSegment returns the current segment.
func (n *Navigator) CurrentSegment() Segment {
	route := n.CurrentRoute()
	if route == nil {
		return Segment{}
	}
	track := route.TrackAt(n.trackIndex)
	if track == nil {
		return Segment{}
	}
	segment := track.SegmentAt(n.segmentIndex)
	if segment.IsEmpty() {
		return Segment{}
	}
	return segment
}

// NumRoutes returns the number of routes in the navigator.
func (n *Navigator) NumRoutes() int {
	return len(n.routes)
}

// RouteAt returns the route at the specified index.
func (n *Navigator) RouteAt(index int) *Route {
	if len(n.routes) == 0 || (index > len(n.routes)-1 || index < 0) {
		return nil
	}
	return n.routes[index]
}

// TrackIndex returns the index of the current track.
func (n *Navigator) TrackIndex() int {
	return n.trackIndex
}

// RouteIndex returns the index of the current route.
func (n *Navigator) RouteIndex() int {
	return n.routeIndex
}

// SegmentIndex returns the index of the current segment.
func (n *Navigator) SegmentIndex() int {
	return n.segmentIndex
}

// Elevation returns the current elevation.
func (n *Navigator) Elevation() float64 {
	return n.elevation.ValueY()
}

// IsOffline checks if the navigator is in offline mode.
func (n *Navigator) IsOffline() bool {
	return n.offlineIndex > 0
}

// OfflineDuration returns the remaining offline duration.
func (n *Navigator) OfflineDuration() int {
	return n.offlineIndex
}

// Update updates the navigator's state in the provided Device proto message.
func (n *Navigator) Update(state *proto.Device) {
	n.normalizeDeviceState(state)

	state.Distance.Distance = n.distance
	state.Distance.CurrentDistance = n.currentDistance
	state.Distance.RouteDistance = n.RouteDistance()
	state.Distance.CurrentRouteDistance = n.currentRouteDistance
	state.Distance.TrackDistance = n.TrackDistance()
	state.Distance.CurrentTrackDistance = n.currentTrackDistance
	state.Distance.SegmentDistance = n.SegmentDistance()
	state.Distance.CurrentSegmentDistance = n.currentSegmentDistance

	state.Navigator.CurrentRouteIndex = int64(n.routeIndex)
	state.Navigator.CurrentTrackIndex = int64(n.trackIndex)
	state.Navigator.CurrentSegmentIndex = int64(n.segmentIndex)
	if route := n.CurrentRoute(); route != nil {
		state.Navigator.CurrentRouteId = route.ID()
	}
	if track := n.CurrentTrack(); track != nil {
		state.Navigator.CurrentTrackId = track.ID()
	}

	state.Location.Bearing = n.CurrentBearing()
	state.Location.Elevation = n.elevation.ValueY()
	state.Location.Lat = n.point.Lat
	state.Location.Lon = n.point.Lon

	state.IsOffline = n.offlineIndex > 0
	state.OfflineDuration = int64(n.offlineIndex)

	setUTM(n.point.Lat, n.point.Lon, state.Location.Utm)
	setDMS(n.point.Lat, n.point.Lon, state.Location.LatDms, state.Location.LonDms)
}

// NextRoute moves to the next route.
func (n *Navigator) NextRoute() bool {
	next := n.routeIndex < len(n.routes)-1
	if !next {
		return false
	}
	nextIndex := n.routeIndex + 1
	return n.MoveToRoute(nextIndex)
}

// PrevRoute moves to the previous route.
func (n *Navigator) PrevRoute() bool {
	prev := n.routeIndex > 0
	if !prev {
		return false
	}
	prevIndex := n.routeIndex - 1
	return n.MoveToRoute(prevIndex)
}

// MoveToRouteByID moves to a route specified by its ID.
func (n *Navigator) MoveToRouteByID(routeID string) bool {
	index := n.indexRoute(routeID)
	if index == -1 {
		return false
	}
	return n.MoveToRoute(index)
}

// MoveToRouteByID moves to a route specified by its ID.
func (n *Navigator) MoveToTrackByID(routeID string, trackID string) bool {
	routeIndex := n.indexRoute(routeID)
	if routeIndex == -1 {
		return false
	}
	route := n.routes[routeIndex]
	trackIndex := route.indexTrack(trackID)
	if trackIndex == -1 {
		return false
	}
	return n.MoveToTrack(routeIndex, trackIndex)
}

// MoveToRoute moves to a route specified by its index.
func (n *Navigator) MoveToRoute(routeIndex int) bool {
	return n.MoveToSegment(routeIndex, 0, 0)
}

// MoveToTrack moves to a track within a route specified by indices.
func (n *Navigator) MoveToTrack(routeIndex int, trackIndex int) bool {
	return n.MoveToSegment(routeIndex, trackIndex, 0)
}

// MoveToSegment moves to a segment within a track within a route specified by indices.
func (n *Navigator) MoveToSegment(routeIndex int, trackIndex int, segmentIndex int) bool {
	if len(n.routes) == 0 {
		return false
	}

	if n.routeIndex == routeIndex &&
		n.trackIndex == trackIndex &&
		n.segmentIndex == segmentIndex {
		return false
	}

	if routeIndex < 0 || routeIndex > len(n.routes)-1 {
		return false
	}
	route := n.routes[routeIndex]

	if len(route.tracks) == 0 {
		return false
	}
	if trackIndex < 0 || trackIndex > len(route.tracks)-1 {
		return false
	}
	track := route.TrackAt(trackIndex)

	if segmentIndex < 0 || segmentIndex > len(track.segments)-1 {
		return false
	}
	segment := track.SegmentAt(segmentIndex)
	if segmentIndex != n.segmentIndex {
		n.currentSegmentDistance = 0
	}

	if trackIndex != n.trackIndex {
		n.currentTrackDistance = 0
		n.currentSegmentDistance = 0
	}

	if routeIndex != n.routeIndex {
		n.currentRouteDistance = 0
		n.currentTrackDistance = 0
		n.currentSegmentDistance = 0
	}

	n.point = segment.pointA
	n.routeIndex = routeIndex
	n.trackIndex = trackIndex
	n.segmentIndex = segmentIndex
	n.offlineIndex = 0
	return true
}

// DestinationTo moves the navigator to a destination specified by distance.
func (n *Navigator) DestinationTo(meters float64) bool {
	return n.jump(meters)
}

// ToOffline puts the navigator into offline mode.
func (n *Navigator) ToOffline() {
	if n.skipOffline {
		return
	}
	n.offlineIndex = n.offline.Value()
}

// NextElevation updates the next elevation.
func (n *Navigator) NextElevation(tick float64) {
	n.elevation.Next(tick)
}

// ShuffleElevation shuffles the generator of the Elevation instance.
func (n *Navigator) ShuffleElevation() {
	n.elevation.Shuffle()
}

// NextLocation updates the navigator's location based on the provided tick and speed.
func (n *Navigator) NextLocation(tick float64, speed float64) (ok bool) {
	if len(n.routes) == 0 {
		return
	}

	if ok := n.nextOffline(); ok {
		return false
	}

	newDist := tick * speed

loop:
	for {
		if n.isValidDistance(newDist) {
			isStartSeg := n.isStartSegment()
			n.updateDist(newDist)
			segment := n.CurrentSegment()
			if !segment.IsEmpty() {
				if isStartSeg {
					n.point = segment.DestinationTo(1)
				} else {
					n.point = segment.DestinationTo(n.currentSegmentDistance)
				}
				ok = true
			}
			break
		} else {
			seg := n.nextSegment()
			foundSegment := !seg.IsEmpty()
			notFoundSegment := seg.IsEmpty()
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

// Snapshot creates a snapshot of the navigator's state.
func (n *Navigator) Snapshot() *proto.Snapshot_Navigator {
	routes := make([]*proto.Snapshot_Navigator_Route, len(n.routes))
	for i := 0; i < len(n.routes); i++ {
		routes[i] = n.routes[i].Snapshot()
	}
	snap := &proto.Snapshot_Navigator{
		Routes:                 routes,
		RouteIndex:             int64(n.routeIndex),
		TrackIndex:             int64(n.trackIndex),
		SegmentIndex:           int64(n.segmentIndex),
		CurrentSegmentDistance: n.currentSegmentDistance,
		CurrentRouteDistance:   n.currentRouteDistance,
		CurrentTrackDistance:   n.currentTrackDistance,
		CurrentDistance:        n.currentDistance,
		OfflineIndex:           int64(n.offlineIndex),
		Point:                  &proto.Snapshot_PointLatLon{Lon: n.point.Lon, Lat: n.point.Lat},
		Elevation:              n.elevation.Snapshot(),
		Distance:               n.distance,
		SkipOffline:            n.skipOffline,
		Version:                int64(n.version),
	}
	if n.offline != nil {
		snap.OfflineMin = int64(n.offline.Min())
		snap.OfflineMax = int64(n.offline.Max())
	}
	return snap
}

// FromSnapshot restores the navigator's state from a snapshot.
func (n *Navigator) FromSnapshot(snap *proto.Snapshot_Navigator) {
	if snap == nil {
		return
	}
	n.routes = make([]*Route, len(snap.Routes))
	for i := 0; i < len(snap.Routes); i++ {
		route := new(Route)
		route.RouteFromSnapshot(snap.Routes[i])
		n.routes[i] = route
	}
	n.routeIndex = int(snap.RouteIndex)
	n.trackIndex = int(snap.TrackIndex)
	n.segmentIndex = int(snap.SegmentIndex)
	n.currentSegmentDistance = snap.CurrentSegmentDistance
	n.currentRouteDistance = snap.CurrentRouteDistance
	n.currentTrackDistance = snap.CurrentTrackDistance
	n.currentDistance = snap.CurrentDistance
	n.offlineIndex = int(snap.OfflineIndex)
	n.point = geo.LatLonPoint{
		Lon: snap.Point.Lon,
		Lat: snap.Point.Lat,
	}
	n.elevation = new(types.Sensor)
	n.elevation.FromSnapshot(snap.Elevation)
	if !n.skipOffline {
		n.offline = types.NewRandom(int(snap.OfflineMin), int(snap.OfflineMax))
	}
	n.distance = snap.Distance
	n.skipOffline = snap.SkipOffline
	n.version = int(snap.Version)
}

func (n *Navigator) jump(distance float64) bool {
	if distance > n.distance {
		distance = n.distance
	}
	if distance < 1 {
		return false
	}

	var (
		route        *Route
		track        *Track
		segment      Segment
		trackIndex   int
		routeIndex   int
		segmentIndex int
		routeDist    float64
		curRouteDist float64
		curTrackDist float64
	)

	for i := 0; i < len(n.routes); i++ {
		route = n.routes[i]
		routeIndex += i
		if route.dist+routeDist > distance {
			break
		}
		routeDist += route.dist
	}
	if route == nil {
		return false
	}

	for i := 0; i < route.NumTracks(); i++ {
		track = route.TrackAt(i)
		if track == nil {
			continue
		}
		trackIndex += i
		if track.dist+routeDist > distance {
			break
		}
		curRouteDist += track.dist
		routeDist += track.dist
	}
	if track == nil {
		return false
	}

	for i := 0; i < track.NumSegments(); i++ {
		segment = track.SegmentAt(i)
		if segment.IsEmpty() {
			continue
		}
		segmentIndex = i
		if segment.dist+routeDist > distance {
			break
		}
		routeDist += segment.dist
		curRouteDist += segment.dist
		curTrackDist += segment.dist
	}
	if segment.IsEmpty() {
		return false
	}

	diff := math.Abs(distance - routeDist)
	n.point = segment.DestinationTo(diff)
	n.routeIndex = routeIndex
	n.trackIndex = trackIndex
	n.segmentIndex = segmentIndex
	n.currentSegmentDistance = diff
	n.currentDistance = routeDist + diff
	n.currentRouteDistance = curRouteDist + diff
	n.currentTrackDistance = curTrackDist + diff
	return true
}

// Reset resets the Navigator's state, including indexes and calculated route distance.
func (n *Navigator) Reset() {
	n.resetIndexes()
	n.calcRouteDistance()
}

func (n *Navigator) updateDist(newDist float64) {
	n.currentDistance += newDist
	n.currentSegmentDistance += newDist
	n.currentTrackDistance += newDist
	n.currentRouteDistance += newDist
}

func (n *Navigator) isValidDistance(dist float64) bool {
	segment := n.CurrentSegment()
	if segment.IsEmpty() {
		return false
	}
	return n.currentSegmentDistance+dist <= segment.Distance()
}

func (n *Navigator) isStartSegment() bool {
	return n.segmentIndex == 0 && n.currentSegmentDistance == 0
}

func (n *Navigator) nextOffline() bool {
	if n.skipOffline {
		return false
	}
	if n.offlineIndex > 0 {
		n.offlineIndex--
		return true
	}
	return false
}

func (n *Navigator) nextSegment() Segment {
	track := n.CurrentTrack()
	if track == nil {
		return Segment{}
	}
	if n.segmentIndex == track.NumSegments()-1 {
		return Segment{}
	}
	segment := n.CurrentSegment()
	n.currentSegmentDistance = 0
	n.segmentIndex++
	return segment
}

func (n *Navigator) nextTrack() bool {
	n.currentSegmentDistance = 0
	n.segmentIndex = 0
	route := n.CurrentRoute()
	if route == nil {
		return false
	}
	n.currentTrackDistance = 0
	if n.trackIndex < route.NumTracks()-1 {
		n.trackIndex++
		return true
	}
	n.trackIndex = 0
	return false
}

func (n *Navigator) nextRoute() {
	n.currentRouteDistance = 0
	if n.routeIndex < len(n.routes)-1 {
		n.routeIndex++
		return
	}
	n.routeIndex = 0
	n.currentDistance = 0.1
}

func (n *Navigator) calcRouteDistance() {
	n.distance = 0
	for i := 0; i < len(n.routes); i++ {
		n.distance += n.routes[i].dist
	}
}

func (n *Navigator) resetIndexes() {
	n.routeIndex = 0
	n.trackIndex = 0
	n.segmentIndex = 0
	n.currentSegmentDistance = 0
	n.currentRouteDistance = 0
	n.currentTrackDistance = 0
	n.currentDistance = 0
}

func (n *Navigator) hasRoute(id string) bool {
	for i := 0; i < len(n.routes); i++ {
		if id == n.routes[i].id {
			return true
		}
	}
	return false
}

func (n *Navigator) indexRoute(routeID string) int {
	if len(routeID) == 0 {
		return -1
	}
	for i := 0; i < len(n.routes); i++ {
		if n.routes[i].id == routeID {
			return i
		}
	}
	return -1
}

func (n *Navigator) nextVersion() {
	n.version++
}

func (n *Navigator) normalizeDeviceState(device *proto.Device) {
	if device.Distance == nil {
		device.Distance = new(proto.Device_Distance)
	}
	if device.Navigator == nil {
		device.Navigator = new(proto.Device_Navigator)
	}
	if device.Location == nil {
		device.Location = &proto.Device_Location{
			LatDms: new(proto.Device_Location_DMS),
			LonDms: new(proto.Device_Location_DMS),
			Utm:    new(proto.Device_Location_UTM),
		}
	}
}
