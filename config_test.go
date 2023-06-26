package gpsgen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/mmadfox/go-gpsgen/types"
	"github.com/stretchr/testify/require"
)

func TestConfig_DeviceIDFromConfig(t *testing.T) {
	conf := getFullValidConfig()
	proc, err := conf.NewDevice()
	require.NoError(t, err)
	require.NotNil(t, proc)
	require.Equal(t, conf.ID, proc.ID().String())
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		conf        func() *Config
		containsErr string
	}{
		{
			name: "should return error when the device model value is too short",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Model = ""
				return conf
			},
			containsErr: fmt.Sprintf("Model must be > %d", constraints.Model.Min),
		},
		{
			name: "should return error when device model value is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Model = strings.Repeat("a", types.MaxModelLen+1)
				return conf
			},
			containsErr: fmt.Sprintf("Model must be < %d", constraints.Model.Max),
		},
		{
			name: "should return error when properties list is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Properties = make(map[string]string, 0)
				for i := 0; i < constraints.Properties.Max+1; i++ {
					conf.Properties[fmt.Sprintf("key-%d", i)] = "one"
				}
				return conf
			},
			containsErr: fmt.Sprintf("Properties must be < %d", constraints.Properties.Max),
		},
		{
			name: "should return error when properties key is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				key := strings.Repeat("a", constraints.Properties.MaxKeyLen+1)
				conf.Properties[key] = "one"
				return conf
			},
			containsErr: fmt.Sprintf("Properties.key{aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa} must be < %d", constraints.Properties.MaxKeyLen),
		},
		{
			name: "should return error when properties value is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				value := strings.Repeat("a", constraints.Properties.MaxValueLen+1)
				conf.Properties["one"] = value
				return conf
			},
			containsErr: fmt.Sprintf("Properties.value{aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa} must be < %d", constraints.Properties.MaxValueLen),
		},
		{
			name: "should return error when description value is too short",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Description = "1"
				return conf
			},
			containsErr: fmt.Sprintf("Description must be > %d", constraints.Description.Min),
		},
		{
			name: "should return error when description value is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Description = strings.Repeat("a", constraints.Description.Max+1)
				return conf
			},
			containsErr: fmt.Sprintf("Description must be < %d", constraints.Description.Max),
		},
		{
			name: "should return error when invalid min speed",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Speed.Min = -1
				return conf
			},
			containsErr: fmt.Sprintf("Speed.Min must be > %f", constraints.Speed.Min),
		},
		{
			name: "should return error when invalid max speed",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Speed.Max = constraints.Speed.Max + 1
				return conf
			},
			containsErr: fmt.Sprintf("Speed.Max must be < %f", constraints.Speed.Max),
		},
		{
			name: "should return error when invalid min speed amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Speed.Amplitude = Amplitude4 - 1
				return conf
			},
			containsErr: fmt.Sprintf("Speed.Amplitude must be > %d", Amplitude4),
		},
		{
			name: "should return error when invalid max speed amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Speed.Amplitude = Amplitude512 + 1
				return conf
			},
			containsErr: fmt.Sprintf("Speed.Amplitude must be < %d", Amplitude512),
		},
		{
			name: "should return error when invalid min battery",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Battery.Min = -1
				return conf
			},
			containsErr: fmt.Sprintf("Battery.Min must be > %f", constraints.Battery.Min),
		},
		{
			name: "should return error when invalid max battery",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Battery.Max = constraints.Battery.Max + 1
				return conf
			},
			containsErr: fmt.Sprintf("Battery.Max must be < %f", constraints.Battery.Max),
		},
		{
			name: "should return error when invalid min elevation",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Elevation.Min = -1
				return conf
			},
			containsErr: fmt.Sprintf("Elevation.Min must be > %f", constraints.Elevation.Min),
		},
		{
			name: "should return error when invalid max elevation",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Elevation.Max = constraints.Elevation.Max + 1
				return conf
			},
			containsErr: fmt.Sprintf("Elevation.Max must be < %f", constraints.Elevation.Max),
		},
		{
			name: "should return error when invalid min elevation amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Elevation.Amplitude = Amplitude4 - 1
				return conf
			},
			containsErr: fmt.Sprintf("Elevation.Amplitude must be > %d", Amplitude4),
		},
		{
			name: "should return error when invalid max elevation amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Elevation.Amplitude = Amplitude512 + 1
				return conf
			},
			containsErr: fmt.Sprintf("Elevation.Amplitude must be < %d", Amplitude512),
		},
		{
			name: "should return error when sensors list is too long",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Sensors = make([]Sensor, constraints.Sensors.Max+1)
				for i := 0; i < constraints.Sensors.Max; i++ {
					conf.Sensors[i] = Sensor{Name: "ok"}
				}
				return conf
			},
			containsErr: fmt.Sprintf("Sensors.Max must be < %d", constraints.Sensors.Max),
		},
		{
			name: "should return error when invalid min sensor amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Sensors = []Sensor{
					{Name: "s1", Amplitude: Amplitude4 - 1},
				}
				return conf
			},
			containsErr: fmt.Sprintf("Sensors[0].Amplitude must be > %d", Amplitude4),
		},
		{
			name: "should return error when invalid max sensor amplitude",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Sensors = []Sensor{
					{Name: "s1", Amplitude: Amplitude512 + 1},
				}
				return conf
			},
			containsErr: fmt.Sprintf("Sensors[0].Amplitude must be < %d", Amplitude512),
		},
		{
			name: "should return error when invalid max offline",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Offline.Max = constraints.Offline.Max + 1
				return conf
			},
			containsErr: fmt.Sprintf("Offline.Max must be < %d", constraints.Offline.Max),
		},
		{
			name: "should return error when invalid min routes",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Routes = []*navigator.Route{}
				return conf
			},
			containsErr: fmt.Sprintf("Routes.Min must be > %d", constraints.Routes.Min),
		},
		{
			name: "should return error when invalid max routes",
			conf: func() *Config {
				conf := getFullValidConfig()
				conf.Routes = []*navigator.Route{}
				for i := 0; i < constraints.Routes.Max+1; i++ {
					r, _ := route.Generate()
					conf.Routes = append(conf.Routes, r)
				}
				return conf
			},
			containsErr: fmt.Sprintf("Routes.Max must be < %d", constraints.Routes.Max),
		},
		{
			name: "should return error when invalid max distance routes",
			conf: func() *Config {
				conf := getFullValidConfig()
				r, err := navigator.NewRoute([][]navigator.Point{
					{
						{X: -31.842791814115763, Y: 23.040317669562427},
						{X: 70.12002910024887, Y: 155.79766533938636},
					},
				})
				require.NoError(t, err)
				conf.Routes = []*navigator.Route{r}
				return conf
			},
			containsErr: fmt.Sprintf("Routes.MaxDistance must be < %f", constraints.Routes.MaxDistance),
		},
		{
			name: "should return error when invalid min distance routes",
			conf: func() *Config {
				conf := getFullValidConfig()
				r, err := navigator.NewRoute([][]navigator.Point{
					{
						{X: -0.7397828385093419, Y: 28.803987823772303},
						{X: -0.7398090096154561, Y: 28.803991612584213},
					},
				})
				require.NoError(t, err)
				conf.Routes = []*navigator.Route{r}
				return conf
			},
			containsErr: fmt.Sprintf("Routes.MinDistance must be > %f", constraints.Routes.MinDistance),
		},
		{
			name: "should return valid device",
			conf: getFullValidConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := tt.conf()
			if err := conf.Validate(); err != nil {
				require.Contains(t, err.Error(), tt.containsErr)
			}
		})
	}
}

func getFullValidConfig() *Config {
	conf := NewConfig()
	conf.ID = uuid.NewString()
	conf.Model = "Model"
	conf.UserID = "UserID"
	conf.Properties = map[string]string{"foo": "bar", "xim": "mix"}
	conf.Description = "Description"
	conf.Speed.Min = 2
	conf.Speed.Max = 5
	conf.Speed.Amplitude = Amplitude8
	conf.Battery.Min = 4
	conf.Battery.Max = 40
	conf.Elevation.Min = 1
	conf.Elevation.Max = 10
	conf.Elevation.Amplitude = Amplitude32
	conf.Offline.Min = 0
	conf.Offline.Max = 10
	conf.Sensors = []Sensor{
		{Name: "sensor", Min: 1, Max: 2, Amplitude: Amplitude16},
	}
	conf.Routes, _ = route.RoutesForChina()
	return conf
}
