package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBattery_New(t *testing.T) {
	type args struct {
		min        float64
		max        float64
		chargeTime time.Duration
	}
	type want struct {
		min        float64
		max        float64
		chargeTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "should return error when min value less than 0",
			args: args{
				min: -1,
				max: 1,
			},
			wantErr: true,
		},
		{
			name: "should return error when max value greater than 100",
			args: args{
				min: 0,
				max: 101,
			},
			wantErr: true,
		},
		{
			name: "should return valid object when min greater than max",
			args: args{
				min: 10,
				max: 5,
			},
			want: want{
				min: 5,
				max: 5,
			},
		},
		{
			name: "should return valid object when chargeTime <= 0",
			args: args{
				min:        2,
				max:        4,
				chargeTime: time.Duration(0),
			},
			want: want{
				min:        2,
				max:        4,
				chargeTime: defaultChargeTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBattery(tt.args.min, tt.args.max, tt.args.chargeTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBattery() error = %v, wantErr %v", err, tt.wantErr)
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

func TestBattery_Next(t *testing.T) {
	battery, err := NewBattery(1, 10, time.Minute)
	require.NoError(t, err)
	battery.Reset()
	tick := float64(1)
	want := float64(1)
	for i := 0; i < int(time.Minute.Seconds())+1; i++ {
		battery.Next(tick)
		want = battery.Value()
	}
	require.Equal(t, want, battery.Value())
	require.True(t, battery.IsLow())
	require.Equal(t, time.Minute, battery.ChargeTime())
}

func TestBattery_Snapshot(t *testing.T) {
	battery, err := NewBattery(1, 10, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, battery)
	snap := battery.Snapshot()
	require.NotNil(t, snap)
	require.Equal(t, battery.Min(), snap.Min)
	require.Equal(t, battery.Max(), snap.Max)
	require.Equal(t, 0.0, snap.Val)
	require.Equal(t, battery.ChargeTime(), time.Duration(snap.ChargeTime))
}

func TestBattery_FromSnapshot(t *testing.T) {
	battery, err := NewBattery(1, 10, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, battery)
	snap := battery.Snapshot()
	battery2 := new(Battery)
	battery2.FromSnapshot(snap)
	require.Equal(t, battery.Min(), battery2.Min())
	require.Equal(t, battery.Max(), battery2.Max())
	require.Equal(t, battery.val, battery2.val)
	require.Equal(t, battery.ChargeTime(), battery2.ChargeTime())
}
