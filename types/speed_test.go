package types

import (
	"math"
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
)

func TestNewSpeed(t *testing.T) {
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
			name: "should return valid speed value when min, max == 0",
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
			name: "should return error when min less than the minimum value",
			args: args{
				min: MinSpeedVal - 1,
				max: 0,
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude less min value",
			args: args{
				min:       0,
				max:       100,
				amplitude: 0,
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude greater max value",
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

func TestSpeedToProto(t *testing.T) {
	speed, err := NewSpeed(1, 10, 20)
	require.NoError(t, err)
	require.NotNil(t, speed)

	protoSpeed := speed.ToProto()
	require.NotNil(t, protoSpeed)
	require.Equal(t, speed.Min(), protoSpeed.Min)
	require.Equal(t, speed.Max(), protoSpeed.Max)
	require.Equal(t, speed.Value(), protoSpeed.Val)
	require.NotNil(t, protoSpeed.Gen)
}

func TestSpeedFromProto(t *testing.T) {
	protoSpeed := &proto.TypeState{
		Min: 2,
		Max: 4,
		Val: 3,
		Gen: protoGenerator(),
	}

	speed := new(Speed)
	speed.FromProto(protoSpeed)
	require.Equal(t, protoSpeed.Min, speed.Min())
	require.Equal(t, protoSpeed.Max, speed.Max())
	require.Equal(t, protoSpeed.Val, speed.Value())
}

func protoGenerator() *proto.Curve {
	return &proto.Curve{
		Points: []*proto.Curve_ControlPoint{
			{
				Vp: &proto.Curve_Point{},
				Cp: &proto.Curve_Point{},
			},
		},
		Mode: 0,
	}
}
