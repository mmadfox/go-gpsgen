package gpsgen

import "github.com/mmadfox/go-gpsgen/navigator"

func Drone(
	model string,
	routes []*navigator.Route,
	props Properties,
	sensors ...Sensor,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithSpeed(10, 30, 128),
		WithElevation(100, 2500, 128),
		WithBattery(80, 100),
		WithSensors(sensors...),
		WithOffline(1, 10),
		WithProps(props),
		WithDescritpion("Drone"),
	)
}

func Bicycle(
	model string,
	routes []*navigator.Route,
	props Properties,
	sensors ...Sensor,
) (*Device, error) {
	return NewDevice(
		WithModel(model),
		WithSpeed(2, 10, 16),
		WithElevation(1, 150, 16),
		WithBattery(100, 1),
		WithSensors(sensors...),
		WithOffline(10, 120),
		WithProps(props),
		WithDescritpion("Bicycle"),
	)
}
