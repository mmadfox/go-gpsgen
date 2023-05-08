package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSpeed(t *testing.T) {
	type args struct {
		min float64
		max float64
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
				min: minSpeedVal - 1,
				max: 0,
			},
			wantErr: true,
		},
		{
			name: "should be greater by 3 units when the values are equal",
			args: args{
				min: 6,
				max: 6,
			},
			want: want{
				min: 6,
				max: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSpeed(tt.args.min, tt.args.max, 8)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want.min, got.Min())
				require.Equal(t, tt.want.max, got.Max())
			}
		})
	}
}

func TestSpeed_Next(t *testing.T) {
	speed, err := NewSpeed(1, 100, 32)
	require.NoError(t, err)

	var prev float64
	for i := 0; i < 100; i++ {
		speed.Next(float64(i))
		if i < 2 {
			prev = speed.Value()
			continue
		}
		diff := math.Abs(speed.Value() - prev)
		require.NotZero(t, diff)
	}
}
