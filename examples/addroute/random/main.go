package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/go-gpsgen"
)

func main() {
	gen := gpsgen.New(nil)

	gen.OnPacket(func(data []byte) {
		pck, err := gpsgen.PacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		tracker := pck.Devices[0]
		fmt.Printf("%s -> %f, %f\n",
			tracker.Model,
			tracker.Location.Lon,
			tracker.Location.Lat)
	})

	droneTracker := gpsgen.NewDroneTracker()

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
