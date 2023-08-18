package navigator

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mmadfox/go-gpsgen/properties"
	"github.com/mmadfox/go-gpsgen/proto"
)

// Route represents a route containing tracks and associated properties.
type Route struct {
	id      string
	dist    float64
	color   string
	name    Name
	tracks  []*Track
	props   properties.Properties
	version int
}

// RouteSnapshot is a serialized snapshot of a Route.
type RouteSnapshot = proto.Snapshot_Navigator_Route

// NewRoute creates a new Route instance with default values.
func NewRoute() *Route {
	return &Route{
		id:     uuid.NewString(),
		color:  colorful.FastHappyColor().Hex(),
		tracks: make([]*Track, 0),
	}
}

// RestoreRoute restores a Route instance with the given parameters.
func RestoreRoute(
	routeID string,
	color string,
	props properties.Properties,
) *Route {
	return &Route{
		id:     routeID,
		color:  color,
		props:  props,
		tracks: make([]*Track, 0),
	}
}

// RouteFromTracks creates a new Route instance from a list of tracks.
func RouteFromTracks(tracks ...*Track) *Route {
	newRoute := NewRoute()
	for i := 0; i < len(tracks); i++ {
		newRoute.AddTrack(tracks[i])
	}
	return newRoute
}

// ID returns the ID of the route.
func (r *Route) ID() string {
	return r.id
}

// Color returns the color of the route.
func (r *Route) Color() string {
	return r.color
}

// ChangeColor changes the color of the route.
func (r *Route) ChangeColor(color colorful.Color) error {
	if !color.IsValid() {
		return fmt.Errorf("invalid route color %s", color.Hex())
	}
	r.color = color.Hex()
	r.nextVersion()
	return nil
}

// Distance returns the total distance of the route.
func (r *Route) Distance() float64 {
	return r.dist
}

// Name returns the name of the route.
func (r *Route) Name() Name {
	return r.name
}

// ChangeName changes the name of the route.
func (r *Route) ChangeName(name string) error {
	n, err := ParseName(name)
	if err != nil {
		return err
	}
	r.name = n
	r.nextVersion()
	return nil
}

// Props returns the properties associated with the route.
func (r *Route) Props() properties.Properties {
	if r.props == nil {
		r.props = properties.Make()
	}
	return r.props
}

// NumTracks returns the number of tracks in the route.
func (r *Route) NumTracks() int {
	return len(r.tracks)
}

// TrackAt returns the track at the specified index.
func (r *Route) TrackAt(i int) *Track {
	if len(r.tracks) == 0 || i > len(r.tracks)-1 || i < 0 {
		return nil
	}
	return r.tracks[i]
}

// RemoveTrack removes a track from the route.
func (r *Route) RemoveTrack(trackID string) (ok bool) {
	if len(trackID) == 0 || len(trackID) > 36 {
		return
	}
	for i := 0; i < len(r.tracks); i++ {
		track := r.tracks[i]
		if track.ID() == trackID {
			r.tracks = append(r.tracks[:i], r.tracks[i+1:]...)
			ok = true
			track = nil
			break
		}
	}
	if ok {
		r.updateState()
		r.nextVersion()
	}
	return
}

// TrackByID returns the track with the specified ID.
func (r *Route) TrackByID(trackID string) (*Track, error) {
	if len(trackID) == 0 || len(trackID) > 36 {
		return nil, ErrTrackNotFound
	}
	for i := 0; i < len(r.tracks); i++ {
		if r.tracks[i].ID() == trackID {
			r.nextVersion()
			return r.tracks[i], nil
		}
	}
	return nil, ErrTrackNotFound
}

// EachTrack iterates through each track in the route and applies a function.
func (r *Route) EachTrack(fn func(int, *Track) bool) {
	for i := 0; i < len(r.tracks); i++ {
		if ok := fn(i, r.tracks[i]); !ok {
			break
		}
	}
}

// AddTrack adds a track to the route.
func (r *Route) AddTrack(track *Track) *Route {
	if track == nil {
		return r
	}

	r.tracks = append(r.tracks, track)
	r.updateState()
	r.nextVersion()

	return r
}

// Snapshot creates a snapshot of the route.
func (r *Route) Snapshot() *RouteSnapshot {
	tracks := make([]*proto.Snapshot_Navigator_Route_Track, len(r.tracks))
	for i := 0; i < len(r.tracks); i++ {
		tracks[i] = r.tracks[i].Snapshot()
	}
	var rawProto []byte
	if r.props != nil {
		rawProto, _ = json.Marshal(r.props)
	}
	return &RouteSnapshot{
		Id:       r.id,
		Distance: r.dist,
		Color:    r.color,
		Props:    rawProto,
		Tracks:   tracks,
		Name:     r.Name().String(),
		Version:  int64(r.version),
	}
}

// RouteFromSnapshot restores a route from a snapshot.
func (r *Route) RouteFromSnapshot(snap *RouteSnapshot) error {
	if snap == nil {
		return nil
	}
	tracks := make([]*Track, len(snap.Tracks))
	for i := 0; i < len(snap.Tracks); i++ {
		track := new(Track)
		TrackFromSnapshot(track, snap.Tracks[i])
		tracks[i] = track
	}
	var props properties.Properties
	if len(snap.Props) > 0 {
		props = properties.Make()
		if err := json.Unmarshal(snap.Props, &props); err != nil {
			return err
		}
	}

	if len(snap.Name) > 0 {
		name, err := ParseName(snap.Name)
		if err != nil {
			return err
		}
		r.name = name
	}

	r.id = snap.Id
	r.dist = snap.Distance
	r.color = snap.Color
	r.props = props
	r.tracks = tracks
	r.version = int(snap.Version)
	return nil
}

func (r *Route) updateState() {
	r.dist = 0
	for i := 0; i < len(r.tracks); i++ {
		r.dist += r.tracks[i].dist
	}
}

func (r *Route) indexTrack(trackID string) int {
	for i := 0; i < len(r.tracks); i++ {
		if r.tracks[i].id == trackID {
			return i
		}
	}
	return -1
}

func (r *Route) nextVersion() {
	r.version++
}
