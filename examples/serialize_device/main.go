package main

import (
	"fmt"
	"reflect"

	"github.com/mmadfox/go-gpsgen"
)

func main() {
	droneTracker := gpsgen.NewDroneTracker()

	route := gpsgen.RandomRoute(28.31261399982, 53.247483804819666, 2, gpsgen.RouteLevelXL)

	droneTracker.AddRoute(route)
	data, err := gpsgen.EncodeTracker(droneTracker)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tracker size %d bytes\n", len(data))

	droneTracker2, err := gpsgen.DecodeTracker(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("success", reflect.DeepEqual(droneTracker.State(), droneTracker2.State()))
}
