package types

import (
	"math"
	"testing"

	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/stretchr/testify/require"
)

func TestSensorToProto(t *testing.T) {
	s1, err := NewSensor("sensor", 1, 30, 8)
	require.NoError(t, err)
	require.NotNil(t, s1)

	protoS1 := s1.ToProto()
	require.NotNil(t, protoS1)
	require.Equal(t, s1.Name(), protoS1.Name)
	require.Equal(t, s1.Min(), protoS1.Min)
	require.Equal(t, s1.Max(), protoS1.Max)
	require.Equal(t, s1.ValueX(), protoS1.ValX)
	require.Equal(t, s1.ValueY(), protoS1.ValY)
	require.NotNil(t, protoS1.Gen)
}

func TestSensorFromProto(t *testing.T) {
	protoS1 := &proto.SensorState{
		Name: "sensor",
		Min:  1,
		Max:  30,
		ValX: 30,
		ValY: 60,
		Gen:  protoGenerator(),
	}

	s1 := new(Sensor)
	s1.FromProto(protoS1)
	require.Equal(t, protoS1.Name, s1.Name())
	require.Equal(t, protoS1.Min, s1.Min())
	require.Equal(t, protoS1.Max, s1.Max())
	require.Equal(t, protoS1.ValX, s1.ValueX())
	require.Equal(t, protoS1.ValY, s1.ValueY())
}

func TestNewSensor(t *testing.T) {
	type args struct {
		name      string
		min       float64
		max       float64
		amplitude int
	}
	type want struct {
		min  float64
		max  float64
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "should return valid sensor, when min and max values are negative",
			args: args{
				name:      "negative sensors",
				min:       -10,
				max:       -5,
				amplitude: 8,
			},
			want: want{
				name: "negative sensors",
				min:  -10,
				max:  -5,
			},
		},
		{
			name: "should return valid sensor, when min greater than max value",
			args: args{
				min:       5,
				max:       1,
				name:      "some",
				amplitude: 8,
			},
			want: want{
				min:  5,
				max:  1,
				name: "some",
			},
		},
		{
			name: "should return error when name is empty",
			args: args{
				amplitude: 8,
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude less than 4",
			args: args{
				amplitude: 0,
				name:      "some",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSensor(tt.args.name, tt.args.min, tt.args.max, tt.args.amplitude)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSensor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.want.min, got.Min())
				require.Equal(t, tt.want.max, got.Max())
				require.Equal(t, tt.want.name, got.Name())
				require.NotZero(t, got.String())
			}
		})
	}
}

func TestSensor_PositiveNext(t *testing.T) {
	sensor, err := NewSensor("positive", 1, 100, 32)
	require.NoError(t, err)

	var prev float64
	for i := 0; i < 100; i++ {
		sensor.Next(float64(i) / 100)
		if i < 2 {
			prev = sensor.ValueY()
			continue
		}
		diff := math.Abs(sensor.ValueY() - prev)
		require.NotZero(t, diff)
	}
}

func TestSensor_NegativeNext(t *testing.T) {
	sensor, err := NewSensor("nagative", -1, -100, 32)
	require.NoError(t, err)
	for i := 0; i < 100; i++ {
		sensor.Next(float64(i) / 100)
		require.NotZero(t, sensor.ValueY())
	}
}
