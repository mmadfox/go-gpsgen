package types

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSensor_New(t *testing.T) {
	type args struct {
		name      string
		min       float64
		max       float64
		amplitude int
		mode      SensorMode
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
				mode:      WithSensorRandomMode,
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
			got, err := NewSensor(tt.args.name, tt.args.min, tt.args.max, tt.args.amplitude, 0)
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
	sensor, err := NewSensor("positive", 1, 100, 32, 0)
	require.NoError(t, err)

	sensor.Shuffle()

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
	sensor, err := NewSensor("nagative", -1, -100, 32, 0)
	require.NoError(t, err)

	sensor.Shuffle()

	for i := 0; i < 100; i++ {
		sensor.Next(float64(i) / 100)
		require.NotZero(t, sensor.ValueY())
	}
}

func TestSensor_Snapshot(t *testing.T) {
	sensor, err := NewSensor("pressure", 0.1, 1.0, 16, 0)
	require.NoError(t, err)
	require.NotNil(t, sensor)
	snap := sensor.Snapshot()
	require.NotNil(t, snap)
	require.Equal(t, sensor.ID(), snap.Id)
	require.Equal(t, sensor.Min(), snap.Min)
	require.Equal(t, sensor.Max(), snap.Max)
	require.Equal(t, sensor.ValueX(), snap.ValX)
	require.Equal(t, sensor.ValueY(), snap.ValY)
	require.Equal(t, sensor.Name(), snap.Name)
	require.NotNil(t, snap.Gen)
}

func TestSensor_FromSnapshot(t *testing.T) {
	sensor, err := NewSensor("pressure", 0.1, 1.0, 16, 0)
	require.NoError(t, err)
	require.NotNil(t, sensor)
	snap := sensor.Snapshot()
	sensor2 := new(Sensor)
	sensor2.FromSnapshot(snap)
	require.Equal(t, sensor.ID(), sensor2.ID())
	require.Equal(t, sensor.Min(), sensor2.Min())
	require.Equal(t, sensor.Max(), sensor2.Max())
	require.Equal(t, sensor.ValueX(), sensor2.ValueX())
	require.Equal(t, sensor.ValueY(), sensor2.ValueY())
	require.Equal(t, sensor.Name(), sensor2.Name())
}

func TestRestoreSensor(t *testing.T) {
	type args struct {
		id        string
		name      string
		min       float64
		max       float64
		amplitude int
		mode      SensorMode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error when id is empty",
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when id is invalid",
			args: args{
				id: "badid",
			},
			wantErr: true,
		},
		{
			name: "should return error when name is empty",
			args: args{
				id:   uuid.NewString(),
				name: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude < 4",
			args: args{
				id:        uuid.NewString(),
				name:      "s1",
				amplitude: 3,
			},
			wantErr: true,
		},
		{
			name: "should not return error when invalid sensor mode",
			args: args{
				id:        uuid.NewString(),
				name:      "s1",
				amplitude: 4,
				min:       1,
				max:       2,
				mode:      SensorMode(999),
			},
			wantErr: false,
		},
		{
			name: "should not return error when all params are valid",
			args: args{
				id:        uuid.NewString(),
				name:      "s1",
				min:       0,
				max:       100,
				amplitude: 5,
				mode:      WithSensorRandomMode | WithSensorStartMode | WithSensorEndMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RestoreSensor(tt.args.id, tt.args.name, tt.args.min, tt.args.max, tt.args.amplitude, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreSensor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			require.NotNil(t, got)
		})
	}
}
