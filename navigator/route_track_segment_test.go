package navigator

import (
	"testing"

	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/stretchr/testify/require"
)

func TestSegment_Snapshot(t *testing.T) {
	segment := track1km2segment.SegmentAt(0)
	snapshot := segment.Snapshot()
	require.NotNil(t, snapshot)
	require.Equal(t, segment.pointA.Lat, snapshot.PointA.Lat)
	require.Equal(t, segment.pointA.Lon, snapshot.PointA.Lon)
	require.Equal(t, segment.pointB.Lon, snapshot.PointB.Lon)
	require.Equal(t, segment.pointB.Lon, snapshot.PointB.Lon)
	require.Equal(t, segment.dist, snapshot.Distance)
	require.Equal(t, segment.bearing, snapshot.Bearing)
	require.Equal(t, segment.index, int(snapshot.Index))
	require.Equal(t, segment.rel, int(snapshot.Rel))
}

func TestSegment_FromSnapshot(t *testing.T) {
	snapshot := track1km2segment.SegmentAt(0).Snapshot()
	segment := SegmentFromSnapshot(snapshot)
	require.False(t, segment.IsEmpty())
	require.Equal(t, snapshot.PointA.Lat, segment.pointA.Lat)
	require.Equal(t, snapshot.PointA.Lon, segment.pointA.Lon)
	require.Equal(t, snapshot.PointB.Lon, segment.pointB.Lon)
	require.Equal(t, snapshot.PointB.Lon, segment.pointB.Lon)
	require.Equal(t, snapshot.Distance, segment.dist)
	require.Equal(t, snapshot.Bearing, segment.bearing)
	require.Equal(t, int(snapshot.Index), segment.index)
	require.Equal(t, int(snapshot.Rel), segment.rel)
}

func Test_makeSegments(t *testing.T) {
	type args struct {
		points []geo.LatLonPoint
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "should return error when there are less than two points",
			args: args{
				points: []geo.LatLonPoint{
					{Lon: 106.49403901121315, Lat: 29.534639667942272},
				},
			},
			wantErr: true,
		},
		{
			name: "should be 1 segment from 2 points - line",
			want: 1,
			args: args{
				points: []geo.LatLonPoint{
					{Lat: 106.49738017419708, Lon: 29.530081381358315},
					{Lat: 106.49854289891545, Lon: 29.531139573340525},
				},
			},
		},
		{
			name: "should be 2 segments from 3 points - line",
			want: 2,
			args: args{
				points: []geo.LatLonPoint{
					{Lon: 106.49954524781106, Lat: 29.53053489356381},
					{Lon: 106.49870327473838, Lat: 29.528546401884455},
					{Lon: 106.49664511834084, Lat: 29.52830219845653},
				},
			},
		},
		{
			name: "should be 4 segments from 5 points - polygon closed",
			want: 4,
			args: args{
				points: []geo.LatLonPoint{
					{Lon: 106.49403901121315, Lat: 29.534639667942272},
					{Lon: 106.49403901121315, Lat: 29.53170936443685},
					{Lon: 106.4996254357207, Lat: 29.53170936443685},
					{Lon: 106.4996254357207, Lat: 29.534639667942272},
					{Lon: 106.49403901121315, Lat: 29.534639667942272},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, _, err := makeSegments(tt.args.points)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeSegments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(segments) != tt.want {
				t.Errorf("makeSegments() segments = %v, want %v", len(segments), tt.want)
			}
			for i := 0; i < len(segments); i++ {
				segment := segments[i]
				require.NotZero(t, segment.PointA().Lon)
				require.NotZero(t, segment.PointA().Lat)
				require.NotZero(t, segment.PointB().Lon)
				require.NotZero(t, segment.PointB().Lat)
				require.False(t, segment.IsEmpty())
				require.NotZero(t, segment.Distance())
				require.Equal(t, i, segment.Index())
			}
		})
	}
}
