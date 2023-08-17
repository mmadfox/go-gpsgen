package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/types"
)

func main() {
	gen := gpsgen.New(nil)

	gen.OnPacket(func(data []byte) {
		pck, err := gpsgen.PacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		tracker := pck.Devices[0]
		for i := 0; i < len(tracker.Sensors); i++ {
			sensor := tracker.Sensors[i]
			fmt.Printf("%s -> %f\n", sensor.Name, sensor.ValY)
		}
	})

	droneTracker := gpsgen.NewDroneTracker()
	droneTracker.AddSensor("s1", 1, 10, 16, types.WithStart|types.WithRandom|types.WithEnd)
	droneTracker.AddSensor("s2", 10, 20, 16, types.WithRandom|types.WithEnd)
	droneTracker.AddSensor("s3", 20, 30, 16, 0)

	route := gpsgen.RandomRoute(28.31261399982, 53.247483804819666, 2, gpsgen.RouteLevelXL)

	droneTracker.AddRoute(route)

	terminate(func() {
		gen.Close()
	})

	gen.Attach(droneTracker)
	gen.Run()
}

func terminate(fn func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigChan
		fn()
	}()
}
