package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/navigator"
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/route"
)

func main() {
	ch := make(chan *proto.Device, 2)

	// drone
	routeForDrone1, err := route.Generate()
	if err != nil {
		panic(err)
	}
	drone1, err := gpsgen.Drone("drone1", gpsgen.Properties{"foo": "bar"}, routeForDrone1)
	if err != nil {
		panic(err)
	}
	drone1.OnStateChange = func(state *proto.Device) {
		ch <- state
	}

	// drone with custom sensors
	routeForDrone2, err := route.Generate()
	if err != nil {
		panic(err)
	}
	sensors := []gpsgen.Sensor{
		{
			Name:      "i9-91",
			Min:       1,
			Max:       3,
			Amplitude: gpsgen.Amplitude32,
		},
		{
			Name:      "i7-81",
			Min:       5,
			Max:       6,
			Amplitude: gpsgen.Amplitude4,
		},
	}
	drone2WithSensors, err := gpsgen.DroneWithSensors("drone2", []*navigator.Route{routeForDrone2}, nil, sensors...)
	if err != nil {
		panic(err)
	}
	drone2WithSensors.OnStateChange = func(state *proto.Device) {
		ch <- state
	}

	go func() {
		for s := range ch {
			draw.Table(s)
		}
	}()

	gen := gpsgen.New()
	gen.Attach(drone1)
	gen.Attach(drone2WithSensors)

	gen.Run()

	select {}
}
