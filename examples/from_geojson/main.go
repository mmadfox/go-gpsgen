package main

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/geojson"
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
		Name: "TBS-09-87-11",
		Min:  0.1,
		Max:  0.9,
	}

	myDrone, err := gpsgen.Drone("MyDrone-Tx4501", route1, props, tbsSensor)
	if err != nil {
		panic(err)
	}
	myDrone.OnStateChange = func(s *gpsgen.State, snapshot []byte) {
		fmt.Printf("model=%s, tick=%.f, dist=%.f, totalDist=%.f, online=%v, %v \n",
			s.Model,
			s.Tick,
			s.Location.CurrentDistance,
			s.Location.TotalDistance,
			s.Online,
			s.Location,
		)
	}

	gen.Attach(myDrone)
	gen.Run()
	defer gen.Close()

	select {}
}
