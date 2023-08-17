package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mmadfox/go-gpsgen"
)

const (
	numTracksPerRoute = 3
	numTrackers       = 1000
	flushInterval     = 3 * time.Second
)

func main() {
	lon := 37.625616307117696
	lat := 55.75350460378772

	genOpts := gpsgen.NewOptions()
	genOpts.Interval = flushInterval
	gen := gpsgen.New(genOpts)

	// for network transmission
	gen.OnPacket(func(b []byte) {
		packet, err := gpsgen.PacketFromBytes(b)
		if err != nil {
			panic(err)
		}

		fmt.Printf("got packet with numDevices=%d\n", len(packet.Devices))

		for i := 0; i < len(packet.Devices); i++ {
			pck := packet.Devices[i]
			fmt.Printf("%s -> speed:%.2f m/s, tick:%.2f sec, curDist:%.2f meters, lon:%f, lat:%f, el:%f\n",
				pck.Model,
				pck.Speed,
				pck.Tick,
				pck.Distance.CurrentDistance,
				pck.Location.Lon,
				pck.Location.Lat,
				pck.Location.Elevation)
		}
	})

	gen.OnError(func(err error) {
		fmt.Println("[ERROR]", err)
	})

	gen.OnNext(func() {
		fmt.Println("tracker state changed successfully")
	})

	for i := 0; i < numTrackers; i++ {
		tracker := gpsgen.NewTracker()
		tracker.SetUserID(uuid.NewString())
		route := gpsgen.RandomRoute(lon, lat, numTracksPerRoute, gpsgen.RouteLevelM)
		tracker.AddRoute(route)
		gen.Attach(tracker)
	}

	terminate(func() {
		gen.Close()
	})

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
