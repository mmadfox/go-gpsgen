package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpeed_New(t *testing.T) {
	type args struct {
		min       float64
		max       float64
		amplitude int
	}
	type want struct {
		min float64
		max float64
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "should return valid speed value when min and max eq 0",
			args: args{
				min:       0,
				max:       0,
				amplitude: minAmplitude,
			},
			want: want{
				min: 1,
				max: 1,
			},
		},
		{
			name: "should return error when min > max",
			args: args{
				min: 30,
				max: 10,
			},
			wantErr: true,
		},
		{
			name: "should return when max > MaxSpeed",
			args: args{
				min: 0,
				max: MaxSpeedVal + 1,
			},
			wantErr: true,
		},
		{
			name: "should return error when min less than the minimum value",
			args: args{
				min: MinSpeedVal - 1,
				max: 0,
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude less than min value",
			args: args{
				min:       0,
				max:       100,
				amplitude: 0,
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude greater than max value",
			args: args{
				min:       0,
				max:       100,
				amplitude: maxAmplitude + 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSpeed(tt.args.min, tt.args.max, tt.args.amplitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want.min, got.Min())
				require.Equal(t, tt.want.max, got.Max())
				require.NotZero(t, got.String())
			}
		})
	}
}

func TestSpeed_Next(t *testing.T) {
	speed, err := NewSpeed(1, 100, 32)
	require.NoError(t, err)

	speed.Shuffle()

	var prev float64
	for i := 0; i < 100; i++ {
		speed.Next(float64(i) / 100)
		if i < 2 {
			prev = speed.Value()
			continue
		}
		diff := math.Abs(speed.Value() - prev)
		require.NotZero(t, diff)
	}
}

func TestSpeed_Snapshot(t *testing.T) {
	speed, err := NewSpeed(1, 3, 4)
	require.NoError(t, err)
	require.NotNil(t, speed)
	snap := speed.Snapshot()
	require.Equal(t, speed.Min(), snap.Min)
	require.Equal(t, speed.Max(), snap.Max)
	require.Equal(t, speed.Value(), snap.Val)
	require.NotNil(t, snap.Gen)
}

func TestSpeed_FromSnapshot(t *testing.T) {
	speed, err := NewSpeed(1, 3, 4)
	speed.Next(0.1)
	require.NoError(t, err)
	require.NotNil(t, speed)
	snap := speed.Snapshot()
	speed2 := new(Speed)
	speed2.FromSnapshot(snap)
	require.Equal(t, speed.Min(), speed2.Min())
	require.Equal(t, speed.Max(), speed2.Max())
	require.Equal(t, speed.Value(), speed2.Value())
}
