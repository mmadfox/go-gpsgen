package types

import (
	"math"
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
)

func TestNewBattery(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	type want struct {
		min float64
		max float64
		val float64
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "should return error, when min value less than 0",
			args: args{
				min: -1,
				max: 1,
			},
			wantErr: true,
		},
		{
			name: "should return error, when max value greater than 100",
			args: args{
				min: 0,
				max: 101,
			},
			wantErr: true,
		},
		{
			name: "should return valid object, when min greater than max",
			args: args{
				min: 10,
				max: 5,
			},
			want: want{
				min: 5,
				max: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBattery(tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBattery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want.min, got.Min())
				require.Equal(t, tt.want.max, got.Max())
				require.Equal(t, tt.want.val, got.Value())
				require.NotZero(t, got.String())
			}
		})
	}
}

func TestBatteryToProto(t *testing.T) {
	battery, err := NewBattery(1, 10)
	require.NoError(t, err)
	require.NotNil(t, battery)

	protobattery := battery.ToProto()
	require.NotNil(t, protobattery)
	require.Equal(t, battery.Min(), protobattery.Min)
	require.Equal(t, battery.Max(), protobattery.Max)
	require.Equal(t, battery.Value(), protobattery.Val)
	require.NotNil(t, protobattery.Gen)
}

func TestBatteryFromProto(t *testing.T) {
	protobattery := &proto.TypeState{
		Min: 2,
		Max: 4,
		Val: 3,
		Gen: protoGenerator(),
	}

	battery := new(Battery)
	battery.FromProto(protobattery)
	require.Equal(t, protobattery.Min, battery.Min())
	require.Equal(t, protobattery.Max, battery.Max())
	require.Equal(t, protobattery.Val, battery.Value())
}

func TestBattery_Next(t *testing.T) {
	battery, err := NewBattery(1, 100)
	require.NoError(t, err)

	var prev float64
	for i := 0; i < 100; i++ {
		battery.Next(float64(i) / 100)
		if i < 2 {
			prev = battery.Value()
			continue
		}
		diff := math.Abs(battery.Value() - prev)
		require.NotZero(t, diff)
	}
}
