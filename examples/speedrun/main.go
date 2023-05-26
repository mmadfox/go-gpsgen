package main

import (
	"time"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/draw"
	"github.com/mmadfox/go-gpsgen/proto"
)

func main() {
	conf := gpsgen.NewConfig()
	conf.Sensors = []gpsgen.Sensor{
		{Name: "temp", Min: 23, Max: 27, Amplitude: gpsgen.Amplitude16},
		{Name: "scan", Min: 100, Max: 105, Amplitude: gpsgen.Amplitude4},
	}
	myDevice, err := conf.NewDevice()
	if err != nil {
		panic(err)
	}
	myDevice.OnStateChange = func(state *proto.Device) {
		draw.Table(state)
	}

	genInterval := gpsgen.WithInterval(10 * time.Millisecond)
	gen := gpsgen.New(genInterval)
	gen.Attach(myDevice)

	gen.Run()
	defer gen.Close()

	select {}
}
