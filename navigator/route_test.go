package navigator

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewRoute(t *testing.T) {
	route := NewRoute()
	require.NotNil(t, route)
	require.NotEmpty(t, route.Color())
	require.NotEmpty(t, route.ID())
	require.Zero(t, route.Distance())
	require.Empty(t, route.Props())
	require.Zero(t, route.NumTracks())
}

func TestRoute_Props(t *testing.T) {
	route := NewRoute()
	route.Props().Set("key", "val")
	prop1, ok := route.Props().String("key")
	require.True(t, ok)
	require.Equal(t, "val", prop1)
}

func TestRoute_RemoveTrack(t *testing.T) {
	type args struct {
		trackID string
	}
	tests := []struct {
		name    string
		args    args
		arrange func(*Route)
		assert  func(*Route)
		wantOk  bool
	}{
		{
			name: "should return false when trackID is empty",
			args: args{
				trackID: "",
			},
			wantOk: false,
		},
		{
			name: "should return false when track not found",
			args: args{
				trackID: "someid",
			},
			arrange: func(r *Route) {
				r.AddTrack(track1km2segment).AddTrack(track300m1segment)
			},
			wantOk: false,
		},
		{
			name: "should return false when invalid trackID",
			args: args{
				trackID: uuid.NewString() + "a",
			},
			wantOk: false,
		},
		{
			name: "should return true when track removed",
			args: args{
				trackID: track1km2segment.ID(),
			},
			arrange: func(r *Route) {
				r.AddTrack(track1km2segment)
			},
			assert: func(r *Route) {
				require.Zero(t, r.NumTracks())
				require.Zero(t, r.Distance())
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := NewRoute()
			if tt.arrange != nil {
				tt.arrange(route)
			}
			if gotOk := route.RemoveTrack(tt.args.trackID); gotOk != tt.wantOk {
				t.Errorf("Route.RemoveTrack() = %v, want %v", gotOk, tt.wantOk)
			}
			if tt.assert != nil {
				tt.assert(route)
			}
		})
	}
}

func TestRoute_TrackByID(t *testing.T) {
	type args struct {
		trackID string
	}
	tests := []struct {
		name    string
		args    args
		arrange func(*Route)
		want    *Track
		wantErr bool
	}{
		{
			name: "should return error when trackID is empty",
			args: args{
				trackID: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when trackID greater than 36 chars",
			args: args{
				trackID: track1km2segment.ID() + "a",
			},
			wantErr: true,
		},
		{
			name: "should return error when track not found",
			args: args{
				trackID: uuid.NewString(),
			},
			arrange: func(r *Route) {
				r.AddTrack(track1km2segment)
			},
			wantErr: true,
		},
		{
			name: "should return track",
			args: args{
				trackID: track300m1segment.ID(),
			},
			arrange: func(r *Route) {
				r.AddTrack(track300m1segment)
			},
			want: track300m1segment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRoute()
			if tt.arrange != nil {
				tt.arrange(r)
			}
			got, err := r.TrackByID(tt.args.trackID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.TrackByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.TrackByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_EachRoute(t *testing.T) {
	route := NewRoute()
	route.AddTrack(track1km2segment).AddTrack(track3km7segments)

	actualTracks := 0
	route.EachTrack(func(n int, track *Track) bool {
		require.NotNil(t, track)
		actualTracks++
		return true
	})
	require.Equal(t, route.NumTracks(), actualTracks)

	actualTracks = 0
	route.EachTrack(func(n int, track *Track) bool {
		if n > 0 {
			return false
		}
		require.NotNil(t, track)
		actualTracks++
		return true
	})
	require.Equal(t, 1, actualTracks)
}

func TestRoute_AddTrack(t *testing.T) {
	type args struct {
		track *Track
	}
	tests := []struct {
		name   string
		args   args
		assert func(*Route)
	}{
		{
			name: "should ignored track when track is nil",
			args: args{
				track: nil,
			},
			assert: func(r *Route) {
				require.Zero(t, r.NumTracks())
			},
		},
		{
			name: "should added track",
			args: args{
				track: track1km2segment,
			},
			assert: func(r *Route) {
				require.Equal(t, 1, r.NumTracks())
				require.NotZero(t, r.Distance())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRoute()
			r.AddTrack(tt.args.track)
			if tt.assert != nil {
				tt.assert(r)
			}
		})
	}
}

func TestRoute_Snapshot(t *testing.T) {
	route := NewRoute().
		AddTrack(track1km2segment).
		AddTrack(track300m1segment)
	route.Props().Set("k", "v")
	snapshot := route.Snapshot()
	require.NotNil(t, snapshot)
	require.Equal(t, route.ID(), snapshot.Id)
	require.Equal(t, route.Distance(), snapshot.Distance)
	require.Equal(t, route.Color(), snapshot.Color)
	require.Len(t, snapshot.Tracks, route.NumTracks())
	require.NotEmpty(t, snapshot.Props)
}

func TestRoute_FromSnapshot(t *testing.T) {
	expectedName := "some name"
	r1 := NewRoute().
		AddTrack(track1km2segment).
		AddTrack(track300m1segment)
	r1.Props().Set("k", "v")
	require.NoError(t, r1.ChangeName(expectedName))
	snap := r1.Snapshot()

	route := new(Route)
	err := route.RouteFromSnapshot(snap)
	require.NoError(t, err)
	require.Equal(t, snap.Id, route.ID())
	require.Equal(t, snap.Distance, route.Distance())
	require.Equal(t, snap.Color, route.Color())
	require.Len(t, snap.Tracks, route.NumTracks())
	require.NotEmpty(t, route.Props())
	require.Equal(t, expectedName, route.Name().String())
}
