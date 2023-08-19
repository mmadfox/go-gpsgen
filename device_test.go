package gpsgen

import (
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/types"
	"github.com/stretchr/testify/require"
)

var (
	// 1segment: 300m
	track300m1segment, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.49331396675268, Lat: 29.5299004724652},
		{Lon: 106.49523863664103, Lat: 29.532016484207674},
	})
	// 1segment: 301m, 2segment: 697m
	track1km2segment, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.46029732125612, Lat: 29.532477955504234},
		{Lon: 106.46229029469634, Lat: 29.534568282538686},
		{Lon: 106.4659759305107, Lat: 29.53996017665672},
	})
	track3km7segments, _ = navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 106.48818494999364, Lat: 29.526967168711465},
		{Lon: 106.48818494999364, Lat: 29.53306336155346},
		{Lon: 106.49038855444923, Lat: 29.53586552047419},
		{Lon: 106.49055806248373, Lat: 29.537291150475014},
		{Lon: 106.48007681565423, Lat: 29.53522643843239},
		{Lon: 106.47905976744477, Lat: 29.534439870375238},
		{Lon: 106.47617813084906, Lat: 29.53055610091529},
		{Lon: 106.4755283500491, Lat: 29.529843240954946},
	})
)

func makeSensor(t *testing.T, name string) *types.Sensor {
	sensor, err := types.NewSensor(name, 1, 2, 4, 0)
	require.NoError(t, err)
	return sensor
}

func TestDevice_New(t *testing.T) {
	tests := []struct {
		name    string
		args    func() *DeviceOptions
		assert  func(*Device)
		wantErr bool
	}{
		{
			name: "should return error when elevation min is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Navigator.Elevation.Min = navigator.MinElevation - 1
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when elevation max is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Navigator.Elevation.Max = navigator.MaxElevation + 1
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when speed min is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Speed.Min = types.MinSpeedVal - 1
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when speed max is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Speed.Min = types.MaxSpeedVal + 1
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when battery min is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Battery.Min = -1
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when battery max is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Battery.Max = 101
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when color is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Color = "#f"
				return opts
			},
			wantErr: true,
		},
		{
			name: "should return error when model is invalid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Model = "A"
				return opts
			},
			wantErr: true,
		},
		{
			name: "should not return error when model is empty",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Model = ""
				return opts
			},
			assert: func(d *Device) {
				require.NotEmpty(t, d.Model())
			},
			wantErr: false,
		},
		{
			name: "should not return error when model is valid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.Model = "ModelNXt87a"
				return opts
			},
			assert: func(d *Device) {
				require.NotEmpty(t, d.Model())
			},
			wantErr: false,
		},
		{
			name: "should not return error when id is empty",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.ID = ""
				return opts
			},
			assert: func(d *Device) {
				require.NotEmpty(t, d.ID())
			},
			wantErr: false,
		},
		{
			name:    "should return new device when opts is nil",
			args:    func() *DeviceOptions { return nil },
			wantErr: false,
		},
		{
			name: "should not return error when all params are valid",
			args: func() *DeviceOptions {
				opts := NewDeviceOptions()
				opts.ID = "id"
				opts.Descr = "descr"
				opts.Model = "model"
				opts.UserID = "userID"
				opts.Color = "#ffffff"
				return opts
			},
			assert: func(d *Device) {
				require.Equal(t, "id", d.ID())
				require.Equal(t, "descr", d.Descr())
				require.Equal(t, "model", d.Model())
				require.Equal(t, "userID", d.UserID())
				require.Equal(t, "#ffffff", d.Color())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev, err := NewDevice(tt.args())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.assert != nil {
				tt.assert(dev)
			}
		})
	}
}

func TestDevice_Next(t *testing.T) {
	type args struct {
		opts func() *DeviceOptions
	}

	tests := []struct {
		name    string
		args    args
		arrange func(*Device)
		assert  func(d *Device)
		want    bool
		wantErr bool
	}{
		{
			name: "should return false and skip next loop when mutex locked",
			arrange: func(d *Device) {
				d.mu.Lock()
			},
			assert: func(d *Device) {
				tick := atomic.LoadUint32(&d.tick)
				require.Equal(t, uint32(1), tick)
			},
			want: false,
		},
		{
			name: "should return false when no routes",
			want: false,
		},
		{
			name: "should return false when navigator is offline",
			arrange: func(d *Device) {
				route := navigator.NewRoute()
				route.AddTrack(track300m1segment)
				d.AddRoute(route)
				d.ToOffline()
			},
			assert: func(d *Device) {
				require.True(t, d.IsOffline())
				state := d.State()
				require.True(t, state.IsOffline)
				require.NotZero(t, state.OfflineDuration)
			},
			want: false,
		},
		{
			name: "should return false when battery is low",
			args: args{
				opts: func() *DeviceOptions {
					opts := NewDeviceOptions()
					opts.Battery.ChargeTime = 0
					return opts
				},
			},
			arrange: func(d *Device) {
				route := navigator.NewRoute()
				route.AddTrack(track300m1segment)
				d.AddRoute(route)
			},
			assert: func(d *Device) {
				require.True(t, d.IsOffline())
				require.True(t, d.State().IsOffline)
				require.NotZero(t, d.State().OfflineDuration)
			},
		},
		{
			name: "should return false when route is finish",
			arrange: func(d *Device) {
				route := navigator.NewRoute()
				route.AddTrack(track300m1segment)
				d.AddRoute(route)
				d.DestinationTo(300)
				d.Update()
			},
			assert: func(d *Device) {
				require.True(t, d.IsOffline())
				require.True(t, d.State().IsOffline)
				require.NotZero(t, d.State().OfflineDuration)
				require.Zero(t, d.State().Distance.CurrentDistance)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts *DeviceOptions
			if tt.args.opts != nil {
				opts = tt.args.opts()
			}
			dev, err := NewDevice(opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.arrange != nil {
				tt.arrange(dev)
			}
			require.Equal(t, tt.want, dev.Next(1))
			if tt.assert != nil {
				tt.assert(dev)
			}
		})
	}
}

func TestDevice_SetModel(t *testing.T) {
	type args struct {
		model string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "should return error when model is invalid",
			args: args{
				model: "A",
			},
			wantErr: true,
		},
		{
			name: "should return error when model is empty",
			args: args{
				model: "",
			},
			wantErr: true,
		},
		{
			name: "should not return error when model is valid",
			args: args{
				model: "Model",
			},
			want: "Model",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			if err := d.SetModel(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Device.SetModel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, d.Model())
		})
	}
}

func TestDevice_SetDescription(t *testing.T) {
	type args struct {
		descr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should assign valid description",
			args: args{
				descr: "Desct",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			d.SetDescription(tt.args.descr)
			require.Equal(t, tt.args.descr, d.Descr())
		})
	}
}

func TestDevice_SetColor(t *testing.T) {
	type args struct {
		color colorful.Color
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return error when color is invalid",
			args: args{
				color: colorful.Color{R: -1, G: -1},
			},
			wantErr: true,
		},
		{
			name: "should not return error when color is valid",
			args: args{
				color: colorful.Color{R: 0, G: 0, B: 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			if err := d.SetColor(tt.args.color); (err != nil) != tt.wantErr {
				t.Errorf("Device.SetColor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.args.color.Hex(), d.Color())
		})
	}
}

func TestDevice_AddSensor(t *testing.T) {
	type args struct {
		sensor func() []*types.Sensor
	}

	tests := []struct {
		name   string
		args   args
		assert func(d *Device, sensors []*types.Sensor)
	}{
		{
			name: "should not add when list is empty",
			args: args{
				sensor: func() []*types.Sensor {
					return []*types.Sensor{}
				},
			},
			assert: func(d *Device, _ []*types.Sensor) {
				require.Zero(t, d.NumSensors())
				require.Zero(t, d.State().Sensors)
			},
		},
		{
			name: "should not add when sensor already exists",
			args: args{
				sensor: func() []*types.Sensor {
					s1 := makeSensor(t, "s1")
					return []*types.Sensor{s1, s1}
				},
			},
			assert: func(d *Device, _ []*types.Sensor) {
				require.Equal(t, 1, d.NumSensors())
				require.Len(t, d.State().Sensors, 1)
			},
		},
		{
			name: "should not add when sensors is nil",
			args: args{
				sensor: func() []*types.Sensor {
					return []*types.Sensor{nil, nil, nil}
				},
			},
			assert: func(d *Device, _ []*types.Sensor) {
				require.Zero(t, d.NumSensors())
				require.Len(t, d.State().Sensors, 0)
			},
		},
		{
			name: "should add when all sensors are valid",
			args: args{
				sensor: func() []*types.Sensor {
					s1 := makeSensor(t, "s1")
					s2 := makeSensor(t, "s2")
					return []*types.Sensor{s1, s2}
				},
			},
			assert: func(d *Device, _ []*types.Sensor) {
				require.Equal(t, 2, d.NumSensors())
				require.Len(t, d.State().Sensors, 2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)

			sensors := tt.args.sensor()
			d.AddSensor(sensors...)

			if tt.assert != nil {
				tt.assert(d, sensors)
			}
		})
	}
}

func TestDevice_ResetSensors(t *testing.T) {
	d := NewTracker()
	d.AddSensor(makeSensor(t, "s1"))
	d.AddSensor(makeSensor(t, "s2"))
	require.Equal(t, 2, d.NumSensors())
	require.Len(t, d.State().Sensors, 2)

	d.ResetSensors()
	require.Equal(t, 0, d.NumSensors())
	require.Len(t, d.State().Sensors, 0)
	require.Len(t, d.Sensors(), 0)
}

func TestDevice_RemoveSensor(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(d *Device) string
		assert  func(d *Device, id string)
	}{
		{
			name: "should return false when sensorID is empty",
			arrange: func(d *Device) string {
				return ""
			},
			assert: func(d *Device, sid string) {
				require.False(t, d.RemoveSensor(sid))
			},
		},
		{
			name: "should return false when sensor not found",
			arrange: func(d *Device) string {
				return "someid"
			},
			assert: func(d *Device, sid string) {
				require.False(t, d.RemoveSensor(sid))
				require.Empty(t, d.State().Sensors)
			},
		},
		{
			name: "should return true when sensor removed",
			arrange: func(d *Device) string {
				sensor, err := types.NewSensor("s1", 1, 2, 4, 0)
				require.NoError(t, err)
				d.AddSensor(sensor)
				return sensor.ID()
			},
			assert: func(d *Device, sid string) {
				require.Len(t, d.State().Sensors, 1)
				require.Equal(t, 1, d.NumSensors())
				require.True(t, d.RemoveSensor(sid))
				require.Empty(t, d.State().Sensors)
				require.Equal(t, 0, d.NumSensors())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			id := tt.arrange(d)
			tt.assert(d, id)
		})
	}
}

func TestDevice_GetSensor(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(d *Device) string
		assert  func(d *Device, sid string)
	}{
		{
			name: "should return sensor by index",
			arrange: func(d *Device) string {
				sensor, err := types.NewSensor("s1", 1, 2, 8, 0)
				require.NoError(t, err)
				d.AddSensor(sensor)
				return sensor.ID()
			},
			assert: func(d *Device, sid string) {
				s1 := d.SensorAt(0)
				require.NotNil(t, s1)
				require.Equal(t, sid, s1.ID())
			},
		},
		{
			name: "should not return sensor by index when list of sensors is empty",
			arrange: func(d *Device) string {
				return ""
			},
			assert: func(d *Device, sid string) {
				s1 := d.SensorAt(0)
				require.Nil(t, s1)
			},
		},
		{
			name: "should return false when sensor not found by id",
			arrange: func(d *Device) string {
				return ""
			},
			assert: func(d *Device, sid string) {
				s1, ok := d.SensorByID(sid)
				require.Nil(t, s1)
				require.False(t, ok)
			},
		},
		{
			name: "should return sensor by id",
			arrange: func(d *Device) string {
				s1 := makeSensor(t, "s1")
				d.AddSensor(s1)
				return s1.ID()
			},
			assert: func(d *Device, sid string) {
				s1, ok := d.SensorByID(sid)
				require.NotNil(t, s1)
				require.True(t, ok)
				require.Equal(t, sid, s1.ID())
			},
		},
		{
			name: "should iter each sensor when list of sensors is not empty",
			arrange: func(d *Device) string {
				sensor, err := types.NewSensor("s1", 1, 2, 8, 0)
				require.NoError(t, err)
				d.AddSensor(sensor)
				return sensor.ID()
			},
			assert: func(d *Device, sid string) {
				c := 0
				d.EachSensor(func(i int, s *types.Sensor) bool {
					c++
					require.NotNil(t, s)
					return true
				})
				require.Equal(t, 1, c)

				c = 0
				d.EachSensor(func(i int, s *types.Sensor) bool {
					c++
					require.NotNil(t, s)
					return false
				})
				require.Equal(t, 1, c)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			id := tt.arrange(d)
			tt.assert(d, id)
		})
	}
}

func TestDevice_AddRoute(t *testing.T) {
	type args struct {
		routes []*navigator.Route
	}
	tests := []struct {
		name    string
		args    args
		assert  func(*Device)
		wantErr bool
	}{
		{
			name: "should return error when routes are empty",
			args: args{
				routes: []*navigator.Route{},
			},
			wantErr: true,
		},
		{
			name: "should not return error when routes are valid",
			args: args{
				routes: testRoutes(),
			},
			assert: func(d *Device) {
				require.NotZero(t, d.NumRoutes())
				require.NotZero(t, d.State().Routes)
			},
			wantErr: false,
		},
		{
			name: "should not return error when track is added to route",
			args: args{
				routes: testRoutes(),
			},
			assert: func(d *Device) {
				route0 := d.RouteAt(0)
				require.NotNil(t, route0)
				actualNumTracks := route0.NumTracks()

				route0.AddTrack(track3km7segments)
				require.True(t, d.Next(1))
				actualVersion := d.Version()
				require.Equal(t, actualNumTracks+1, route0.NumTracks())
				require.Equal(t, actualNumTracks+1, len(d.State().Routes.Routes[0].Tracks))

				require.True(t, d.Next(1))
				require.Equal(t, actualVersion, d.Version())
				route0.RemoveTrack(route0.TrackAt(0).ID())
				require.True(t, d.Next(1))
				require.NotEqual(t, actualVersion, d.Version())
				require.Equal(t, actualNumTracks, route0.NumTracks())
				require.Equal(t, actualNumTracks, len(d.State().Routes.Routes[0].Tracks))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			if err := d.AddRoute(tt.args.routes...); (err != nil) != tt.wantErr {
				t.Errorf("Device.AddRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tt.assert != nil {
				tt.assert(d)
			}
		})
	}
}

func TestDevice_RemoveRoute(t *testing.T) {
	type args struct {
		routeID string
	}
	tests := []struct {
		name   string
		args   args
		assert func(*Device)
		want   bool
	}{
		{
			name: "should return false when route not found",
			args: args{
				routeID: "someid",
			},
			want: false,
		},
		{
			name: "should return true when route found",
			args: args{
				routeID: "3caaac65-b01a-4865-8d01-6aba1fbeab69",
			},
			assert: func(d *Device) {
				require.Zero(t, d.NumRoutes())
				require.Empty(t, d.State().Routes.Routes)
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			route := navigator.RestoreRoute("3caaac65-b01a-4865-8d01-6aba1fbeab69", "#ffffff", nil)
			route.AddTrack(track1km2segment).AddTrack(track300m1segment)
			d.AddRoute(route)
			if got := d.RemoveRoute(tt.args.routeID); got != tt.want {
				t.Errorf("Device.RemoveRoute() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(d)
			}
		})
	}
}

func TestDevice_RemoveTrack(t *testing.T) {
	type args struct {
		routeID string
		trackID string
	}
	tests := []struct {
		name   string
		args   args
		assert func(*Device)
		want   bool
	}{
		{
			name: "should return false when trackID and routeID is empty",
			args: args{
				routeID: "",
				trackID: "",
			},
			assert: func(d *Device) {
				require.Equal(t, 1, d.NumRoutes())
				require.Equal(t, 1, d.RouteAt(0).NumTracks())
			},
			want: false,
		},
		{
			name: "should return false when track not found",
			args: args{
				routeID: "20994de4-8c4b-4bdf-b39e-b2affd42b0dd",
				trackID: "someid",
			},
			assert: func(d *Device) {
				require.Equal(t, 1, d.NumRoutes())
				require.Equal(t, 1, d.RouteAt(0).NumTracks())
			},
			want: false,
		},
		{
			name: "should return false when track not found",
			args: args{
				routeID: "20994de4-8c4b-4bdf-b39e-b2affd42b0dd",
				trackID: "8d102c1d-f6de-4afc-9585-33083486ef12",
			},
			assert: func(d *Device) {
				require.Equal(t, 1, d.NumRoutes())
				require.Equal(t, 0, d.RouteAt(0).NumTracks())
				require.Len(t, d.State().Routes.Routes, 1)
				require.Len(t, d.State().Routes.Routes[0].Tracks, 0)
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDevice(nil)
			require.NoError(t, err)
			track, err := navigator.RestoreTrack("8d102c1d-f6de-4afc-9585-33083486ef12", "#000000", []geo.LatLonPoint{
				{Lon: 106.49331396675268, Lat: 29.5299004724652},
				{Lon: 106.49523863664103, Lat: 29.532016484207674},
			})
			require.NoError(t, err)
			route := navigator.RestoreRoute("20994de4-8c4b-4bdf-b39e-b2affd42b0dd", "#ff0000", nil)
			route.AddTrack(track)
			d.AddRoute(route)
			if got := d.RemoveTrack(tt.args.routeID, tt.args.trackID); got != tt.want {
				t.Errorf("Device.RemoveTrack() = %v, want %v", got, tt.want)
			}
			if tt.assert != nil {
				tt.assert(d)
			}
		})
	}
}

func TestDevice_Update(t *testing.T) {
	tests := []struct {
		name    string
		arrange func(*Device)
		assert  func(*Device)
	}{
		{
			name: "should be device state update",
			arrange: func(d *Device) {
				s1, _ := types.NewSensor("s1", 1, 2, 8, 0)
				d.AddSensor(s1)
				s2, _ := types.NewSensor("s2", 0.1, 1, 16, 0)
				d.AddSensor(s2)
				d.AddRoute(testRoutes()...)
				d.Next(3)
			},
			assert: func(d *Device) {
				require.NotZero(t, d.CurrentBearing())
				require.NotZero(t, d.Distance())
				require.NotZero(t, d.CurrentDistance())
				require.NotZero(t, d.RouteDistance())
				require.NotZero(t, d.CurrentRouteDistance())
				require.NotZero(t, d.TrackDistance())
				require.NotZero(t, d.CurrentTrackDistance())
				require.NotZero(t, d.SegmentDistance())
				require.NotZero(t, d.CurrentSegmentDistance())
				require.False(t, d.IsFinish())
				require.False(t, d.IsOffline())
				require.NotNil(t, d.CurrentRoute())
				require.NotNil(t, d.CurrentTrack())
				require.False(t, d.CurrentSegment().IsEmpty())
				require.Equal(t, 0, d.RouteIndex())
				require.Equal(t, 0, d.TrackIndex())
				require.Equal(t, 0, d.SegmentIndex())
				require.NotZero(t, d.Location())

				state := d.State()
				require.NotEmpty(t, state.Id)
				require.NotEmpty(t, state.UserId)
				require.NotEmpty(t, state.Model)
				require.NotEmpty(t, state.Description)
				require.NotZero(t, state.Duration)
				require.NotZero(t, state.Speed)
				require.False(t, state.IsOffline)
				require.Zero(t, state.OfflineDuration)
				require.NotEmpty(t, state.Color)
				require.NotZero(t, state.TimeEstimate)
				require.NotZero(t, state.Tick)
				require.NotZero(t, state.Battery.Charge)
				require.NotZero(t, state.Battery.ChargeTime)
				require.Len(t, state.Sensors, 2)
				for i := 0; i < 2; i++ {
					require.NotEmpty(t, state.Sensors[i].Id)
					require.NotZero(t, state.Sensors[i].ValY)
					require.NotZero(t, state.Sensors[i].ValX)
					require.NotEmpty(t, state.Sensors[i].Name)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewDeviceOptions()
			opts.UserID = "userID"
			opts.Descr = "descr"
			opts.Color = "#ffffff"
			d, err := NewDevice(opts)
			require.NoError(t, err)
			if tt.arrange != nil {
				tt.arrange(d)
			}
			d.Update()
			if tt.assert != nil {
				tt.assert(d)
			}
		})
	}
}

func TestDevice_EachRoute(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	expectedRoutes := testRoutes()
	d.AddRoute(expectedRoutes...)
	actual := 0
	d.EachRoute(func(i int, r *navigator.Route) bool {
		require.NotNil(t, r)
		actual++
		return true
	})
	require.Equal(t, len(expectedRoutes), actual)
}

func TestDevice_ResetRoutes(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	d.AddRoute(testRoutes()...)
	d.ResetRoutes()
	require.Equal(t, 0, d.NumRoutes())
	require.Len(t, d.State().Routes.Routes, 0)
}

func TestDevice_DestinationTo(t *testing.T) {
	meters := 100.0
	d, err := NewDevice(nil)
	require.NoError(t, err)
	d.AddRoute(testRoutes()...)
	require.True(t, d.DestinationTo(meters))
	require.Equal(t, meters, d.CurrentDistance())
	prevDist := d.CurrentDistance()
	d.Next(1)
	require.GreaterOrEqual(t, d.CurrentDistance(), prevDist)
}

func TestDevice_MoveToSegment(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	d.AddRoute(testRoutes()...)
	require.True(t, d.MoveToSegment(0, 0, 1))
	require.Equal(t, 0, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 1, d.SegmentIndex())
}

func TestDevice_MoveToTrack(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	require.True(t, d.MoveToTrack(0, 1))
	require.Equal(t, 0, d.RouteIndex())
	require.Equal(t, 1, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_MoveToRoute(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)
	require.True(t, d.MoveToRoute(1))
	require.Equal(t, 1, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_MoveToRouteByID(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)
	require.True(t, d.MoveToRouteByID(route1.ID()))
	require.Equal(t, 1, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_MoveToTrackByID(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)
	require.True(t, d.MoveToTrackByID(route1.ID(), track300m1segment.ID()))
	require.Equal(t, 1, d.RouteIndex())
	require.Equal(t, 1, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_ToPrevRoute(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)

	require.True(t, d.MoveToRoute(1))
	require.Equal(t, 1, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())

	require.True(t, d.ToPrevRoute())

	require.Equal(t, 0, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_ToNextRoute(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)

	require.Equal(t, 0, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())

	require.True(t, d.ToNextRoute())

	require.Equal(t, 1, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
}

func TestDevice_ResetNavigator(t *testing.T) {
	d, err := NewDevice(nil)
	require.NoError(t, err)
	route := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route)
	route1 := navigator.RouteFromTracks(track1km2segment, track300m1segment)
	d.AddRoute(route1)

	require.True(t, d.MoveToRoute(1))
	d.ResetNavigator()
	require.Equal(t, 0, d.RouteIndex())
	require.Equal(t, 0, d.TrackIndex())
	require.Equal(t, 0, d.SegmentIndex())
	require.Equal(t, 2, d.NumRoutes())
	require.Zero(t, d.CurrentDistance())
}

func TestDevice_Snapshot(t *testing.T) {
	opts := NewDeviceOptions()
	opts.Battery.Min = 1
	opts.Model = "model"
	opts.UserID = "userId"
	opts.Color = "#ffffff"
	opts.Descr = "some descr"
	d, err := NewDevice(opts)
	require.NoError(t, err)
	s1, _ := types.NewSensor("s1", 1, 10, 8, 0)
	d.AddSensor(s1)
	s2, _ := types.NewSensor("s2", 15, 20, 4, 0)
	d.AddSensor(s2)
	d.AddRoute(testRoutes()...)
	for i := 0; i < 10; i++ {
		require.True(t, d.Next(1.0))
	}

	snap := d.Snapshot()
	require.NotNil(t, snap)

	require.NotEmpty(t, snap.Id)
	require.NotEmpty(t, snap.UserId)
	require.NotEmpty(t, snap.Tick)
	require.NotEmpty(t, snap.Descr)
	require.NotEmpty(t, snap.Status)
	require.NotZero(t, snap.Duration)

	require.Len(t, snap.Sensors, 2)
	for i := 0; i < len(snap.Sensors); i++ {
		s := snap.Sensors[i]
		require.NotEmpty(t, s.Id)
		require.NotEmpty(t, s.Name)
		require.NotZero(t, s.Min)
		require.NotZero(t, s.Max)
		require.NotZero(t, s.Gen)
	}

	require.NotNil(t, snap.Battery)
	require.NotZero(t, snap.Battery.Val)
	require.NotZero(t, snap.Battery.Min)
	require.NotZero(t, snap.Battery.Max)
	require.NotZero(t, snap.Battery.ChargeTime)

	require.NotNil(t, snap.Speed)
	require.NotNil(t, snap.Speed.Gen)
	require.NotZero(t, snap.Speed.Val)
	require.NotZero(t, snap.Speed.Min)
	require.NotZero(t, snap.Speed.Max)

	require.NotNil(t, snap.Navigator)
	require.Len(t, snap.Navigator.Routes, 1)
}

func TestDevice_FromSnapshot(t *testing.T) {
	opts := NewDeviceOptions()
	opts.Battery.Min = 1
	opts.Model = "model"
	opts.UserID = "userId"
	opts.Color = "#ffffff"
	opts.Descr = "some descr"
	d, err := NewDevice(opts)
	require.NoError(t, err)
	s1, _ := types.NewSensor("s1", 1, 10, 8, 0)
	d.AddSensor(s1)
	s2, _ := types.NewSensor("s2", 15, 20, 4, 0)
	d.AddSensor(s2)
	d.AddRoute(testRoutes()...)
	for i := 0; i < 10; i++ {
		require.True(t, d.Next(1.0))
	}
	snap := d.Snapshot()

	d2 := new(Device)
	d2.FromSnapshot(snap)

	require.Equal(t, d, d2)
}

func TestDevice_Serialize(t *testing.T) {
	d, _ := NewDevice(nil)
	s1, _ := types.NewSensor("s1", 1, 10, 8, 0)
	d.AddSensor(s1)
	s2, _ := types.NewSensor("s2", 15, 20, 4, 0)
	d.AddSensor(s2)
	d.AddRoute(testRoutes()...)
	require.True(t, d.Next(3.0))

	data, err := d.MarshalBinary()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	d2 := new(Device)
	require.NoError(t, d2.UnmarshalBinary(data))

	d3 := new(Device)
	require.Error(t, d3.UnmarshalBinary(nil))

	require.Equal(t, d, d2)
}

func TestDevice_Duration(t *testing.T) {
	expectedDuration := 3.0

	d := NewTracker()
	d.AddRoute(testRoutes()...)
	d.Next(expectedDuration)

	require.Equal(t, expectedDuration, d.Duration())
}

func TestNewSensor(t *testing.T) {
	type args struct {
		name      string
		min       float64
		max       float64
		amplitude int
		mode      types.SensorMode
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Sensor
		wantErr bool
	}{
		{
			name: "should return error when name is emtpy",
			args: args{
				name: "",
			},
			wantErr: true,
		},
		{
			name: "should return error when amplitude is invalid",
			args: args{
				name:      "s1",
				min:       1,
				max:       2,
				amplitude: -1,
				mode:      0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSensor(tt.args.name, tt.args.min, tt.args.max, tt.args.amplitude, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSensor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSensor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrackerNew(t *testing.T) {
	require.NotNil(t, NewTracker())
	require.NotNil(t, NewKidsTracker())
	require.NotNil(t, NewAnimalTracker())
	require.NotNil(t, NewBicycleTracker())
	require.NotNil(t, NewDroneTracker())
}

func TestDevice_SetUserID(t *testing.T) {
	expectedUserID := uuid.NewString()
	d := NewAnimalTracker()
	require.Empty(t, d.UserID())
	d.SetUserID(expectedUserID)
	require.Equal(t, expectedUserID, d.UserID())
}

func TestDevice_Routes(t *testing.T) {
	expectedRoutes := testRoutes()
	d := NewBicycleTracker()
	d.AddRoute(expectedRoutes...)
	require.Equal(t, expectedRoutes, d.Routes())
}
