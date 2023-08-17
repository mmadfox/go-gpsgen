package navigator

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/properties"
	"github.com/mmadfox/go-gpsgen/proto"
)

// Track represents a track composed of segments.
type Track struct {
	id       string
	segments []Segment
	color    string
	name     Name
	props    properties.Properties
	dist     float64
	isClosed bool
	version  int
}

// TrackSnapshot is a serialized snapshot of a Track.
type TrackSnapshot = proto.Snapshot_Navigator_Route_Track

// NewTrack creates a new Track instance from a list of geographical points.
func NewTrack(points []geo.LatLonPoint) (*Track, error) {
	segments, dist, err := makeSegments(points)
	if err != nil {
		return nil, err
	}
	return &Track{
		id:       uuid.NewString(),
		segments: segments,
		color:    colorful.FastHappyColor().Hex(),
		dist:     dist,
		isClosed: isClosed(points),
	}, nil
}

// RestoreTrack restores a track from a snapshot.
func RestoreTrack(trackID, color string, points []geo.LatLonPoint) (*Track, error) {
	segments, dist, err := makeSegments(points)
	if err != nil {
		return nil, err
	}
	return &Track{
		id:       trackID,
		segments: segments,
		dist:     dist,
		isClosed: isClosed(points),
		color:    color,
	}, nil
}

// ID returns the ID of the track.
func (t *Track) ID() string {
	return t.id
}

// IsClosed checks if the track is closed.
func (t *Track) IsClosed() bool {
	return t.isClosed
}

// Color returns the color of the track.
func (t *Track) Color() string {
	return t.color
}

// Name returns the name of the track.
func (t *Track) Name() Name {
	return t.name
}

// ChangeColor changes the color of the track.
func (t *Track) ChangeColor(color colorful.Color) error {
	if !color.IsValid() {
		return fmt.Errorf("invalid track color %s", color.Hex())
	}
	t.color = color.Hex()
	t.nextVersion()
	return nil
}

// ChangeName changes the name of the track.
func (t *Track) ChangeName(name string) error {
	n, err := ParseName(name)
	if err != nil {
		return err
	}
	t.name = n
	t.nextVersion()
	return nil
}

// Distance returns the distance of the track.
func (t *Track) Distance() float64 {
	return t.dist
}

// Props returns the properties of the track.
func (t *Track) Props() properties.Properties {
	if t.props == nil {
		t.props = properties.Make()
	}
	return t.props
}

// NumSegments returns the number of segments in the track.
func (t *Track) NumSegments() int {
	return len(t.segments)
}

// SegmentAt returns the segment at the specified index.
func (t *Track) SegmentAt(i int) Segment {
	if len(t.segments) == 0 || i > len(t.segments)-1 || i < 0 {
		return Segment{}
	}
	return t.segments[i]
}

// Snapshot returns a serialized snapshot of the track.
func (t *Track) Snapshot() *TrackSnapshot {
	segments := make([]*SegmentSnapshot, len(t.segments))
	for i := 0; i < len(t.segments); i++ {
		segments[i] = t.segments[i].Snapshot()
	}
	var rawProps []byte
	if t.props != nil {
		rawProps, _ = json.Marshal(t.props)
	}
	return &TrackSnapshot{
		Id:        t.id,
		Segmenets: segments,
		Color:     t.color,
		Props:     rawProps,
		Distance:  t.dist,
		IsClosed:  t.isClosed,
		Name:      t.name.value,
		Version:   int64(t.version),
	}
}

func (t *Track) nextVersion() {
	t.version++
}

// TrackFromSnapshot restores a track from a snapshot.
func TrackFromSnapshot(track *Track, snap *TrackSnapshot) {
	if track == nil || snap == nil {
		return
	}
	var props properties.Properties
	if len(snap.Props) > 0 {
		props = properties.Make()
		_ = json.Unmarshal(snap.Props, &props)
	}
	track.id = snap.Id
	track.color = snap.Color
	track.props = props
	track.dist = snap.Distance
	track.isClosed = snap.IsClosed
	track.segments = make([]Segment, len(snap.Segmenets))
	for i := 0; i < len(snap.Segmenets); i++ {
		track.segments[i] = SegmentFromSnapshot(snap.Segmenets[i])
	}
	track.name.value = snap.Name
	track.version = int(snap.Version)
}

func isClosed(points []geo.LatLonPoint) bool {
	if len(points) < 3 {
		return false
	}
	p1 := points[0]
	p2 := points[len(points)-1]
	return p1.Lon == p2.Lon &&
		p2.Lat == p2.Lat
}
