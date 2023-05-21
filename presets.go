package gpsgen

import (
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/random"
	"github.com/mmadfox/go-gpsgen/route"
)

func RandomDrone() (*Device, error) {
	model := random.String(5)
	routes, err := route.RoutesForRussia()
	if err != nil {
		return nil, err
	}
	return Drone(model, nil, routes...)
}

func RandomTracker() (*Device, error) {
	model := random.String(5)
	routes, err := route.RoutesForRussia()
	if err != nil {
		return nil, err
	}
	return Tracker(model, nil, routes...)
}

func Drone(
	model string,
	props Properties,
	route ...*navigator.Route,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithSpeed(5, 10, Amplitude4),
		WithRoutes(route),
		WithElevation(100, 800, Amplitude8),
		WithBattery(80, 100),
		WithOffline(1, 10),
		WithProps(props),
		WithDescription("Drone"),
	)
}

func DroneWithSensors(
	model string,
	routes []*navigator.Route,
	props Properties,
	sensors ...Sensor,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithRoutes(routes),
		WithSpeed(5, 10, Amplitude8),
		WithElevation(100, 800, Amplitude16),
		WithBattery(80, 100),
		WithSensors(sensors...),
		WithOffline(1, 10),
		WithProps(props),
		WithDescription("Drone"),
	)
}

func Tracker(
	model string,
	props Properties,
	route ...*navigator.Route,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithRoutes(route),
		WithSpeed(1, 3, Amplitude4),
		WithElevation(1, 150, Amplitude8),
		WithBattery(1, 100),
		WithOffline(10, 120),
		WithProps(props),
		WithDescription("Tracker"),
	)
}

func TrackerWithSensors(
	model string,
	routes []*navigator.Route,
	props Properties,
	sensors ...Sensor,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithSpeed(1, 3, Amplitude4),
		WithRoutes(routes),
		WithSensors(sensors...),
		WithElevation(1, 150, Amplitude8),
		WithBattery(1, 100),
		WithOffline(10, 120),
		WithProps(props),
		WithDescription("Tracker"),
	)
}
