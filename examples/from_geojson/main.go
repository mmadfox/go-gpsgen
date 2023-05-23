package main

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/geojson"
	"github.com/mmadfox/go-gpsgen/proto"
)

func main() {
	gen := gpsgen.New()

	route1, err := geojson.Decode([]byte(route))
	if err != nil {
		panic(err)
	}

	props := gpsgen.Properties{
		"engine":     "v09XT1",
		"manufactor": "mmadofx",
		"gidro":      "5971",
	}

	tbsSensor := gpsgen.Sensor{
		Name:      "TBS-09-87-11",
		Min:       0.1,
		Max:       0.9,
		Amplitude: gpsgen.Amplitude16,
	}

	myDrone, err := gpsgen.DroneWithSensors("MyDrone-Tx4501", route1, props, tbsSensor)
	if err != nil {
		panic(err)
	}
	myDrone.OnStateChange = func(state *proto.Device) {
		fmt.Printf("model=%s, tick=%d, dist=%f, totalDist=%f, online=%v\n",
			state.Model,
			state.Tick,
			state.Location.CurrentDistance,
			state.Location.TotalDistance,
			state.Online,
		)
	}

	gen.Attach(myDrone)
	gen.Run()
	defer gen.Close()

	select {}
}
