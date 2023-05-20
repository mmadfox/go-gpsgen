package gpsgen

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/mmadfox/go-gpsgen/types"
)

// DeviceSetting functions modify a Device
// by applying specific configurations or settings.
type DeviceSetting func(*deviceSettings)

// Sensor represents a custom sensor with its associated properties.
//
// The sensor generates data from the minimum to the maximum value along the bezier curve,
// taking into account the control points that are specified in the amplitude parameter.
//
// Amplitude is the number of control points on the bezier curve from 4 to 512.
type Sensor struct {
	Name      string  `json:"name"`      // specifies the name or identifier of the sensor.
	Min       float64 `json:"min"`       // indicates the minimum value that the sensor can produce.
	Max       float64 `json:"max"`       // indicates the maximum value that the sensor can produce.
	Amplitude int     `json:"amplitude"` // specifies the amplitude or range of the sensor's values from 4 to 512
}

// WithModel sets the model of the device.
func WithModel(model string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.model = model
	}
}

// WithUserID sets the userID of the device.
func WithUserID(userID string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.userID = userID
	}
}

// WithRoutes appends the provided routes to the existing routes in the device.
func WithRoutes(routes []*navigator.Route) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.routes = append(ds.routes, routes...)
	}
}

// WithRoute appends a single route to the existing routes in the device.
func WithRoute(route *navigator.Route) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.routes = append(ds.routes, route)
	}
}

// WithSpeed sets the minimum, maximum, and amplitude speed value in the device.
func WithSpeed(min, max float64, amplitude int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.speed.min = min
		ds.speed.max = max
		ds.speed.amplitude = amplitude
	}
}

// WithBattery sets the minimum and maximum battery values in the device.
func WithBattery(min, max float64) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.battery.min = min
		ds.battery.max = max
	}
}

// WithSensors appends the provided sensors to the existing sensors in the device.
func WithSensors(sensor ...Sensor) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.sensors = append(ds.sensors, sensor...)
	}
}

// WithElevation sets the minimum, maximum, and amplitude of the elevation in the device.
func WithElevation(min, max float64, amplitude int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.elevation.min, ds.elevation.max = min, max
		ds.elevation.amplitude = amplitude
	}
}

// WithOffline sets the minimum and maximum offline values in the device.
func WithOffline(min, max int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.offline.min, ds.offline.max = min, max
	}
}

// WithProps sets the properties in the device.
func WithProps(props Properties) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.props = props
	}
}

// WithDescription sets the description in the device.
func WithDescription(descr string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.descr = descr
	}
}

// Routes takes one or more *navigator.Route arguments and returns a []*navigator.Route.
func Routes(route ...*navigator.Route) []*navigator.Route {
	return append([]*navigator.Route{}, route...)
}

func defaultSettings() *deviceSettings {
	return &deviceSettings{
		id:        uuid.New(),
		model:     types.RandomModel().String(),
		speed:     rangeFloatVal{min: 1, max: 5},
		battery:   rangeFloatVal{min: 50, max: 100},
		elevation: rangeFloatVal{min: 1, max: 300},
		offline:   rangeIntVal{min: 1, max: 60},
	}
}

type deviceSettings struct {
	id        uuid.UUID
	userID    string
	model     string
	speed     rangeFloatVal
	battery   rangeFloatVal
	sensors   []Sensor
	elevation rangeFloatVal
	offline   rangeIntVal
	routes    []*navigator.Route
	props     Properties
	descr     string
}

func (ds *deviceSettings) applyFor(d *Device) (err error) {
	d.model, err = types.NewModel(ds.model)
	if err != nil {
		return err
	}
	d.speed, err = types.NewSpeed(ds.speed.min, ds.speed.max, ds.speed.amplitude)
	if err != nil {
		return err
	}
	d.battery, err = types.NewBattery(ds.battery.min, ds.battery.max)
	if err != nil {
		return err
	}
	if len(ds.sensors) > 0 {
		d.sensors = make([]*types.Sensor, len(ds.sensors))
		for i := 0; i < len(ds.sensors); i++ {
			sensorOpts := ds.sensors[i]
			sensor, err := types.NewSensor(
				sensorOpts.Name,
				sensorOpts.Min,
				sensorOpts.Max,
				int(sensorOpts.Amplitude),
			)
			if err != nil {
				return err
			}
			d.sensors[i] = sensor
		}
	}
	if len(ds.routes) > 0 {
		for i := 0; i < len(ds.routes); i++ {
			if ds.routes[i] == nil {
				return fmt.Errorf("invalid route for this device")
			}
		}
		d.navigator.AddRoutes(ds.routes...)
	} else {
		defaultRoute, err := route.RoutesForRussia()
		if err != nil {
			return err
		}
		d.navigator.AddRoutes(defaultRoute...)
	}
	return nil
}

type rangeFloatVal struct {
	min, max  float64
	amplitude int
}

type rangeIntVal struct {
	min, max int
}
