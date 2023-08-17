package main

import (
	"fmt"

	"github.com/mmadfox/go-gpsgen"
)

func main() {
	droneTracker := gpsgen.NewDroneTracker()

	route := gpsgen.RandomRoute(28.31261399982, 53.247483804819666, 2, gpsgen.RouteLevelXL)

	droneTracker.AddRoute(route)
	data, err := droneTracker.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Printf("tracker size %d bytes\n", len(data))

	droneTracker2 := new(gpsgen.Device)
	if err := droneTracker2.UnmarshalBinary(data); err != nil {
		panic(err)
	}
}
