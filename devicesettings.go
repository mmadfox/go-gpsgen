package gpsgen

import (
	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/route"
	"github.com/mmadfox/go-gpsgen/types"
)

type DeviceSetting func(*deviceSettings)

func WithModel(model string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.model = model
	}
}

func WithUserID(userID string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.userID = userID
	}
}

func WithRoutes(routes []*navigator.Route) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.routes = append(ds.routes, routes...)
	}
}

func WithRoute(route *navigator.Route) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.routes = append(ds.routes, route)
	}
}

func WithSpeed(min, max float64, amplitude int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.speed.min = min
		ds.speed.max = max
		ds.speed.amplitude = amplitude
	}
}

func WithBattery(min, max float64) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.battery.min, ds.battery.max = min, max
	}
}

func WithSensors(sensor ...Sensor) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.sensors = append(ds.sensors, sensor...)
	}
}

func WithElevation(min, max float64, amplitude int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.elevation.min, ds.elevation.max = min, max
		ds.elevation.amplitude = amplitude
	}
}

func WithOffline(min, max int) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.offline.min, ds.offline.max = min, max
	}
}

func WithProps(props Properties) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.props = props
	}
}

func WithDescritpion(descr string) DeviceSetting {
	return func(ds *deviceSettings) {
		ds.descr = descr
	}
}

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
		d.navigator.AddRoutes(ds.routes...)
	} else {
		defaultRoute, err := route.Generate()
		if err != nil {
			return err
		}
		d.navigator.AddRoute(defaultRoute)
	}
	return nil
}

type Sensor struct {
	Name      string
	Min, Max  float64
	Amplitude int
}

type rangeFloatVal struct {
	min, max  float64
	amplitude int
}

type rangeIntVal struct {
	min, max int
}
