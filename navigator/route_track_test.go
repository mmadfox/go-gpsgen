package navigator

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/stretchr/testify/require"
)

func TestTrack_IsClosed(t *testing.T) {
	type args struct {
		points []geo.LatLonPoint
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "should return error when no points",
			args: args{
				points: nil,
			},
			wantErr: true,
		},
		{
			name: "should return false when track is line",
			args: args{
				points: []geo.LatLonPoint{
					{Lon: 106.49659165973264, Lat: 29.532430319999506},
					{Lon: 106.49863645147855, Lat: 29.53127911431376},
					{Lon: 106.49775438445101, Lat: 29.530488379584682},
				},
			},
			want: false,
		},
		{
			name: "should return true when geometry is closed",
			args: args{
				points: []geo.LatLonPoint{
					{Lon: 106.49803504214191, Lat: 29.532430319999506},
					{Lon: 106.49803504214191, Lat: 29.531837276283085},
					{Lon: 106.49803504214191, Lat: 29.531837276283085},
					{Lon: 106.49907748499231, Lat: 29.531837276283085},
					{Lon: 106.49907748499231, Lat: 29.532430319999506},
					{Lon: 106.49803504214191, Lat: 29.532430319999506},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := NewTrack(tt.args.points)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTrack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got := tr.IsClosed(); got != tt.want {
				t.Errorf("Track.IsPolygon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack_Color(t *testing.T) {
	points := []geo.LatLonPoint{
		{Lon: 106.49659165973264, Lat: 29.532430319999506},
		{Lon: 106.49863645147855, Lat: 29.53127911431376},
		{Lon: 106.49775438445101, Lat: 29.530488379584682},
	}

	track, err := NewTrack(points)
	require.NotEmpty(t, track.Color())
	require.NoError(t, err)
	require.NotNil(t, track)
}

func TestTrack_Snapshot(t *testing.T) {
	track, err := NewTrack([]geo.LatLonPoint{
		{Lon: 106.49331396675268, Lat: 29.5299004724652},
		{Lon: 106.49523863664103, Lat: 29.532016484207674},
	})
	require.NoError(t, err)
	track.Props().Set("foo", 1)
	track.Props().Set("bar", "baz")
	snapshot := track.Snapshot()
	require.NotNil(t, snapshot)
	require.Equal(t, track.ID(), snapshot.Id)
	require.Equal(t, track.Distance(), snapshot.Distance)
	require.Equal(t, track.Color(), snapshot.Color)
	require.Equal(t, track.IsClosed(), snapshot.IsClosed)
	require.Len(t, snapshot.Segmenets, track300m1segment.NumSegments())
	require.NotEmpty(t, snapshot.Props)
}

func TestTrack_FromSnapshot(t *testing.T) {
	track, err := NewTrack([]geo.LatLonPoint{
		{Lon: 106.49331396675268, Lat: 29.5299004724652},
		{Lon: 106.49523863664103, Lat: 29.532016484207674},
	})
	require.NoError(t, err)
	track.Props().Set("foo", 1)
	track.Props().Set("bar", "baz")
	snapshot := track.Snapshot()
	TrackFromSnapshot(track, snapshot)
	require.Equal(t, snapshot.Id, track.ID())
	require.Equal(t, snapshot.Distance, track.Distance())
	require.Equal(t, snapshot.Color, track.Color())
	require.Equal(t, snapshot.IsClosed, track.IsClosed())
	require.Equal(t, len(snapshot.Segmenets), track.NumSegments())
	require.NotEmpty(t, track.Props())
}
