package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/go-gpsgen"
	"github.com/mmadfox/go-gpsgen/geo"
	"github.com/mmadfox/go-gpsgen/navigator"
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

	track, err := navigator.NewTrack([]geo.LatLonPoint{
		{Lon: 67.67893291456264, Lat: 44.74332748947083},
		{Lon: 64.86705299207989, Lat: 48.249999904172114},
	})
	if err != nil {
		panic(err)
	}
	route := navigator.RouteFromTracks(track)

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
