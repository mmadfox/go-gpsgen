package main

import (
	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/proto"
	"github.com/mmadfox/go-gpsgen/route"
)

func main() {
	myRoute, err := route.Russia2()
	if err != nil {
		panic(err)
	}

	myDevice, err := gpsgen.NewDevice(
		gpsgen.WithModel("myModel"),
		gpsgen.WithDescription("some description"),
		gpsgen.WithElevation(1, 3, 4),
		gpsgen.WithSpeed(1, 7, 64),
		gpsgen.WithOffline(1, 120),
		gpsgen.WithBattery(0, 100),
		gpsgen.WithProps(gpsgen.Properties{
			"foo": "foo",
			"bar": "bar",
		}),
		gpsgen.WithUserID("12345678"),
		gpsgen.WithSensors(
			gpsgen.Sensor{Name: "s1", Min: 1, Max: 10, Amplitude: 8},
			gpsgen.Sensor{Name: "s2", Min: 3, Max: 10, Amplitude: 128}),
		gpsgen.WithRoute(myRoute),
	)
	if err != nil {
		panic(err)
	}

	myDevice.OnStateChange = func(state *proto.Device) {
		draw.Table(state)
	}

	gen := gpsgen.New()
	gen.Attach(myDevice)

	gen.Run()
	defer gen.Close()

	select {}
}
