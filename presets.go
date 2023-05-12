package gpsgen

import "github.com/mmadfox/go-gpsgen/navigator"

func Drone(
	model string,
	props Properties,
	route ...*navigator.Route,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithSpeed(10, 30, 128),
		WithRoutes(route),
		WithElevation(100, 2500, 128),
		WithBattery(80, 100),
		WithOffline(1, 10),
		WithProps(props),
		WithDescritpion("Drone"),
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
		WithSpeed(10, 30, 128),
		WithElevation(100, 2500, 128),
		WithBattery(80, 100),
		WithSensors(sensors...),
		WithOffline(1, 10),
		WithProps(props),
		WithDescritpion("Drone"),
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
		WithSpeed(2, 10, 16),
		WithElevation(1, 150, 16),
		WithBattery(1, 100),
		WithOffline(10, 120),
		WithProps(props),
		WithDescritpion("Tracker"),
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
		WithSpeed(2, 10, 16),
		WithRoutes(routes),
		WithSensors(sensors...),
		WithElevation(1, 150, 16),
		WithBattery(1, 100),
		WithOffline(10, 120),
		WithProps(props),
		WithDescritpion("Tracker"),
	)
}
